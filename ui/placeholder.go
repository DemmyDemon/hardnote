package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Placeholder struct {
	height int
	width  int
}

func (ph Placeholder) Init() tea.Cmd {
	return nil
}

func (ph Placeholder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ph.height = msg.Height - 1 // Leave room for status bar
		ph.width = msg.Width
	}
	return ph, nil
}

func (ph Placeholder) View() string {
	title := "──┤ Placeholder ├"
	line := strings.Repeat("─", max(0, ph.width-lipgloss.Width(title)))
	view := lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	view += "\nThe view you are trying to use is not implemented yet."
	view += strings.Repeat("\n", ph.height-2)
	return view
}
