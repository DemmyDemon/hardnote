package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpText = []string{
	"If you are seeing this, the password was correct.",
	"",
	"Use ↑ and ↓ to scroll the help text if it is too long for your  terminal.",
	"Press clrl+l to open the note listing.",
	"",
	"Global keys:",
	"  ctrl+q  exits HardNote, if there are no unsaved changes",
	"  ctrl+c  exits without checking if it's saved.",
	"  ctrl+h  opens this help screen, if there are no unsaved changes",
	"",
	"Listing keys:",
	"    ↑ and ↓   navigates the list.",
	"    Shift+↑   moves entry up",
	"    Shift+↓   moves entry down",
	"    enter     loads the selected entry into the editor",
	"    r         renames the selected entry",
	"    n         creates a new entry",
	"",
	"Editor keys:",
	"  ctrl+s  saves the current note",
	"  ctrl+u  discards the changes to the current note",
	"  ctrl+l  opens the listing, if the current note is saved",
	"",
	"Additional keybinds as I make them up, I guess.",
}

func NewHelpScreen() HelpScreen {
	return HelpScreen{}
}

type HelpScreen struct {
	height   int
	width    int
	viewport viewport.Model
}

func (h HelpScreen) Init() tea.Cmd {
	h.viewport = viewport.New(70, 20)
	h.viewport.SetContent(strings.Join(helpText, "\n"))
	return nil
}

func (h HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height - 1        // Leave room for status bar
		h.viewport.Height = h.height - 1 // Leave room for the header

		h.width = msg.Width
		h.viewport.Width = msg.Width
		h.viewport.YPosition = 1 // So it's not behind the header, I guess? Cargo cult.
		h.viewport.SetContent(strings.Join(helpText, "\n"))
	}

	viewport, cmd := h.viewport.Update(msg)
	h.viewport = viewport
	return h, cmd
}

func (h HelpScreen) View() string {
	return fmt.Sprintf("%s\n%s", h.headerView(), h.viewport.View())
}
func (h HelpScreen) headerView() string {
	title := "──┤ HardNote Help ├"
	line := strings.Repeat("─", max(0, h.width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
