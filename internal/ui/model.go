package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/progachev/mangrab/internal/core"
	"github.com/progachev/mangrab/internal/fsutil"
)

type screen int

const (
	screenHome screen = iota
	screenChapters
	screenDownloading
	screenDone
	screenSettings
)

type Model struct {
	urlInput textinput.Model

	chapters []core.Chapter
	cursor   int
	slug     string

	// Экран результатов: скролл по объединённому списку файлов.
	doneCursor int

	outputDir string

	spinner  spinner.Model
	progress progress.Model
	status   string

	screen  screen
	busy    bool
	err     error
	results []string
	failed  []string

	cancel  context.CancelFunc
	program *tea.Program

	lang       string
	sources    []string
	sourcesOff []string
	settings   *settingsState

	width, height int
}

func NewModel(defaultDir string) Model {
	url := textinput.New()
	url.Placeholder = "https://mangalib.org/ru/manga/..."
	url.CharLimit = 512
	url.Width = 60
	url.Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	pr := progress.New(
		progress.WithScaledGradient("#7aa2f7", "#9ece6a"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return Model{
		urlInput:  url,
		spinner:   sp,
		progress:  pr,
		outputDir: defaultDir,
		status:    "Вставьте ссылку и нажмите Enter",
		screen:    screenHome,
		lang:      "ru",
		sources:   core.ImageServers(),
	}
}

func (m *Model) SetLang(l string) {
	if l == "en" {
		m.lang = "en"
	} else {
		m.lang = "ru"
	}
}

func (m *Model) SetSources(s []string)     { m.sources = s }
func (m *Model) SetSourcesOff(s []string)  { m.sourcesOff = s }
func (m *Model) SetProgram(p *tea.Program) { m.program = p }

func (m *Model) SetInitialURL(u string) {
	m.urlInput.SetValue(u)
	m.urlInput.CursorEnd()
}

func (m *Model) persistConfig() {
	_ = SaveConfig(config{
		LastMangaURL: m.urlInput.Value(),
		OutputDir:    m.outputDir,
		Lang:         m.lang,
		Sources:      m.sources,
		SourcesOff:   m.sourcesOff,
	})
}

// ─── сообщения ─────────────────────────────────────────────────────────────

type chaptersLoadedMsg struct{ manga *core.Manga }
type chapterErrorMsg struct{ err error }
type downloadProgressMsg struct {
	Percent float64
	Text    string
}
type downloadDoneMsg struct {
	Created []string
	Failed  []string
	Err     error
}
type dirChosenMsg struct{ Dir string }

// ─── Init ──────────────────────────────────────────────────────────────────

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// ─── Update ────────────────────────────────────────────────────────────────

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.progress.Width = max(20, msg.Width-20)
		m.resizeInputs()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "s" &&
			!m.busy &&
			m.screen != screenSettings &&
			m.screen != screenDownloading &&
			m.canOpenSettings() {
			m.openSettings()
			return m, nil
		}
		return m.handleKey(msg)

	case chaptersLoadedMsg:
		m.busy = false
		m.err = nil
		m.slug = msg.manga.Slug
		m.chapters = dedupChapters(msg.manga.Chapters)
		m.urlInput.SetValue(msg.manga.URL)
		m.urlInput.Blur()
		m.screen = screenChapters
		m.cursor = 0
		m.status = m.t("status.loaded", len(m.chapters))
		return m, nil

	case chapterErrorMsg:
		m.busy = false
		m.err = msg.err
		m.status = m.t("status.error_list")
		return m, nil

	case dirChosenMsg:
		m.outputDir = msg.Dir
		m.status = m.t("status.folder", msg.Dir)
		return m, nil

	case downloadProgressMsg:
		cmd := m.progress.SetPercent(msg.Percent / 100.0)
		m.status = msg.Text
		return m, cmd

	case downloadDoneMsg:
		m.busy = false
		m.err = msg.Err
		m.results = msg.Created
		m.failed = msg.Failed
		m.screen = screenDone
		m.doneCursor = 0
		if msg.Err != nil {
			m.status = m.t("status.aborted", msg.Err)
		} else if len(msg.Failed) == 0 {
			m.status = m.t("status.done", len(msg.Created))
		} else {
			m.status = m.t("status.done_err", len(msg.Created), len(msg.Failed))
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	// ⬇️ fallback: не-KeyMsg (Blink, etc.) прокидываем в активное поле.
	if m.screen == screenSettings && m.settings != nil {
		st := m.settings
		if st.srcEditing {
			var cmd tea.Cmd
			st.srcInput, cmd = st.srcInput.Update(msg)
			return m, cmd
		}
		if st.cursor == fieldOutput || st.cursor == fieldLang {
			var cmd tea.Cmd
			*st.fields[st.cursor], cmd = (*st.fields[st.cursor]).Update(msg)
			return m, cmd
		}
		return m, nil
	}
	if m.screen == screenHome && m.urlInput.Focused() {
		var cmd tea.Cmd
		m.urlInput, cmd = m.urlInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// dedupChapters убирает записи с одинаковой парой (Number, Volume).
func dedupChapters(in []core.Chapter) []core.Chapter {
	if len(in) == 0 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]core.Chapter, 0, len(in))
	for _, c := range in {
		key := c.Number + "\x00" + c.Volume
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, c)
	}
	return out
}

// resizeInputs пересчитывает ширину инпутов под текущий размер окна.
func (m *Model) resizeInputs() {
	w := max(20, m.width-20)
	m.urlInput.Width = w
}

func (m *Model) canOpenSettings() bool {
	switch m.screen {
	case screenHome:
		return !m.urlInput.Focused()
	case screenChapters:
		return true
	case screenDone:
		return true
	default:
		return false
	}
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case "esc":
		if m.screen == screenDownloading && m.cancel != nil {
			m.cancel()
			m.status = m.t("status.home")
			return m, nil
		}
	}

	switch m.screen {
	case screenHome:
		return m.handleHomeScreen(msg)
	case screenSettings:
		return m.handleSettingsScreen(msg)
	case screenChapters:
		return m.handleChaptersScreen(msg)
	case screenDownloading:
		return m, nil
	case screenDone:
		return m.handleDoneScreen(msg)
	}

	return m, nil
}

