package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var pickStyleSelected = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0")).
	PaddingLeft(2).PaddingRight(1)

var pickStyleUnselected = lipgloss.NewStyle().
	Background(lipgloss.Color("0")).
	Foreground(lipgloss.Color("15")).
	PaddingLeft(1).PaddingRight(1)

func PickOne(prompt string, options []string, action PickOneAction) tea.Cmd {
	return func() tea.Msg {
		return PickOneRequestMsg{
			Prompt:  prompt,
			Options: options,
			Action:  action,
		}
	}
}

type PickOneAction func(selected int) tea.Cmd

type PickOneRequestMsg struct {
	Prompt  string
	Options []string
	Action  PickOneAction
}

func NewPickOneScreen() PickOneScreen {
	return PickOneScreen{}
}

type PickOneScreen struct {
	height  int
	width   int
	request PickOneRequestMsg
	cursor  int
}

func (po PickOneScreen) Init() tea.Cmd {
	return nil
}

func (po PickOneScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PickOneRequestMsg:
		po.request = msg
		if po.request.Options == nil || len(po.request.Options) == 0 {
			po.request.Options = []string{"No", "Yes"}
		}
		po.cursor = 0
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return po, tea.Batch(UpdateStatus("Aborted selection", DirtStateUnchanged), SetUiState(UIStateListing))
		case "enter":
			if po.request.Action == nil {
				return po, tea.Batch(UpdateStatus("Select has no action?!", DirtStateUnchanged), SetUiState(UIStateListing))
			}
			return po, po.request.Action(po.cursor)
		case "up":
			po.cursor--
			if po.cursor < 0 {
				po.cursor = 0
			}
		case "down":
			po.cursor++
			if po.cursor >= len(po.request.Options) {
				po.cursor = len(po.request.Options) - 1
			}
		}
	case tea.WindowSizeMsg:
		po.height = msg.Height - 2 // Leave room for prompt, input and statusbar
		po.width = msg.Width
	}
	return po, nil
}

func (po PickOneScreen) View() string {
	s := fmt.Sprintf("─┤ %s ├", po.request.Prompt)
	s += strings.Repeat("─", max(0, po.width-lipgloss.Width(s)))
	s += "\n"
	for i, option := range po.request.Options {
		if i == po.cursor {
			s += fmt.Sprintf(" %s\n", pickStyleSelected.Render(option))
		} else {
			s += fmt.Sprintf(" │%s\n", pickStyleUnselected.Render(option))
		}
	}
	s += strings.Repeat(" │\n", max(0, po.height-len(po.request.Options)))
	return strings.TrimSuffix(s, "\n")
}
