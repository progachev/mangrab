package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

const (
	fieldOutput = iota
	fieldLang
	fieldSources
	fieldCount
)

// sourceItem — одна строка в списке источников.
type sourceItem struct {
	value   string
	enabled bool
}

type settingsState struct {
	cursor int // индекс поля (fieldOutput / fieldLang / fieldSources)

	output textinput.Model
	lang   textinput.Model

	// sources — управляемый список.
	sources    []sourceItem
	srcCursor  int             // индекс выбранной строки в списке источников
	srcEditing bool            // true, если сейчас редактируем value выбранной строки
	srcInput   textinput.Model // поле для редактирования/добавления URL

	fields []*textinput.Model
}

func newSettingsState(m *Model) *settingsState {
	output := textinput.New()
	output.SetValue(m.outputDir)
	output.Placeholder = "/path/to/output"
	output.CharLimit = 512
	output.Width = 60

	lg := textinput.New()
	lg.SetValue(m.lang)
	lg.Placeholder = "ru"
	lg.CharLimit = 8 // запас, чтобы можно было свободно править
	lg.Width = 8

	// Включённые — сначала, затем выключенные.
	items := make([]sourceItem, 0, len(m.sources)+len(m.sourcesOff))
	for _, s := range m.sources {
		items = append(items, sourceItem{value: s, enabled: true})
	}
	for _, s := range m.sourcesOff {
		items = append(items, sourceItem{value: s, enabled: false})
	}

	srcInput := textinput.New()
	srcInput.Placeholder = "https://img3.cdnlibs.org"
	srcInput.CharLimit = 512
	srcInput.Width = 60

	st := &settingsState{
		output:   output,
		lang:     lg,
		sources:  items,
		srcInput: srcInput,
	}
	st.fields = []*textinput.Model{&st.output, &st.lang}
	st.syncFocus()
	return st
}

func (s *settingsState) syncFocus() {
	for i, f := range s.fields {
		if i == s.cursor && s.cursor != fieldSources {
			(*f).Focus()
		} else {
			(*f).Blur()
		}
	}
	if s.cursor == fieldSources && s.srcEditing {
		s.srcInput.Focus()
	} else {
		s.srcInput.Blur()
	}
}

// enabledSources возвращает список включённых источников.
func (s *settingsState) enabledSources() []string {
	out := make([]string, 0, len(s.sources))
	for _, it := range s.sources {
		if it.enabled && strings.TrimSpace(it.value) != "" {
			out = append(out, strings.TrimSpace(it.value))
		}
	}
	return out
}

// disabledSources возвращает список выключенных источников.
func (s *settingsState) disabledSources() []string {
	out := make([]string, 0, len(s.sources))
	for _, it := range s.sources {
		v := strings.TrimSpace(it.value)
		if v == "" || it.enabled {
			continue
		}
		out = append(out, v)
	}
	return out
}

func joinSources(s []string) string {
	return strings.Join(s, ",")
}