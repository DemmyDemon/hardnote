package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const readSizeLimit = 3 * 1024 * 1024 // 3MiB

var nonWordChars = regexp.MustCompile("[^\\w]+")

var listStyleSelected = lipgloss.NewStyle().
	Background(lipgloss.Color("15")).
	Foreground(lipgloss.Color("0")).
	PaddingLeft(2).PaddingRight(1)
var listStyleUnselected = lipgloss.NewStyle().
	Background(lipgloss.Color("0")).
	Foreground(lipgloss.Color("15")).
	PaddingLeft(1).PaddingRight(1)

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

func (ls *ListScreen) moveCursorUp() {
	if len(ls.index) == 0 {
		return
	}
	ls.cursor--
	if ls.cursor < 0 {
		ls.cursor = 0
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
	return
}

func (ls ListScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return ls, tea.Quit
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
			return ls, Ask(
				"What do you want name this entry?",
				"",
				"Untitled",
				func(answer string) tea.Cmd {
					idx, err := ls.store.Rename(entry.Id, answer)
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					return tea.Batch(UpdateIndex(idx), SetUiState(UIStateListing))
				},
			)
		case "r":
			entryMeta := ls.index[ls.cursor]
			return ls, Ask(
				"What do you want to rename it to?",
				entryMeta.Name,
				"Untitled",
				func(answer string) tea.Cmd {
					idx, err := ls.store.Rename(entryMeta.Id, answer)
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					return tea.Batch(UpdateIndex(idx), SetUiState(UIStateListing))
				},
			)
		case "d":
			entryMeta := ls.index[ls.cursor]
			return ls, PickOne(
				fmt.Sprintf("Delete %q?", entryMeta.Name),
				[]string{"No", "Yes, delete for ever!"},
				func(selected int) tea.Cmd {
					if selected == 1 {
						idx, err := ls.store.Delete(entryMeta.Id)
						if err != nil {
							UpdateStatus(err.Error(), DirtStateUnchanged)
						}
						return tea.Batch(
							UpdateStatus(fmt.Sprintf("Deleted %s", entryMeta.Name), DirtStateUnchanged),
							UpdateIndex(idx),
							SetUiState(UIStateListing),
						)
					}
					return tea.Batch(
						UpdateStatus("Okay, never mind.", DirtStateUnchanged),
						SetUiState(UIStateListing),
					)
				},
			)
		case "enter":
			if len(ls.index) > 0 && ls.cursor <= len(ls.index)-1 {
				return ls, RequestEdit(ls.index[ls.cursor])
			}
		case "ctrl+e":
			entryMeta := ls.index[ls.cursor]
			return ls, Ask(
				"Where do you want to export?",
				filename(entryMeta.Name),
				"Enter a filename",
				func(filename string) tea.Cmd {
					_, err := os.Stat(filename)
					if err != nil && !errors.Is(err, os.ErrNotExist) {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					if err == nil {
						return UpdateStatus("File exists. Refusing to overwrite.", DirtStateUnchanged)
					}
					entry, err := ls.store.Read(entryMeta.Id)
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					if err := os.WriteFile(filename, []byte(entry.Text), 0600); err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					return tea.Batch(
						UpdateStatus("Export successful!", DirtStateUnchanged),
						SetUiState(UIStateListing),
					)
				},
			)
		case "ctrl+r":
			wd, err := os.Getwd()
			if err != nil {
				return ls, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			return ls, Ask(
				"What file do you want to read?",
				wd+string(os.PathSeparator),
				"Enter a filename",
				func(filename string) tea.Cmd {
					stat, err := os.Stat(filename)
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					if stat.IsDir() {
						return UpdateStatus("Can't read a whole directory.", DirtStateUnchanged)
					}
					if stat.Size() > readSizeLimit {
						return UpdateStatus("No. HardNote does not do well with files that size.", DirtStateUnchanged)
					}
					data, err := os.ReadFile(filename)
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					_, idx, err := ls.store.Create(filepath.Base(filename), string(data))
					if err != nil {
						return UpdateStatus(err.Error(), DirtStateUnchanged)
					}
					return tea.Batch(
						UpdateIndex(idx),
						SetUiState(UIStateListing),
						UpdateStatus("Read successful", DirtStateUnchanged),
					)
				},
			)
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
		ls.cursor = min(ls.cursor, len(ls.index)-1)
		ls.width = msg.Width
	}
	return ls, nil
}

func (ls ListScreen) View() string {
	screen := ""

	start := max(0, (ls.cursor)-(ls.height/2))
	end := start + ls.height

	if end > len(ls.index) {
		end = len(ls.index)
		start = max(0, end-ls.height)
	}

	for i := start; i < end; i++ {
		entryMeta := ls.index[i]
		if i-start == 0 && start != 0 {
			screen += "↑"
		} else if i-start == ls.height-1 && end < len(ls.index) {
			screen += "↓"
		} else {
			screen += " "
		}
		name := entryMeta.Name
		if name == "" {
			name = "Untitled"
		}
		// name += fmt.Sprintf(" start:%d end:%d i:%d ls.cursor:%d len(ls.index)-1:%d", start, end, i, ls.cursor, len(ls.index)-1)
		if i == ls.cursor {
			screen += fmt.Sprintf("%s\n", listStyleSelected.Render(name))
		} else {
			screen += fmt.Sprintf("│%s\n", listStyleUnselected.Render(name))
		}
	}
	screen = strings.TrimSuffix(screen, "\n")

	screen += strings.Repeat("\n │", max(0, ls.height-(end-start)))
	screen = strings.TrimPrefix(screen, "\n")
	return fmt.Sprintf("%s\n%s", ls.headerView(), screen)
}

func (ls ListScreen) headerView() string {
	title := "─┤ ↑↓ Select an entry to edit ├"
	if len(ls.index) == 0 {
		title = "─┤ Press n to create a new entry ├"
	}
	line := strings.Repeat("─", max(0, ls.width-lipgloss.Width(title)))
	return title + line
}

func filename(original string) string {
	filename := strings.ToLower(original)
	filename = nonWordChars.ReplaceAllString(filename, "_")
	filename = strings.Trim(filename, "_")
	filename += ".txt"
	wd, err := os.Getwd()
	if err != nil { // ... what would that error even be?!
		return filename
	}
	filename = filepath.Join(wd, filename)
	return filename
}
