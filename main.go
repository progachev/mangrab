package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/progachev/mangrab/internal/core"
	"github.com/progachev/mangrab/internal/ui"
)

func main() {
	cfg := ui.LoadConfig()

	if len(cfg.Sources) > 0 {
		core.SetImageServers(cfg.Sources)
	}

	dir := cfg.OutputDir
	if dir == "" {
		if d, err := ui.DefaultOutputDir(); err == nil {
			dir = d
		}
	}

	m := ui.NewModel(dir)
	m.SetLang(cfg.Lang)
	m.SetSources(core.ImageServers())
	m.SetSourcesOff(cfg.SourcesOff)
	if cfg.LastMangaURL != "" {
		m.SetInitialURL(cfg.LastMangaURL)
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())
	m.SetProgram(p)

	if _, err := p.Run(); err != nil {
		os.Stderr.WriteString(ui.Tr(cfg.Lang, "app.error", err) + "\n")
		os.Exit(1)
	}
}