// handleDoneScreen — навигация по объединённому списку created+failed.
func (m *Model) handleDoneScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(m.results) + len(m.failed)
	switch msg.String() {
	case "enter":
		m.screen = screenChapters
		m.results = nil
		m.failed = nil
		m.err = nil
		m.doneCursor = 0
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.doneCursor > 0 {
			m.doneCursor--
		}
	case "down", "j":
		if m.doneCursor < total-1 {
			m.doneCursor++
		}
	case "g", "home":
		m.doneCursor = 0
	case "G", "end":
		if total > 0 {
			m.doneCursor = total - 1
		}
	case "pgup":
		m.doneCursor -= m.doneVisibleRows()
		if m.doneCursor < 0 {
			m.doneCursor = 0
		}
	case "pgdown":
		m.doneCursor += m.doneVisibleRows()
		if m.doneCursor > total-1 {
			m.doneCursor = total - 1
		}
	}
	return m, nil
}

// doneVisibleRows — сколько строк списка влезает на экран результатов.
func (m *Model) doneVisibleRows() int {
	if m.height <= 0 {
		return 12
	}
	// шапка+футер+отступы ~ 10 строк
	return max(4, m.height-10)
}

func (m *Model) openSettings() {
	m.settings = newSettingsState(m)
	m.screen = screenSettings
}

func (m *Model) backFromSettings() {
	m.settings = nil
	if m.chapters == nil {
		m.screen = screenHome
		m.urlInput.Focus()
	} else {
		m.screen = screenChapters
	}
}

