package ui

import "github.com/charmbracelet/lipgloss"

var logoLines = []string{
	` ███╗   ███╗ █████╗ ███╗   ██╗ ██████╗ ██████╗  █████╗ ██████╗ `,
	` ████╗ ████║██╔══██╗████╗  ██║██╔════╝ ██╔══██╗██╔══██╗██╔══██╗`,
	` ██╔████╔██║███████║██╔██╗ ██║██║  ███╗██████╔╝███████║██████╔╝`,
	` ██║╚██╔╝██║██╔══██║██║╚██╗██║██║   ██║██╔══██╗██╔══██║██╔══██╗`,
	` ██║ ╚═╝ ██║██║  ██║██║ ╚████║╚██████╔╝██║  ██║██║  ██║██████╔╝`,
	` ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ `,
}

func renderLogo() string {
	stops := []string{
		"#7aa2f7", "#8aa8f5", "#98aef3",
		"#a6b4f1", "#b4baef", "#c2c0ed",
	}
	var out string
	for i, line := range logoLines {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(stops[i])).
			Bold(true)
		out += style.Render(line) + "\n"
	}
	return out
}