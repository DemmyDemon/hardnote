package ui

import (
	"fmt"
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyleSelected = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0"))
var listStyleUnselected = lipgloss.NewStyle().
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
	offset int
	store  storage.Storage
	index  storage.Index
}

func (ls ListScreen) Init() tea.Cmd {
	return nil
}

func (ls *ListScreen) moveCursorUp() {
	if len(ls.index) == 0 {
		return
	}
	ls.cursor--
	if ls.cursor < 0 {
		ls.cursor = 0
	}
	if ls.cursor < ls.offset {
		ls.offset = ls.cursor
	}
	return
}

func (ls *ListScreen) moveCursorDown() {
	if len(ls.index) == 0 {
		return
	}
	ls.cursor++
	if ls.cursor >= len(ls.index) {
		ls.cursor = len(ls.index) - 1
	}
	if ls.cursor >= ls.height-ls.offset {
		ls.offset++
		if ls.offset >= len(ls.index)-ls.height {
			ls.offset = len(ls.index) - ls.height
		}
	}
	return
}

func (ls ListScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			ls.moveCursorUp()
			return ls, nil // UpdateStatus(fmt.Sprintf("[up] c:%d o: %d", ls.cursor, ls.offset), DirtStateUnchanged)
		case "down":
			ls.moveCursorDown()
			return ls, nil // UpdateStatus(fmt.Sprintf("[down] c:%d o: %d", ls.cursor, ls.offset), DirtStateUnchanged)
		case "alt+up":
			idx, err := ls.store.MoveUp(ls.index[ls.cursor].Id)
			if err != nil {
				return ls, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			ls.moveCursorUp()
			return ls, UpdateIndex(idx)
		case "alt+down":
			idx, err := ls.store.MoveDown(ls.index[ls.cursor].Id)
			if err != nil {
				return ls, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			ls.moveCursorDown()
			return ls, UpdateIndex(idx)
		case "n":
			entry, idx, err := ls.store.Create("", "")
			if err != nil {
				return ls, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			ls.cursor = len(idx) - 1
			if ls.cursor >= ls.height {
				ls.offset = (ls.cursor - ls.height) + 1
			}
			return ls, tea.Batch(UpdateIndex(idx), RequestRename(storage.EntryMeta{Id: entry.Id}))
		case "r":
			return ls, RequestRename(ls.index[ls.cursor])
		case "d":
			return ls, RequestDelete(ls.index[ls.cursor])
		case "enter":
			if len(ls.index) > 0 && ls.cursor <= len(ls.index)-1 {
				return ls, RequestEdit(ls.index[ls.cursor])
			}
		}
	case IndexUpdateMsg:
		ls.index = msg.Index
		if ls.cursor >= len(ls.index) {
			ls.cursor = len(ls.index) - 1
		}
		if ls.cursor < 0 {
			ls.cursor = 0
		}
		if ls.offset > ls.cursor-ls.height {
			ls.offset = (ls.cursor - ls.height) + 1
		}
	case tea.WindowSizeMsg:
		ls.height = msg.Height - 2 // Leave room for header and statusbar
		if ls.height >= len(ls.index) {
			ls.offset = 0
		}
		if ls.height <= ls.cursor {
			ls.offset = ls.cursor
		}
		if ls.cursor >= ls.height-ls.offset {
			if ls.offset >= len(ls.index)-ls.height {
				ls.offset = len(ls.index) - ls.height
			}
		}
		ls.width = msg.Width
	}
	return ls, nil
}

func (ls ListScreen) View() string {
	screen := ""

	for i := 0; i < ls.height; i++ {

		if i+ls.offset >= len(ls.index) {
			break
		}

		idx := i + ls.offset
		if idx < 0 {
			ls.offset = 0
			idx = 0
		}
		if idx >= len(ls.index) {
			idx = len(ls.index) - 1
			ls.offset = idx - ls.height
		}

		entryMeta := ls.index[idx]
		if i == 0 && ls.offset != 0 {
			screen += "↑│"
		} else if i == ls.height-1 && i+ls.offset < len(ls.index)-1 {
			screen += "↓│"
		} else {
			screen += " │"
		}
		name := entryMeta.Name
		if name == "" {
			name = "Untitled"
		}
		if i+ls.offset == ls.cursor {
			screen += listStyleSelected.Render(" " + name + " ")
		} else {
			screen += listStyleUnselected.Render(" " + name + " ")
		}
		screen += "\n"
	}
	screen = strings.TrimSuffix(screen, "\n")
	screen += strings.Repeat("\n │", max(0, ls.height-len(ls.index)))
	screen = strings.TrimPrefix(screen, "\n")
	return fmt.Sprintf("%s\n%s", ls.headerView(), screen)
}

func (ls ListScreen) headerView() string {
	title := "─┤ ↑↓ Select an entry to edit   ↲ Open in editor ├"
	if len(ls.index) == 0 {
		title = "─┤ Press n to create a new entry ├"
	}
	line := strings.Repeat("─", max(0, ls.width-lipgloss.Width(title)))
	return title + line
}
