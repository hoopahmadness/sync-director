package device

import tea "github.com/charmbracelet/bubbletea"

/**
 Device List View (network)
- m.DeviceList, m.FolderList
- Lists all devices with nickname, status, # connections, # folders. Arrow keys to highlight and scroll, ENTER to select for Device Network View. Perhaps a grid?
- Panel that shows verbose details for selected device
- Allow filtering by status, hidden, pending actions

Single folder, many devices view
- m.FolderList[1]
- Lists all devices syncing this folder with nickname, status, # connections. Arrow keys to highlight and scroll, ENTER to select for Device View. Perhaps a grid?
- Panel that shows verbose details for selected device, including path
- Allow sorting ASC/DEC by name, status, number of connections
- Allow filtering by status, hidden
- Allow syncing this folder to new device

Single Device, many folders View
- m.DeviceList[1], m.GetFoldersByDevice(dev)
- Essentially a recreation of the syncthing web gui
    - Lists various stats and info about device
    - Lists connected devices and folders.
*/

func (dev Device) Init() tea.Cmd {
	return nil
}

func (dev Device) Update(msg tea.Msg) (Device, tea.Cmd) {
	return dev, nil
}

func (dev Device) View() string {
	return ""
	// This is where I want to look at my internal variables and decide what values to expose
	// options such as collapsed/full, a paired folder or not, hidden or not, etc
}
