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
	UIStateAsking
)

type UI struct {
	name      string
	state     uiState
	help      tea.Model
	list      tea.Model
	edit      tea.Model
	delete    tea.Model
	ask       tea.Model
	statusbar tea.Model
	data      storage.Storage
}

func New(name string, data storage.Storage) tea.Model {
	ui := UI{
		name:      name,
		state:     UIStateListing,
		help:      NewHelpScreen(),
		list:      NewListScreen(data),
		edit:      NewEditScreen(data),
		delete:    NewDeleteScreen(data),
		ask:       NewAskScreen(),
		statusbar: NewStatusbar(name),
		data:      data,
	}

	return ui
}

func SetUiState(newState uiState) tea.Cmd {
	return func() tea.Msg {
		return UIStateUpdateMsg{
			SetState: newState,
		}
	}
}

type UIStateUpdateMsg struct {
	SetState uiState
}

func (ui UI) Init() tea.Cmd {
	return tea.ClearScreen
}

func (ui UI) Distribute(msg tea.Msg) (tea.Model, tea.Cmd) {

	commands := make([]tea.Cmd, 0, 7) // Initialize capacity to the number of models, including status bar

	helpModel, helpCmd := ui.help.Update(msg)
	ui.help = helpModel
	commands = append(commands, helpCmd)

	listModel, listCmd := ui.list.Update(msg)
	ui.list = listModel
	commands = append(commands, listCmd)

	editModel, editCmd := ui.edit.Update(msg)
	ui.edit = editModel
	commands = append(commands, editCmd)

	deleteModel, deleteCmd := ui.delete.Update(msg)
	ui.delete = deleteModel
	commands = append(commands, deleteCmd)

	askModel, askCmd := ui.ask.Update(msg)
	ui.ask = askModel
	commands = append(commands, askCmd)

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
	case UIStateDeleting:
		deleteModel, deleteCmd := ui.delete.Update(msg)
		ui.delete = deleteModel
		if deleteCmd != nil {
			return ui, deleteCmd
		}
	case UIStateAsking:
		askModel, askCmd := ui.ask.Update(msg)
		ui.ask = askModel
		if askCmd != nil {
			return ui, askCmd
		}
	default:
		return ui, UpdateStatus(fmt.Sprintf("INVALID STATE %d", ui.state), DirtStateUnchanged)
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
		}
	case UIStateUpdateMsg:
		ui.state = msg.SetState
		return ui, UpdateStatusName("")
	case DeleteRequestMsg:
		ui.state = UIStateDeleting
	case EditRequestMsg:
		ui.state = UIStateEditing
	case AskRequestMsg:
		ui.state = UIStateAsking
	case IndexUpdateMsg:
		model, cmd := ui.list.Update(msg)
		ui.list = model
		return ui, cmd
	}
	return ui.ToCurrent(msg)
}

func (ui UI) View() string {
	var s string
	switch ui.state {
	case UIStateListing:
		s = ui.list.View()
	case UIStateEditing:
		s = ui.edit.View()
	case UIStateDeleting:
		s = ui.delete.View()
	case UIStateAsking:
		s = ui.ask.View()
	default:
		s = ui.help.View()
	}
	s += "\n"
	s += ui.statusbar.View()
	return s
}
