package director

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/folder"
	"github.com/hoopahmadness/sync-director/v2/web"
	log "github.com/inconshreveable/log15"
	// flow "github.com/hoopahmadness/flow"
)

type director struct { // rename this to 'app' or something
	DevicesById       map[string]Device
	DeviceConnections *web.DeviceWeb
	NetFoldersById    map[string]NetworkFolder

	// Device List is just a convenient representation of the DevicesByID which can be sorted,
	// iterated over. It starts wil nil, which represents "all devices" or network view. Together
	// with the Device Index it represents part of the current state of the app.
	DeviceList  []*device.Device
	DeviceIndex int

	// Folder List is just a convenient representation of the NetFoldersById which can be sorted,
	// iterated over. It starts with nil, which represents "all folders" or network view. Together
	// with the Folder Index it represents part of the current state of the app.
	FolderList  []*folder.NetworkFolder
	FolderIndex int

	log log.Logger
}

type NewDevicesMsg struct {
	// DevicesById map[string]*Device
}

type NewDevicesConnectionMsg struct {
	// DeviceConnections *DeviceWeb
}

type NewNetFoldersMsg struct {
	// NetFoldersById map[string]*NetworkFolder
}

func (sd *director) getDevicesFromConfig() tea.Msg {
	myDevices := createMyDevices(sd.log)
	byID := map[string]Device{}
	for _, dev := range myDevices {
		// populate my device map
		byID[dev.DeviceId] = dev
	}
	sd.DevicesById = byID
	return NewDevicesMsg{}
}

func (sd *director) refreshConnectedDevices() tea.Msg {
	// start with map of device IDs to Devices
	for _, dev := range sd.DevicesById {
		// get connected devices and register each connection in the web
		connectedDevs, err := dev.GetConnectedDevices(sd)
		if err != nil {
			// fmt.Println("Get ConnectedDevices crashed")
			// panic("Get ConnectedDevices crashed")
		}
		for _, connectedDev := range connectedDevs {
			sd.DeviceConnections.NewDeviceConnection(dev, connectedDev)
		}
	}
	return NewDevicesConnectionMsg{}
}

func (sd *director) refreshDeviceFolders() tea.Msg {
	// go back through each device and get the folders
	for _, thisDevice := range sd.DevicesById {
		folders, _ := folder.GetFolders(thisDevice, &sd.log) // shoop make sure this works after sorting out folder package interfaces
		for _, thisFolder := range folders {
			// turn each folder into a network folder if it isn't already
			netFolder, OK := sd.NetFoldersById[thisFolder.Id]
			if !OK {
				netFolder = folder.NewNetworkFolder[Device](thisFolder.Id, sd.log)
				sd.NetFoldersById[folder.Id] = netFolder
			}
			netFolder.IngestFolder(folder, m)
		}
	}
	return NewNetFoldersMsg{}
}

func initialState(logger log.Logger) model {
	return director{
		DevicesById:       map[string]*device.Device{},
		DeviceConnections: web.NewDeviceWeb(logger),
		NetFoldersById:    map[string]*folder.NetworkFolder{},
		// Title:             "Syncthing Director",
		// Subtitle:          "Devices and Folders",
		// modelPage:         titlePage,
		log: logger,
	}
}

func (sd *director) Init() tea.Cmd {
	return sd.getDevicesFromConfig
}

func (sd *director) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	sd.log.Debug("New Message", "msg", msg)

	switch msg := msg.(type) {

	case NewDevicesMsg:
		// New devices have been found, get all connected devices
		return sd, sd.refreshConnectedDevices

	case NewDevicesConnectionMsg:
		// New connections have been found
		return sd, sd.refreshDeviceFolders

	case NewNetFoldersMsg:
		// New folders have been found
		return sd, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return sd, tea.Quit
		}
	default:
		return sd, nil
	}
	return sd, nil
}

func (sd *director) View() string {
	output := ""  //fmt.Sprintf("%s\n%s\n", sd.Title, sd.Subtitle)
	switch true { // sd.modelPage {
	case true: //titlePage: // just print all the devices and then all the folders
		deviceList := []*device.Device{}
		for _, dev := range sd.DevicesById {
			deviceList = append(deviceList, dev)
		}
		sort.Stable(ByID(deviceList))
		sort.Stable(deviceByStatus(deviceList))
		for _, dev := range deviceList {
			output += fmt.Sprintf("%s (%s)- %s \n", dev.DeviceId, dev.Nickname, dev.Status)
		}
		output += "SHOOP \n"

		folderList := []*folder.NetworkFolder{}
		for _, folder := range sd.NetFoldersById {
			folderList = append(folderList, folder)
		}
		sort.Stable(netFolderByDevices(folderList))
		sort.Stable(netFolderById(folderList))
		output += "\n"
		for _, folder := range folderList {
			output += folder.Id + "\n"
		}
	}
	return output
}

func (sd *director) GetFoldersByDevice(dev *device.Device) map[string]*folder.Folder {
	out := map[string]*folder.Folder{}
	for id, netFold := range sd.NetFoldersById {
		fold, exists := netFold.Folders[dev]
		if exists {
			out[id] = fold
		}
	}
	return out
}
