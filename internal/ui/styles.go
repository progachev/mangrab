package ui

import "github.com/charmbracelet/lipgloss"

var (
	cText      = lipgloss.Color("#c0caf5")
	cTextMuted = lipgloss.Color("#565f89")
	cTextDim   = lipgloss.Color("#414868")

	cAccent   = lipgloss.Color("#7aa2f7")
	cSuccess  = lipgloss.Color("#9ece6a")
	cWarning  = lipgloss.Color("#e0af68")
	cError    = lipgloss.Color("#f7768e")
	cSelected = lipgloss.Color("#ff9e64")

	// Цвета для горячих клавиш.
	cKey     = lipgloss.Color("#ff9e64") // оранжевый — обычные клавиши
	cKeyAlt  = lipgloss.Color("#7dcfff") // голубой — спец-клавиши/комбинации
	cKeyDesc = lipgloss.Color("#a9b1d6") // светлый — описание
	cKeySep  = lipgloss.Color("#3b4261") // тёмный — разделитель ·

	cBorder   = lipgloss.Color("#565f89")
	cBorderHi = lipgloss.Color("#7aa2f7")
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(cText)

	sectionLabelStyle = lipgloss.NewStyle().Foreground(cAccent).Bold(true)

	textStyle = lipgloss.NewStyle().Foreground(cText)
	hintStyle = lipgloss.NewStyle().Foreground(cTextMuted)

	rowSelectedStyle = lipgloss.NewStyle().Foreground(cSelected).Bold(true)
	rowCheckedStyle  = lipgloss.NewStyle().Foreground(cText)
	rowNormalStyle   = lipgloss.NewStyle().Foreground(cTextMuted)

	badgeReadyStyle = lipgloss.NewStyle().Foreground(cSuccess).Bold(true)
	badgeBusyStyle  = lipgloss.NewStyle().Foreground(cAccent).Bold(true)
	badgeErrorStyle = lipgloss.NewStyle().Foreground(cError).Bold(true)

	errorStyle = lipgloss.NewStyle().Foreground(cError)

	// Стили для цветных горячих клавиш в подсказках.
	keyStyle      = lipgloss.NewStyle().Foreground(cKey).Bold(true)
	keyComboStyle = lipgloss.NewStyle().Foreground(cKeyAlt).Bold(true)
	keyDescStyle  = lipgloss.NewStyle().Foreground(cKeyDesc)
	keySepStyle   = lipgloss.NewStyle().Foreground(cKeySep)
)