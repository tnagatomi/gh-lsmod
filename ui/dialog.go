package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// Dialog represents a confirmation dialog
type Dialog struct {
	title   string
	message string
	width   int
	height  int
	styles  DialogStyles
	keyMap  DialogKeyMap
}

// DialogStyles contains the styles for the dialog
type DialogStyles struct {
	Border       lipgloss.Style
	Title        lipgloss.Style
	Message      lipgloss.Style
	Button       lipgloss.Style
	ButtonActive lipgloss.Style
}

// DefaultDialogStyles returns the default styles for the dialog
func DefaultDialogStyles() DialogStyles {
	return DialogStyles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("240")).
			Padding(0, 3).
			MarginRight(1),
		ButtonActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("99")).
			Padding(0, 3).
			MarginRight(1),
	}
}

// NewDialog creates a new dialog
func NewDialog(title, message string) *Dialog {
	return &Dialog{
		title:   title,
		message: message,
		width:   60,
		height:  7,
		styles:  DefaultDialogStyles(),
		keyMap:  DefaultDialogKeyMap(),
	}
}

// View renders the dialog
func (d *Dialog) View() string {
	// Create the dialog content
	titleView := d.styles.Title.Render(d.title)
	messageView := d.styles.Message.Render(d.message)
	
	// Create the buttons
	yesButton := d.styles.ButtonActive.Render("Yes")
	noButton := d.styles.Button.Render("No")
	buttonsView := lipgloss.JoinHorizontal(lipgloss.Center, yesButton, noButton)
	
	// Join all content vertically
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleView,
		lipgloss.NewStyle().Height(1).Render(""), // Empty line
		messageView,
		lipgloss.NewStyle().Height(1).Render(""), // Empty line
		buttonsView,
	)
	
	// Apply border to the content
	dialogView := d.styles.Border.Width(d.width).Render(content)
	
	return dialogView
}

// KeyMap defines the key bindings for the dialog
type DialogKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

// DefaultDialogKeyMap returns the default key bindings for the dialog
func DefaultDialogKeyMap() DialogKeyMap {
	return DialogKeyMap{
		Confirm: key.NewBinding(
			key.WithKeys("y", "Y", "enter"),
			key.WithHelp("y/enter", "confirm"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("n", "N", "esc", "q", "ctrl+c"),
			key.WithHelp("n/esc/q", "cancel"),
		),
	}
}