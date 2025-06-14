package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dirtState int

const (
	DirtStateUnchanged dirtState = iota
	DirtStateDirty
	DirtStateClean
)

func UpdateStatus(message string, dirt dirtState) tea.Cmd {
	return func() tea.Msg {
		return StatusbarUpdateMsg{
			Message: message,
			Dirt:    dirt,
		}
	}
}

type StatusbarUpdateMsg struct {
	Message string
	Dirt    dirtState
}

func UpdateStatusName(name string) tea.Cmd {
	return func() tea.Msg {
		return StatusNameUpdateMsg(name)
	}
}

type StatusNameUpdateMsg string

var clean = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("10"))
var dirty = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("9"))

func NewStatusbar(filename string) Statusbar {
	return Statusbar{
		file:    filename,
		message: "Press ctrl+h for the help screen",
	}
}

type Statusbar struct {
	file    string
	name    string
	message string
	width   int
	dirty   bool
}

func (sb Statusbar) IsDirty() bool {
	return sb.dirty
}

func (sb Statusbar) Init() tea.Cmd {
	return nil
}

func (sb Statusbar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			if !sb.dirty {
				return sb, tea.Quit
			}
			sb.message = "ctrl+s to save before quitting, or ctrl+c to insist"
		case "ctrl+h":
			if !sb.dirty {
				sb.message = "Opening help"
				return sb, SetUiState(UIStateHelping)
			}
			sb.message = "ctrl+s to save before viewing help, ctrl+u to discard changes"
		case "ctrl+l":
			if !sb.dirty {
				sb.message = "Listing notes"
				return sb, SetUiState(UIStateListing)
			}
			sb.message = "ctrl+s to save before viewing list, ctrl+u to discard changes"
		default:
			sb.message = fmt.Sprintf("%q", msg.String())
		}
	case tea.WindowSizeMsg:
		sb.width = msg.Width
	case StatusbarUpdateMsg:
		sb.message = msg.Message
		if msg.Dirt != DirtStateUnchanged {
			sb.dirty = (msg.Dirt == DirtStateDirty)
		}
	case StatusNameUpdateMsg:
		sb.name = string(msg)
	}
	return sb, nil
}

func (sb Statusbar) View() string {
	var name = sb.name
	if name != "" {
		if sb.dirty {
			name = dirty.Render(sb.name)
		} else {
			name = clean.Render(sb.name)
		}
		return fmt.Sprintf("%s: %s → %s", sb.file, name, sb.message)
	}
	return fmt.Sprintf("%s → %s", sb.file, sb.message)
}
