package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/progachev/mangrab/internal/core"
)

func (m *Model) View() string {
	switch m.screen {
	case screenHome:
		return m.homeView()
	case screenChapters, screenDownloading:
		return m.workspaceView()
	case screenDone:
		return m.doneView()
	case screenSettings:
		return m.settingsView()
	}
	return ""
}

// ─── Подсказки с подсветкой клавиш ─────────────────────────────────────────

// isComboKey сообщает, является ли клавиша «спец»-комбинацией
// (стрелки, Tab, Esc, Ctrl+…, Enter и т.п.) — их подсвечиваем другим цветом.
func isComboKey(k string) bool {
	switch k {
	case "↑↓", "↑", "↓", "Tab", "Esc", "Enter", "Ctrl+C", "PgUp", "PgDn", "Home", "End":
		return true
	}
	if strings.HasPrefix(k, "Ctrl+") {
		return true
	}
	return false
}

// renderHints рендерит структурированную подсказку с цветными клавишами.
// Если для name нет записи в hintsTable — возвращает пустую строку
// (вызывающий код может откатиться к m.t(name)).
func (m *Model) renderHints(name string) string {
	hs := m.hints(name)
	if len(hs) == 0 {
		return hintStyle.Render(m.t(name))
	}
	var b strings.Builder
	for i, h := range hs {
		if i > 0 {
			b.WriteString(keySepStyle.Render(" · "))
		}
		if isComboKey(h.key) {
			b.WriteString(keyComboStyle.Render(h.key))
		} else {
			b.WriteString(keyStyle.Render(h.key))
		}
		b.WriteString(" ")
		b.WriteString(keyDescStyle.Render(h.desc))
	}
	return b.String()
}

// ─── HOME ──────────────────────────────────────────────────────────────────

func (m *Model) homeView() string {
	logo := renderLogo()
	sub := lipgloss.NewStyle().Foreground(cTextMuted).Italic(true).
		Render(m.t("home.sub"))
	versionLine := hintStyle.Render(m.t("home.version"))

	inputCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorderHi).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			sectionLabelStyle.Render(m.t("source.label")),
			m.urlInput.View(),
		))

	hint := m.renderHints("home.hint")

	parts := []string{
		logo, "", sub, "", versionLine, "", "", inputCard, "", hint,
	}
	if m.err != nil {
		parts = append(parts, "", errorStyle.Render("✖ "+m.err.Error()))
	}
	if m.busy {
		parts = append(parts, "", badgeBusyStyle.Render(m.spinner.View()+" "+m.t("home.loading")))
	}
	return lipgloss.NewStyle().
		Width(max(40, m.width)).
		Align(lipgloss.Center).
		Render(strings.Join(parts, "\n"))
}

// ─── WORKSPACE ─────────────────────────────────────────────────────────────

func (m *Model) workspaceView() string {
	parts := []string{
		m.workspaceHeader(),
		m.sourceReadonlyView(),
		m.chaptersView(),
		m.downloadCardView(),
	}
	if m.err != nil {
		parts = append(parts, "", errorStyle.Render("✖ "+m.err.Error()))
	}
	parts = append(parts, "", m.workspaceFooter())
	return strings.Join(parts, "\n")
}

func (m *Model) workspaceHeader() string {
	title := titleStyle.Render(m.t("header.title"))
	sub := hintStyle.Render(m.t("header.subtitle", nz(m.slug)))
	left := lipgloss.JoinVertical(lipgloss.Left, title, sub)

	var badge string
	switch {
	case m.busy:
		badge = badgeBusyStyle.Render(m.spinner.View() + " " + m.t("header.working"))
	case m.err != nil:
		badge = badgeErrorStyle.Render("● " + m.t("header.error"))
	default:
		badge = badgeReadyStyle.Render("● " + m.t("header.ready"))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, left, "    ", badge)
}

func (m *Model) sourceReadonlyView() string {
	label := sectionLabelStyle.Render(m.t("source.label"))
	u := textStyle.Render(m.urlInput.Value())
	body := lipgloss.JoinVertical(lipgloss.Left, label, u)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1, 2).
		Width(max(40, m.width-8)).
		Render(body)
}