func (m *Model) handleSettingsScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	st := m.settings
	if st == nil {
		m.backFromSettings()
		return m, nil
	}

	// ── 1. Режим редактирования источника: всё в srcInput ─────────────────
	if st.srcEditing {
		switch msg.String() {
		case "esc":
			st.srcEditing = false
			st.srcInput.Blur()
			st.srcInput.SetValue("")
			return m, nil
		case "enter":
			val := strings.TrimSpace(st.srcInput.Value())
			if val != "" {
				if st.srcCursor < len(st.sources) {
					st.sources[st.srcCursor].value = val
				} else {
					st.sources = append(st.sources, sourceItem{value: val, enabled: true})
					st.srcCursor = len(st.sources) - 1
				}
			}
			st.srcEditing = false
			st.srcInput.Blur()
			st.srcInput.SetValue("")
			return m, nil
		default:
			var cmd tea.Cmd
			st.srcInput, cmd = st.srcInput.Update(msg)
			return m, cmd
		}
	}

	// ── 2. Глобальные клавиши: сохранение, отмена, навигация, Tab ─────────
	switch msg.String() {
	case "esc":
		m.backFromSettings()
		return m, nil
	case "enter":
		l := strings.ToLower(strings.TrimSpace(st.lang.Value()))
		if l != "en" && l != "ru" {
			m.err = fmt.Errorf("%s", m.t("settings.bad_lang"))
			return m, nil
		}

		m.outputDir = strings.TrimSpace(st.output.Value())
		m.lang = l

		m.sources = st.enabledSources()
		m.sourcesOff = st.disabledSources()
		core.SetImageServers(m.sources)

		m.persistConfig()
		m.err = nil
		m.status = m.t("settings.saved")
		m.backFromSettings()
		return m, nil
	case "tab":
		st.cursor = (st.cursor + 1) % fieldCount
		st.syncFocus()
		return m, nil
	case "up", "k":
		// Если курсор в текстовом поле — стрелки вверх/вниз не должны
		// блокировать перемещение по полям. Но если пользователь хочет
		// переместиться по списку источников — работаем как раньше.
		if st.cursor == fieldSources {
			if st.srcCursor > 0 {
				st.srcCursor--
			}
			return m, nil
		}
		if st.cursor > 0 {
			st.cursor--
		}
		st.syncFocus()
		return m, nil
	case "down", "j":
		if st.cursor == fieldSources {
			if st.srcCursor < len(st.sources)-1 {
				st.srcCursor++
			}
			return m, nil
		}
		if st.cursor < fieldCount-1 {
			st.cursor++
		}
		st.syncFocus()
		return m, nil
	}

	// ── 3. Курсор в текстовом поле (output / lang): всё остальное → в поле ─
	if st.cursor == fieldOutput || st.cursor == fieldLang {
		var cmd tea.Cmd
		*st.fields[st.cursor], cmd = (*st.fields[st.cursor]).Update(msg)
		return m, cmd
	}

	// ── 4. Курсор на списке источников: горячие клавиши списка ────────────
	if st.cursor == fieldSources {
		switch msg.String() {
		case " ":
			if st.srcCursor >= 0 && st.srcCursor < len(st.sources) {
				st.sources[st.srcCursor].enabled = !st.sources[st.srcCursor].enabled
			}
			return m, nil
		case "a":
			st.srcEditing = true
			st.srcCursor = len(st.sources)
			st.srcInput.SetValue("")
			st.srcInput.Focus()
			return m, textinput.Blink
		case "e":
			if st.srcCursor >= 0 && st.srcCursor < len(st.sources) {
				st.srcEditing = true
				st.srcInput.SetValue(st.sources[st.srcCursor].value)
				st.srcInput.CursorEnd()
				st.srcInput.Focus()
				return m, textinput.Blink
			}
			return m, nil
		case "d", "delete", "backspace":
			if st.srcCursor >= 0 && st.srcCursor < len(st.sources) {
				st.sources = append(st.sources[:st.srcCursor], st.sources[st.srcCursor+1:]...)
				if st.srcCursor >= len(st.sources) && st.srcCursor > 0 {
					st.srcCursor--
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) handleHomeScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		value := strings.TrimSpace(m.urlInput.Value())
		if value == "" {
			m.err = fmt.Errorf("%s", m.t("home.empty_url"))
			return m, nil
		}
		m.urlInput.Blur()
		m.busy = true
		m.err = nil
		m.status = m.t("status.search")
		return m, tea.Batch(m.spinner.Tick, findChaptersCmd(value))
	}
	var cmd tea.Cmd
	m.urlInput, cmd = m.urlInput.Update(msg)
	return m, cmd
}

func (m *Model) handleChaptersScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if !m.busy {
			m.screen = screenHome
			m.urlInput.Focus()
			m.chapters = nil
			m.cursor = 0
			m.slug = ""
			m.err = nil
			m.status = m.t("status.cancel")
			return m, textinput.Blink
		}
		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.chapters)-1 {
			m.cursor++
		}
	case " ", "space":
		if m.cursor >= 0 && m.cursor < len(m.chapters) {
			m.chapters[m.cursor].Selected = !m.chapters[m.cursor].Selected
		}
	case "a":
		for i := range m.chapters {
			m.chapters[i].Selected = true
		}
	case "n":
		for i := range m.chapters {
			m.chapters[i].Selected = false
		}
	case "o":
		return m, chooseDirCmd(m.outputDir)
	case "d":
		selected := m.selected()
		if len(selected) == 0 {
			m.err = fmt.Errorf("%s", m.t("chapters.pick_one"))
			return m, nil
		}
		if m.outputDir == "" {
			m.err = fmt.Errorf("%s", m.t("chapters.no_dir"))
			return m, nil
		}
		if err := os.MkdirAll(m.outputDir, 0o755); err != nil {
			m.err = fmt.Errorf(m.t("chapters.mkdir"), err)
			return m, nil
		}
		ctx, cancel := context.WithCancel(context.Background())
		m.cancel = cancel
		m.screen = screenDownloading
		m.busy = true
		m.err = nil
		m.progress.SetPercent(0)
		m.persistConfig()
		return m, tea.Batch(
			m.spinner.Tick,
			downloadCmd(m.program, ctx, m.slug, selected, m.outputDir, m.t),
		)
	case "q":
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) selected() []core.Chapter {
	out := make([]core.Chapter, 0)
	for _, c := range m.chapters {
		if c.Selected {
			out = append(out, c)
		}
	}
	return out
}

