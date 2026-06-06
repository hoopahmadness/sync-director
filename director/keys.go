package director

import (
	"github.com/charmbracelet/bubbles/key"
)

func newKeyMap() keyMap {
	var keys = keyMap{
		quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("ESC", "quit"),
		),
		swapDevView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("TAB", "Devices View"),
		),
		swapFoldView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("TAB", "Folders View"),
		),
		devListDown: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Previous Device"),
		),
		devListUp: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "Next Device"),
		),
		devListNet: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "All Devices"),
		),
		foldListDown: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Previous Folder"),
		),
		foldListUp: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "Next Folder"),
		),
		foldListNet: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "All Folders"),
		),
		up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
		),
	}
	return keys

}
