package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var helpText = []string{
	"HARDNOTE",
	"This will eventually be the fooken help text, yeah?",
	"I mean, eventually...",
}

func NewHelpScreen() HelpScreen {
	return HelpScreen{}
}

type HelpScreen struct {
	height int
	width  int
}

func (h HelpScreen) Init() tea.Cmd {
	return nil
}

func (h HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height - 1 // Leave room for status bar
		h.width = msg.Width
	}
	return h, nil
}

func (h HelpScreen) View() string {
	if len(helpText) >= h.height {
		return strings.Join(helpText[:h.height], "\n")
	}
	msg := strings.Join(helpText, "\n")
	msg += strings.Repeat("\n", h.height-len(helpText))
	return msg
}