// chaptersView — карточка со списком глав.
// Порядок: шапка → список → разделитель → строка подсказки (внизу).
func (m *Model) chaptersView() string {
	// ── Шапка ──────────────────────────────────────────────────────────────
	label := sectionLabelStyle.Render(m.t("chapters.label"))

	total := len(m.chapters)
	sel := len(m.selected())

	meta := hintStyle.Render(m.t("chapters.meta", total, sel))
	head := lipgloss.JoinHorizontal(lipgloss.Left, label, "  ", meta)

	return m.renderChaptersCard(head, m.chapters)
}

// renderChaptersCard собирает карточку: шапка, список, строка подсказки снизу.
func (m *Model) renderChaptersCard(head string, visible []core.Chapter) string {
	const maxRows = 12
	start := 0
	if m.cursor >= maxRows {
		start = m.cursor - maxRows + 1
	}
	end := start + maxRows
	if end > len(visible) {
		end = len(visible)
	}

	var rows strings.Builder
	if len(visible) == 0 {
		rows.WriteString(hintStyle.Render("  " + m.t("chapters.none")))
		rows.WriteString("\n")
	} else {
		for i := start; i < end; i++ {
			c := visible[i]
			check := "[ ]"
			if c.Selected {
				check = "[x]"
			}
			marker := " "
			if i == m.cursor {
				marker = "▸"
			}
			line := fmt.Sprintf(" %s %s  %-6s  vol %-4s  %s",
				marker, check, nz(c.Number), nz(c.Volume), c.Name)
			switch {
			case i == m.cursor:
				rows.WriteString(rowSelectedStyle.Render(line))
			case c.Selected:
				rows.WriteString(rowCheckedStyle.Render(line))
			default:
				rows.WriteString(rowNormalStyle.Render(line))
			}
			rows.WriteString("\n")
		}
		if len(visible) > maxRows {
			rows.WriteString(hintStyle.Render(m.t("chapters.range", start+1, end, len(visible))))
			rows.WriteString("\n")
		}
	}

	// строка с цветными горячими клавишами
	toolsLine := lipgloss.NewStyle().
		MaxWidth(max(20, m.width-12)).
		Render("  " + m.renderHints("chapters.tools"))

	body := lipgloss.JoinVertical(lipgloss.Left,
		head,
		"",
		rows.String(),
		"",
		toolsLine, // ← внизу
	)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1, 2).
		Width(max(40, m.width-8)).
		Render(body)
}

func (m *Model) downloadCardView() string {
	label := sectionLabelStyle.Render(m.t("output.label"))

	folder := hintStyle.Render("▸ ")
	if m.outputDir == "" {
		folder += hintStyle.Render(m.t("output.none"))
	} else {
		folder += textStyle.Render(m.outputDir)
	}

	status := hintStyle.Render("▸ ") + textStyle.Render(m.status)
	progress := m.progress.View()

	body := lipgloss.JoinVertical(lipgloss.Left, label, folder, status, progress)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1, 2).
		Width(max(40, m.width-8)).
		Render(body)
}

func (m *Model) workspaceFooter() string {
	return m.renderHints("footer.workspace")
}

// ─── DONE ──────────────────────────────────────────────────────────────────

