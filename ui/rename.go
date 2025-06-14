package ui

import (
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewRenameScreen(data storage.Storage) RenameScreen {
	ti := textinput.New()
	ti.Placeholder = "Untitled"
	ti.Focus()
	ti.CharLimit = 60
	ti.Width = 20
	return RenameScreen{
		input: ti,
		store: data,
	}
}

func RequestRename(entryMeta storage.EntryMeta) tea.Cmd {
	return func() tea.Msg {
		return RenameRequestMsg{
			EntryMeta: entryMeta,
		}
	}
}

type RenameRequestMsg struct {
	EntryMeta storage.EntryMeta
}

type RenameScreen struct {
	height    int
	store     storage.Storage
	entryMeta storage.EntryMeta
	input     textinput.Model
}

func (rs RenameScreen) Init() tea.Cmd {
	return nil
}

func (rs RenameScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RenameRequestMsg:
		rs.entryMeta = msg.EntryMeta
		rs.input.SetValue(msg.EntryMeta.Name)
		rs.input.SetCursor(len(msg.EntryMeta.Name))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return rs, SetUiState(UIStateListing)
		case "enter":
			idx, err := rs.store.Rename(rs.entryMeta.Id, rs.input.Value())
			if err != nil {
				return rs, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			return rs, tea.Batch(UpdateIndex(idx), SetUiState(UIStateListing))
		}
	case tea.WindowSizeMsg:
		rs.height = msg.Height - 3 // Leave room for prompt, input and statusbar
		rs.input.Width = msg.Width
	}

	inputModel, inputCmd := rs.input.Update(msg)
	rs.input = inputModel
	return rs, inputCmd
}

func (rs RenameScreen) View() string {
	s := "What would you like this entry to be named?\n"
	s += rs.input.View()
	s += strings.Repeat("\n", rs.height)
	return s
}
