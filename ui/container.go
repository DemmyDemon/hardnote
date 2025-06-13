package ui

import (
	"fmt"

	"github.com/DemmyDemon/hardnote/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type uiState int

const (
	UIStateHelping uiState = iota
	UIStateListing
	UIStateEditing
	UIStateRenaming
	UIStateDeleting
)

type UI struct {
	name      string
	state     uiState
	help      tea.Model
	list      tea.Model
	edit      tea.Model
	rename    tea.Model
	delete    tea.Model
	statusbar tea.Model
	data      storage.Storage
}

func New(name string, data storage.Storage) tea.Model {
	ui := UI{
		name:      name,
		state:     UIStateHelping,
		help:      NewHelpScreen(),
		list:      Placeholder{},
		edit:      Placeholder{},
		rename:    Placeholder{},
		delete:    Placeholder{},
		statusbar: NewStatusbar(name),
		data:      data,
	}

	return ui
}

func SetUiState(newState uiState) tea.Cmd {
	return func() tea.Msg {
		return UIStateUpdate{
			SetState: newState,
		}
	}
}

type UIStateUpdate struct {
	SetState uiState
}

func (ui UI) Init() tea.Cmd {
	return tea.ClearScreen
}

func (ui UI) Distribute(msg tea.Msg) (tea.Model, tea.Cmd) {

	commands := make([]tea.Cmd, 0, 6) // Initialize capacity to the number of models, including status bar

	helpModel, helpCmd := ui.help.Update(msg)
	ui.help = helpModel
	commands = append(commands, helpCmd)

	listModel, listCmd := ui.list.Update(msg)
	ui.list = listModel
	commands = append(commands, listCmd)

	editModel, editCmd := ui.edit.Update(msg)
	ui.edit = editModel
	commands = append(commands, editCmd)

	renameModel, renameCmd := ui.rename.Update(msg)
	ui.rename = renameModel
	commands = append(commands, renameCmd)

	deleteModel, deleteCmd := ui.delete.Update(msg)
	ui.delete = deleteModel
	commands = append(commands, deleteCmd)

	statusModel, statusCmd := ui.statusbar.Update(msg)
	ui.statusbar = statusModel
	commands = append(commands, statusCmd)

	return ui, tea.Batch(commands...)
}

func (ui UI) ToCurrent(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ui.state {
	case UIStateHelping:
		helpModel, helpCmd := ui.help.Update(msg)
		ui.help = helpModel
		if helpCmd != nil {
			return ui, helpCmd
		}
	case UIStateListing:
		listModel, listCmd := ui.list.Update(msg)
		ui.list = listModel
		if listCmd != nil {
			return ui, listCmd
		}
	case UIStateEditing:
		editModel, editCmd := ui.edit.Update(msg)
		ui.edit = editModel
		if editCmd != nil {
			return ui, editCmd
		}
	case UIStateRenaming:
		renameModel, renameCmd := ui.rename.Update(msg)
		ui.rename = renameModel
		if renameCmd != nil {
			return ui, renameCmd
		}
	case UIStateDeleting:
		deleteModel, deleteCmd := ui.delete.Update(msg)
		ui.delete = deleteModel
		if deleteCmd != nil {
			return ui, deleteCmd
		}
	default:
		return ui, UpdateStatus(ui.name, fmt.Sprintf("INVALID STATE %d", ui.state), DirtStateUnchanged)
	}

	statusModel, statusCmd := ui.statusbar.Update(msg)
	ui.statusbar = statusModel
	return ui, statusCmd
}

func (ui UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return ui.Distribute(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ui, tea.Quit
		default:
			return ui.ToCurrent(msg)
		}
	case UIStateUpdate:
		ui.state = msg.SetState
		return ui, nil
	default:
		return ui.ToCurrent(msg)
	}
}

func (ui UI) View() string {
	var s string
	switch ui.state {
	case UIStateListing:
		s = ui.list.View()
	case UIStateEditing:
		s = ui.edit.View()
	default:
		s = ui.help.View()
	}
	s += "\n"
	s += ui.statusbar.View()
	return s
}
