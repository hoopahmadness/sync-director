package director

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit         key.Binding
	swapDevView  key.Binding
	swapFoldView key.Binding
	devListDown  key.Binding
	devListNet   key.Binding
	devListUp    key.Binding
	foldListDown key.Binding
	foldListNet  key.Binding
	foldListUp   key.Binding
	refresh      key.Binding
	up           key.Binding
	down         key.Binding
	left         key.Binding
	right        key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (sd Director) NetworkDeviceListView() string {
	return ""
	// three panels

	// top panel: Page name & controls legend
	// Q,W,E sort through devices
	// A,S,D sort through folders
	// TAB or whatever switch between network views (when available)
	// left and right arrows switch panel
	// up and down arrows scroll through panels
	// Enter or Space? interact with selection

	// left panel: selectable options
	// - Allow sorting ASC/DEC by name, status, number of connections, num folders
	// - Allow filtering by status, hidden, pending actions
	// - Allow creating new device
	// -

	// right panel: Device blobs vertical list
	// - Lists all devices with nickname, status, # connections, # folders.
	// - ENTER to select for Device Network View.

}