// doneView — экран результатов со скроллом по объединённому списку
// созданных файлов и ошибок. Курсор подсвечивает текущую строку,
// чтобы при большом числе глав можно было ориентироваться.
func (m *Model) doneView() string {
	total := len(m.results) + len(m.failed)

	var b strings.Builder

	// ── Заголовок + сводка ─────────────────────────────────────────────────
	b.WriteString(titleStyle.Render(m.t("done.title")))
	b.WriteString("\n\n")

	if len(m.results) > 0 {
		b.WriteString(badgeReadyStyle.Render(
			fmt.Sprintf("● "+m.t("done.created"), len(m.results))))
		b.WriteString("\n")
	}
	if len(m.failed) > 0 {
		b.WriteString(badgeErrorStyle.Render(
			fmt.Sprintf("● "+m.t("done.failed"), len(m.failed))))
		b.WriteString("\n")
	}
	if total == 0 {
		b.WriteString(hintStyle.Render("  " + m.t("done.nothing")))
		b.WriteString("\n")
	}

	// ── Папка сохранения ───────────────────────────────────────────────────
	if m.outputDir != "" {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(m.t("done.folder")))
		b.WriteString(textStyle.Render(m.outputDir))
		b.WriteString("\n")
	}

	if total > 0 {
		b.WriteString("\n")

		// ── Список с прокруткой ────────────────────────────────────────────
		rows := m.doneVisibleRows()
		start := 0
		if m.doneCursor >= rows {
			start = m.doneCursor - rows + 1
		}
		end := start + rows
		if end > total {
			end = total
		}

		for i := start; i < end; i++ {
			marker := "  "
			if i == m.doneCursor {
				marker = "▸ "
			}

			var line string
			if i < len(m.results) {
				line = marker + "✓ " + m.results[i]
				if i == m.doneCursor {
					b.WriteString(rowSelectedStyle.Render(line))
				} else {
					b.WriteString(textStyle.Render(line))
				}
			} else {
				errIdx := i - len(m.results)
				line = marker + "✗ " + m.failed[errIdx]
				if i == m.doneCursor {
					b.WriteString(rowSelectedStyle.Render(line))
				} else {
					b.WriteString(errorStyle.Render(line))
				}
			}
			b.WriteString("\n")
		}

		if total > rows {
			b.WriteString(hintStyle.Render(
				fmt.Sprintf("  … %d–%d из %d", start+1, end, total)))
			b.WriteString("\n")
		}
	}

	// ── Подсказка с цветными клавишами ────────────────────────────────────
	b.WriteString("\n")
	b.WriteString(m.renderHints("done.hint"))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1, 2).
		Width(max(40, m.width-8)).
		Render(b.String())
}

// ─── SETTINGS ──────────────────────────────────────────────────────────────

func (m *Model) settingsView() string {
	st := m.settings
	if st == nil {
		return ""
	}

	row := func(field int, text string) string {
		if st.cursor == field {
			return rowSelectedStyle.Render("▸ " + text)
		}
		return "  " + textStyle.Render(text)
	}

	outputBlock := lipgloss.JoinVertical(lipgloss.Left,
		row(fieldOutput, m.t("settings.output")),
		"  "+st.output.View(),
	)
	langBlock := lipgloss.JoinVertical(lipgloss.Left,
		row(fieldLang, m.t("settings.lang")),
		"  "+st.lang.View(),
	)

	sourcesHeader := row(fieldSources, m.t("settings.sources"))

	var srcLines []string
	if len(st.sources) == 0 {
		srcLines = append(srcLines, hintStyle.Render("  "+m.t("settings.sources.empty")))
	} else {
		for i, it := range st.sources {
			marker := "  "
			if st.cursor == fieldSources && i == st.srcCursor && !st.srcEditing {
				marker = rowSelectedStyle.Render("▸ ")
			}

			check := "[ ]"
			checkStyle := rowNormalStyle
			if it.enabled {
				check = "[x]"
				checkStyle = rowCheckedStyle
			}

			var valueView string
			if st.srcEditing && st.cursor == fieldSources && i == st.srcCursor {
				valueView = st.srcInput.View()
			} else {
				valueView = checkStyle.Render(it.value)
			}

			line := marker + checkStyle.Render(check+" ") + valueView
			if st.cursor == fieldSources && i == st.srcCursor && !st.srcEditing {
				line = rowSelectedStyle.Render("▸ " + check + " " + it.value)
			}
			srcLines = append(srcLines, line)
		}
	}

	srcHint := m.renderHints("settings.sources.hint")
	sourcesBlock := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{sourcesHeader}, append(srcLines, "", srcHint)...)...)

	body := lipgloss.JoinVertical(lipgloss.Left,
		sectionLabelStyle.Render(m.t("settings.title")),
		"",
		outputBlock,
		"",
		langBlock,
		"",
		sourcesBlock,
		"",
		m.renderHints("settings.hint"),
	)

	if m.err != nil {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "",
			errorStyle.Render("✖ "+m.err.Error()))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cBorderHi).
		Padding(1, 2).
		Width(max(40, m.width-8)).
		Render(body)
}

// ─── Общие утилиты ─────────────────────────────────────────────────────────

func nz(s string) string {
	if s == "" {
		return "—"
	}
	return s
}