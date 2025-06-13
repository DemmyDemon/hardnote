package ui

import tea "github.com/charmbracelet/bubbletea"

type Placeholder struct{}

func (ph Placeholder) Init() tea.Cmd {
	return nil
}

func (ph Placeholder) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return ph, nil
}

func (ph Placeholder) View() string {
	return "[PLACEHOLDER]"
}
