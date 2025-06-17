package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		ph.height = msg.Height - 2 // Leave room for header and status bar
		ph.width = msg.Width
	}
	return ph, nil
}

func (ph Placeholder) View() string {
	view := unifiedHeader("Placeholder", ph.width)
	view += " │ The view you are trying to use is not implemented yet."
	view += strings.Repeat("\n │", ph.height-1)
	return view
}
