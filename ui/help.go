package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultHelpTitle = "HardNote Help"

var defaultHelpText = []string{
	"Use ↑ and ↓ to scroll the help text if it is too long for your terminal.",
	"Press clrl+l to open the note listing.",
	"",
	"Global keys:",
	"  ctrl+q  exits HardNote, if there are no unsaved changes",
	"  ctrl+c  exits without checking if it's saved.",
	"  ctrl+h  opens this help screen, if there are no unsaved changes",
	"",
	"Listing keys:",
	"    ↑ and ↓   navigates the list.",
	"    alt+↑   moves entry up",
	"    alt+↓   moves entry down",
	"    enter     loads the selected entry into the editor",
	"    r         renames the selected entry",
	"    n         creates a new entry",
	"    d         deletes the selected entry",
	"",
	"Editor keys:",
	"  ctrl+s  saves the current note",
	"  ctrl+d  saves the current note, and opens the listing",
	"  ctrl+u  discards the changes to the current note",
	"  ctrl+l  opens the listing, if the current note is saved",
}

func NewHelpScreen() HelpScreen {
	return HelpScreen{
		title: defaultHelpTitle,
		lines: defaultHelpText,
	}
}

type HelpScreen struct {
	height int
	width  int
	offset int
	title  string
	lines  []string
}

func (h HelpScreen) Init() tea.Cmd {
	return nil
}

func (h *HelpScreen) moveUp() {
	h.offset--
	if h.offset < 0 {
		h.offset = 0
	}
}
func (h *HelpScreen) moveDown() {
	h.offset++
	h.offset = min(h.offset, (len(h.lines))-h.height)
	if h.offset < 0 {
		h.offset = 0
	}
}

func (h HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height - 2 // Leave room for status bar and header
		if h.offset >= len(h.lines)-h.height {
			h.offset = len(h.lines) - h.height
			if h.offset < 0 {
				h.offset = 0
			}
		}
		h.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			h.moveUp()
		case "down":
			h.moveDown()
		}
	}
	return h, nil
}

func (h HelpScreen) View() string {
	var sb strings.Builder
	sb.WriteString(h.headerView())
	sb.WriteRune('\n')
	for i := 0; i < h.height; i++ {
		if i+h.offset >= len(h.lines) {
			break
		}
		if i == 0 && h.offset != 0 {
			sb.WriteString("↑│ ")
		} else if i == h.height-1 && h.offset < (len(h.lines)-h.height) {
			sb.WriteString("↓│ ")
		} else {
			sb.WriteString(" │ ")
		}
		sb.WriteString(h.lines[i+h.offset])

		sb.WriteRune('\n')
	}
	return strings.TrimSuffix(sb.String(), "\n") + strings.Repeat("\n │", max(0, h.height-len(h.lines)))
}
func (h HelpScreen) headerView() string {
	title := fmt.Sprintf("─┤ %s ├", h.title)
	line := strings.Repeat("─", max(0, h.width-lipgloss.Width(title)))
	return title + line
}
