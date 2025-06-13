package ui

import (
	"fmt"
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyleSelected = lipgloss.NewStyle().
	MarginLeft(2).
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0")).
	Bold(true)
var listStyleUnselected = lipgloss.NewStyle().
	MarginLeft(2).
	Background(lipgloss.Color("0")).
	Foreground(lipgloss.Color("15"))

type IndexUpdateMsg struct {
	Index storage.Index
}

func UpdateIndex(idx storage.Index) tea.Cmd {
	return func() tea.Msg {
		return IndexUpdateMsg{
			Index: idx,
		}
	}
}

func NewListScreen(data storage.Storage) ListScreen {
	idx, err := data.Index()
	if err != nil {
		panic(err) // This is astronomically unlikely.
	}
	return ListScreen{
		index: idx,
		store: data,
	}
}

type ListScreen struct {
	height int
	width  int
	cursor int
	store  storage.Storage
	index  storage.Index
}

func (ls ListScreen) Init() tea.Cmd {
	return nil
}

func (ls ListScreen) moveCursorUp() ListScreen {
	if len(ls.index) == 0 {
		return ls
	}
	ls.cursor--
	if ls.cursor < 0 {
		ls.cursor = 0
	}
	return ls
}

func (ls ListScreen) moveCursorDown() ListScreen {
	if len(ls.index) == 0 {
		return ls
	}
	ls.cursor++
	if ls.cursor >= len(ls.index) {
		ls.cursor = len(ls.index) - 1
	}
	return ls
}

func (ls ListScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			return ls.moveCursorUp(), nil
		case "down":
			return ls.moveCursorDown(), nil
		case "shift+up":
			idx, err := ls.store.MoveUp(ls.index[ls.cursor].Id)
			if err != nil {
				return ls, UpdateStatus("", err.Error(), DirtStateUnchanged)
			}
			return ls.moveCursorUp(), UpdateIndex(idx)
		case "shift+down":
			idx, err := ls.store.MoveDown(ls.index[ls.cursor].Id)
			if err != nil {
				return ls, UpdateStatus("", err.Error(), DirtStateUnchanged)
			}
			return ls.moveCursorDown(), UpdateIndex(idx)
		case "n":
			entry, idx, err := ls.store.Create("", "")
			if err != nil {
				return ls, UpdateStatus("", err.Error(), DirtStateUnchanged)
			}
			ls.cursor = len(idx) - 1
			return ls, tea.Batch(UpdateIndex(idx), RequestRename(storage.EntryMeta{Id: entry.Id}))
		case "r":
			return ls, RequestRename(ls.index[ls.cursor])
		case "d":
			return ls, RequestDelete(ls.index[ls.cursor])
		}
	case IndexUpdateMsg:
		ls.index = msg.Index
		if ls.cursor >= len(ls.index) {
			ls.cursor = len(ls.index) - 1
		}
		if ls.cursor < 0 {
			ls.cursor = 0
		}
	case tea.WindowSizeMsg:
		ls.height = msg.Height - 2 // Leave room for header and statusbar
		ls.width = msg.Width
	}
	return ls, nil
}

func (ls ListScreen) View() string {
	screen := ""
	if len(ls.index) == 0 {
		screen = listStyleSelected.Render("There are no entries. Press n to create one.")
	}
	for i, entryMeta := range ls.index {
		name := entryMeta.Name
		if name == "" {
			name = "Untitled"
		}
		if i == ls.cursor {
			screen += listStyleSelected.Render(name) + "\n"
		} else {
			screen += listStyleUnselected.Render(name) + "\n"
		}
	}
	screen += strings.Repeat("\n", ls.height-len(ls.index))
	return fmt.Sprintf("%s\n%s", ls.headerView(), screen)
}

func (ls ListScreen) headerView() string {
	title := "──┤ Select an entry to edit ├"
	line := strings.Repeat("─", max(0, ls.width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
