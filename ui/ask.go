package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Ask(question string, answer string, placeholder string, action AskAnswerAction) tea.Cmd {
	return func() tea.Msg {
		return AskRequestMsg{
			Question:    question,
			Answer:      answer,
			Placeholder: placeholder,
			Action:      action,
		}
	}
}

type AskAnswerAction func(answer string) tea.Cmd

type AskRequestMsg struct {
	Question    string
	Answer      string
	Placeholder string
	Action      AskAnswerAction
}

func NewAskScreen() AskScreen {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20
	ti.Prompt = " │ » "
	return AskScreen{
		input: ti,
		question: AskRequestMsg{
			Question: "Is there a bug?",
			Action: func(answer string) tea.Cmd {
				return UpdateStatus("You forgot to pass a question.", DirtStateUnchanged)
			},
		},
	}
}

type AskScreen struct {
	height   int
	width    int
	question AskRequestMsg
	input    textinput.Model
}

func (as AskScreen) Init() tea.Cmd {
	return nil
}

func (as AskScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AskRequestMsg:
		as.question = msg
		as.input.SetValue(msg.Answer)
		as.input.SetCursor(len(msg.Answer))
		as.input.Placeholder = msg.Placeholder
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return as, tea.Batch(UpdateStatus("Aborted question", DirtStateUnchanged), SetUiState(UIStateListing))
		case "enter":
			value := as.input.Value()
			if as.question.Action != nil {
				return as, as.question.Action(value)
			}
			return as, tea.Batch(UpdateStatus("Question has no action?!", DirtStateUnchanged), SetUiState(UIStateListing))
		}
	case tea.WindowSizeMsg:
		as.height = msg.Height - 3 // Leave room for prompt, input and statusbar
		as.width = msg.Width
		as.input.Width = msg.Width
	}

	inputModel, inputCmd := as.input.Update(msg)
	as.input = inputModel
	return as, inputCmd
}

func (as AskScreen) View() string {
	s := unifiedHeader(as.question.Question, as.width)
	s += as.input.View()
	s += strings.Repeat("\n │", as.height)
	return s
}
