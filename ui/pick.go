package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var pickStyleSelected = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBackground(lipgloss.Color("0")).
	BorderForeground(lipgloss.Color("15")).
	BorderLeft(true).
	PaddingLeft(1).PaddingRight(1)

var pickStyleUnselected = lipgloss.NewStyle().
	Background(lipgloss.Color("0")).
	Foreground(lipgloss.Color("15")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBackground(lipgloss.Color("0")).
	BorderForeground(lipgloss.Color("15")).
	BorderLeft(true).
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
		po.height = msg.Height - 2 // Leave room for header and statusbar
		po.width = msg.Width
	}
	return po, nil
}

func (po PickOneScreen) View() string {

	var screen strings.Builder
	screen.WriteString(unifiedHeader(po.request.Prompt, po.width))

	start := max(0, (po.cursor - (po.height / 2)))
	end := start + po.height

	if end > len(po.request.Options) {
		end = len(po.request.Options)
		start = max(0, end-po.height)
	}

	for i := start; i < end; i++ {

		if i-start == 0 && start != 0 {
			screen.WriteRune('↑')
		} else if i-start == po.height-1 && end < len(po.request.Options) {
			screen.WriteRune('↓')
		} else {
			screen.WriteRune(' ')
		}

		option := po.request.Options[i]
		if i == po.cursor {
			screen.WriteString(pickStyleSelected.Render(option))
		} else {
			screen.WriteString(pickStyleUnselected.Render(option))
		}
		screen.WriteRune('\n')
	}

	screen.WriteString(strings.Repeat(" │\n", max(0, po.height-len(po.request.Options))))
	return strings.TrimSuffix(screen.String(), "\n")
}
