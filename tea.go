package main

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/inconshreveable/log15"
)

type Page string

const titlePage = Page("titlepage")

type model struct {
	DevicesById       map[string]*Device
	DeviceConnections *DeviceWeb
	NetFoldersById    map[string]*NetworkFolder

	Title     string
	Subtitle  string
	modelPage Page
	log       log.Logger
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

func (m *model) getDevicesFromConfig() tea.Msg {
	myDevices := createMyDevices(m.log)
	byID := map[string]*Device{}
	for _, dev := range myDevices {
		// populate my device map
		byID[dev.DeviceId] = dev
	}
	m.DevicesById = byID
	return NewDevicesMsg{}
}

func (m *model) refreshConnectedDevices() tea.Msg {
	// start with map of device IDs to Devices
	for _, dev := range m.DevicesById {
		// get connected devices and register each connection in the web
		connectedDevs, err := dev.GetConnectedDevices(m)
		if err != nil {
			// fmt.Println("Get ConnectedDevices crashed")
			// panic("Get ConnectedDevices crashed")
		}
		for _, connectedDev := range connectedDevs {
			m.DeviceConnections.NewDeviceConnection(dev, connectedDev)
		}
	}
	return NewDevicesConnectionMsg{}
}

func (m *model) refreshDeviceFolders() tea.Msg {
	// go back through each device and get the folders
	for _, device := range m.DevicesById {
		folders, _ := device.GetFolders()
		for _, folder := range folders {
			// turn each folder into a network folder if it isn't already
			netFolder, OK := m.NetFoldersById[folder.Id]
			if !OK {
				netFolder = newNetworkFolder(folder.Id, m.log)
				m.NetFoldersById[folder.Id] = netFolder
			}
			netFolder.IngestFolder(folder, m)
		}
	}
	return NewNetFoldersMsg{}
}

func initialState(logger log.Logger) model {
	return model{
		DevicesById:       map[string]*Device{},
		DeviceConnections: newDeviceWeb(logger),
		NetFoldersById:    map[string]*NetworkFolder{},
		Title:             "Syncthing Director",
		Subtitle:          "Devices and Folders",
		modelPage:         titlePage,
		log:               logger,
	}
}

func (m *model) Init() tea.Cmd {
	return m.getDevicesFromConfig
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.log.Debug("New Message", "msg", msg)

	switch msg := msg.(type) {

	case NewDevicesMsg:
		// New devices have been found, get all connected devices
		return m, m.refreshConnectedDevices

	case NewDevicesConnectionMsg:
		// New connections have been found
		return m, m.refreshDeviceFolders

	case NewNetFoldersMsg:
		// New folders have been found
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m *model) View() string {
	output := fmt.Sprintf("%s\n%s\n", m.Title, m.Subtitle)
	switch m.modelPage {
	case titlePage: // just print all the devices and then all the folders
		deviceList := []*Device{}
		for _, dev := range m.DevicesById {
			deviceList = append(deviceList, dev)
		}
		sort.Stable(deviceByID(deviceList))
		sort.Stable(deviceByStatus(deviceList))
		for _, dev := range deviceList {
			output += fmt.Sprintf("%s (%s)- %s \n", dev.DeviceId, dev.Nickname, dev.Status)
		}
		output += "SHOOP \n"

		folderList := []*NetworkFolder{}
		for _, folder := range m.NetFoldersById {
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
