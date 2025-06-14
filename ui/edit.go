package ui

import (
	"fmt"
	"strings"

	"github.com/DemmyDemon/hardnote/storage"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditRequestMsg struct {
	EntryMeta storage.EntryMeta
}

func RequestEdit(entryMeta storage.EntryMeta) tea.Cmd {
	return func() tea.Msg {
		return EditRequestMsg{
			EntryMeta: entryMeta,
		}
	}
}

func NewEditScreen(data storage.Storage) EditScreen {
	ta := textarea.New()
	ta.Prompt = " │ "
	ta.ShowLineNumbers = false
	ta.EndOfBufferCharacter = '•'
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	ta.Focus()
	return EditScreen{
		text:  ta,
		store: data,
	}
}

type EditScreen struct {
	height int
	width  int
	store  storage.Storage
	entry  storage.Entry
	name   string
	text   textarea.Model
}

func (es EditScreen) Init() tea.Cmd {
	return nil
}

func (es EditScreen) Name() string {
	if es.name != "" {
		return es.name
	}
	return "Untitled"
}

func (es *EditScreen) cursorToBeginningFoulSmellingHack() {
	for es.text.Line() > 0 {
		es.text.CursorUp()
	}
	es.text.CursorStart()
}

func (es EditScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var passCmd tea.Cmd
	switch msg := msg.(type) {
	case EditRequestMsg:
		entry, err := es.store.Read(msg.EntryMeta.Id)
		if err != nil {
			return es, tea.Batch(UpdateStatus(err.Error(), DirtStateUnchanged), SetUiState(UIStateListing))
		}
		es.name = msg.EntryMeta.Name
		es.entry = entry
		es.text.SetValue(entry.Text)
		es.cursorToBeginningFoulSmellingHack()
		return es, tea.Batch(UpdateStatus(fmt.Sprintf("Loaded! %d", es.text.Line()), DirtStateClean), UpdateStatusName(es.Name()))
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+h", "ctrl+l", "\x00": // Noop, let statusbar handle
		case "up", "down", "left", "right", "home", "end", "ctrl+home", "ctrl+end": // Noop, does not change value
		case "ctrl+s":
			es.entry.Text = es.text.Value()
			err := es.store.Update(es.entry)
			if err != nil {
				return es, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			return es, UpdateStatus("Saved!", DirtStateClean)
		case "ctrl+d":
			es.entry.Text = es.text.Value()
			err := es.store.Update(es.entry)
			if err != nil {
				return es, UpdateStatus(err.Error(), DirtStateUnchanged)
			}
			return es, tea.Batch(SetUiState(UIStateListing), UpdateStatus(es.name+" saved!", DirtStateClean))
		case "ctrl+u":
			es.text.SetValue(es.entry.Text)
			es.cursorToBeginningFoulSmellingHack()
			return es, UpdateStatus("Reverted!", DirtStateClean)
		case "esc":
			if es.text.Value() == es.entry.Text {
				return es, tea.Batch(SetUiState(UIStateListing), UpdateStatus("Escape successful!", DirtStateClean))
			}
			return es, UpdateStatus("You can't escape with unsaved changes!", DirtStateDirty)
		default:
			passCmd = UpdateStatus(fmt.Sprintf("(editor) %q", msg.String()), DirtStateDirty)
		}
	case tea.WindowSizeMsg:
		es.height = msg.Height - 2 // Leave room for status bar and header
		es.text.SetHeight(es.height)

		es.width = msg.Width
		es.text.SetWidth(es.width)
	}
	model, cmd := es.text.Update(msg)
	es.text = model
	return es, tea.Batch(passCmd, cmd)
}

func (es EditScreen) View() string {
	return fmt.Sprintf("%s\n%s", es.viewHeader(), es.text.View())
}
func (es EditScreen) viewHeader() string {
	info := es.text.LineInfo()
	header := fmt.Sprintf("─┤ %d:%d ├", es.text.Line()+1, info.CharOffset+info.StartColumn)
	header += strings.Repeat("─", max(0, es.width-lipgloss.Width(header)))
	return header
}
