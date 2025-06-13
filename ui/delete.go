package ui

import (
	"fmt"
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func NewDeleteScreen(data storage.Storage) DeleteScreen {
	return DeleteScreen{
		store: data,
	}
}

func RequestDelete(entryMeta storage.EntryMeta) tea.Cmd {
	return func() tea.Msg {
		return DeleteRequestMsg{
			EntryMeta: entryMeta,
		}
	}
}

type DeleteRequestMsg struct {
	EntryMeta storage.EntryMeta
}

type DeleteScreen struct {
	height    int
	store     storage.Storage
	entryMeta storage.EntryMeta
}

func (ds DeleteScreen) Init() tea.Cmd {
	return nil
}

func (ds DeleteScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DeleteRequestMsg:
		ds.entryMeta = msg.EntryMeta
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return ds, SetUiState(UIStateListing)
		case "ctrl+y":
			idx, err := ds.store.Delete(ds.entryMeta.Id)
			if err != nil {
				return ds, UpdateStatus("", err.Error(), DirtStateUnchanged)
			}
			return ds, tea.Batch(UpdateIndex(idx), SetUiState(UIStateListing))
		}
	case tea.WindowSizeMsg:
		ds.height = msg.Height - 5 // Leave room for prompt and statusbar
	}

	return ds, nil
}

func (ds DeleteScreen) View() string {
	s := fmt.Sprintf(
		"You are about to delete\n%s\n%s\nPress ctrl+y to confirm, or esc to abort.",
		ds.entryMeta.Name, ds.entryMeta.Id,
	)
	s += strings.Repeat("\n", ds.height)
	return s
}