// ─── команды ───────────────────────────────────────────────────────────────

func findChaptersCmd(rawURL string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()
		manga, err := core.FetchChapters(ctx, rawURL)
		if err != nil {
			return chapterErrorMsg{err}
		}
		return chaptersLoadedMsg{manga}
	}
}

func chooseDirCmd(defaultDir string) tea.Cmd {
	return func() tea.Msg {
		if env := strings.TrimSpace(os.Getenv("MLLIB_DIR")); env != "" {
			abs, err := filepath.Abs(env)
			if err == nil {
				return dirChosenMsg{Dir: abs}
			}
		}
		return dirChosenMsg{Dir: defaultDir}
	}
}

func downloadCmd(
	program *tea.Program,
	ctx context.Context,
	slug string,
	chapters []core.Chapter,
	dir string,
	tr func(string, ...any) string,
) tea.Cmd {
	return func() tea.Msg {
		used := map[string]bool{}
		var created []string
		var failed []string

		for i, ch := range chapters {
			if ctx.Err() != nil {
				return downloadDoneMsg{Created: created, Failed: failed, Err: ctx.Err()}
			}
			label := tr("download.chapter", ch.Number)
			basePercent := float64(i) / float64(len(chapters)) * 100

			if program != nil {
				program.Send(downloadProgressMsg{
					Percent: basePercent,
					Text:    tr("download.getpages") + label,
				})
			}

			pages, err := core.FetchChapterPages(ctx, slug, ch)
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s: %v", label, err))
				continue
			}

			pdfBytes, err := core.CreateChapterPDF(ctx, ch, pages, func(pageIdx, total int, text string) {
				if program == nil || total == 0 {
					return
				}
				p := (float64(i) + float64(pageIdx)/float64(total)) / float64(len(chapters)) * 100
				program.Send(downloadProgressMsg{
					Percent: p,
					Text:    label + ": " + tr("download.page", pageIdx+1, total),
				})
			})
			if err != nil {
				if ctx.Err() != nil {
					return downloadDoneMsg{Created: created, Failed: failed, Err: ctx.Err()}
				}
				failed = append(failed, fmt.Sprintf("%s: %v", label, err))
				continue
			}

			name := fsutil.ChapterFilename(ch, used)
			path := fsutil.UniquePath(dir, name)
			if err := os.WriteFile(path, pdfBytes, 0o644); err != nil {
				failed = append(failed, fmt.Sprintf("%s: %v", label, err))
				continue
			}
			created = append(created, filepath.Base(path))

			if program != nil {
				program.Send(downloadProgressMsg{
					Percent: float64(i+1) / float64(len(chapters)) * 100,
					Text:    tr("download.done") + filepath.Base(path),
				})
			}
		}

		return downloadDoneMsg{Created: created, Failed: failed}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}