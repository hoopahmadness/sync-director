package director

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/folder"
	"github.com/hoopahmadness/sync-director/v2/web"
	log "github.com/inconshreveable/log15"
)

type Director struct {
	DevicesById       map[string]*device.Device
	DeviceConnections *web.Web[device.Device]
	NetFoldersById    map[string]*folder.NetworkFolder[device.Device]
	RWLock            sync.RWMutex

	// Device List is just a convenient representation of the DevicesByID which can be sorted,
	// iterated over. It starts with nil, which represents "all devices" or network view. Together
	// with the Device Index it represents part of the current state of the app.
	DeviceList  []*device.Device
	DeviceIndex int

	// Folder List is just a convenient representation of the NetFoldersById which can be sorted,
	// iterated over. It starts with nil, which represents "all folders" or network view. Together
	// with the Folder Index it represents part of the current state of the app.
	FolderList  []*folder.NetworkFolder[device.Device]
	FolderIndex int

	// DeviceOrFolder is just a boolean that helps define part of the state of the app
	// in this case it defines whether we are showing the Network mode
	// for the Device List or the Folder List
	DeviceOverFolder bool

	log log.Logger

	readDevicesFromFile func(log.Logger) []tea.Msg
}

// Helper func that goes through all our devices and attempts to query them for connected devices.
// We create a list of tea commands where each one queries one of our known devices
func (sd *Director) refreshConnectedDevices(incomingMsg MsgHistory) tea.Msg {
	sd.log.Debug("Starting refreshConnectedDevices")
	getConnectedDevicesCommands := tea.BatchMsg{}
	sd.RWLock.RLock()
	for _, dev := range sd.DevicesById {
		// get connected devices and register each connection in the web
		getConnectedDevicesCommands = append(getConnectedDevicesCommands,
			func() tea.Msg {
				scopedDev := dev
				context := fmt.Sprintf("Getting device connections for %s", scopedDev.DeviceId)
				newConnections := map[*device.Device][]*device.Device{}
				connectedDevs, err := scopedDev.GetConnectedDevices(sd)
				if err != nil {
					sd.log.Crit("Get ConnectedDevices crashed; probably offline?")
					return nil
				}
				newConnections[scopedDev] = connectedDevs
				outgoing := NewDevicesConnectionMsg{DeviceConnections: newConnections}
				outgoing.AddHistory(incomingMsg, context)
				return outgoing
			})
	}
	sd.RWLock.RUnlock()
	return getConnectedDevicesCommands
}

func (sd *Director) populateDeviceList() {
	devList := []*device.Device{nil}
	sd.RWLock.RLock()
	for _, dev := range sd.DevicesById {
		devList = append(devList, dev)
	}
	sd.RWLock.RUnlock()
	sort.Stable(DeviceByID(devList))
	sort.Stable(deviceByStatus(devList))
	sd.RWLock.Lock()
	sd.DeviceList = devList

	sd.RWLock.Unlock()
}

func (sd *Director) populateFolderList() {
	foldList := []*folder.NetworkFolder[device.Device]{nil}
	sd.RWLock.RLock()
	for _, netFolder := range sd.NetFoldersById {
		foldList = append(foldList, netFolder)
	}
	sd.RWLock.RUnlock()
	sort.Stable(NetFoldersByID[device.Device](foldList))
	sd.FolderList = foldList
}

func (sd *Director) refreshDeviceFolders(incomingMsg MsgHistory) tea.Msg {
	// go back through each device and get the folders
	getFoldersCmds := tea.BatchMsg{}
	sd.RWLock.RLock()
	for _, dev := range sd.DevicesById {
		scopedDev := dev

		getFoldersCmds = append(getFoldersCmds,
			func() tea.Msg {
				foldersToIngest := map[*folder.Folder]*folder.NetworkFolder[device.Device]{}
				context := fmt.Sprintf("Getting folders for %s", scopedDev.DeviceId)
				folders, err := folder.GetFolders(scopedDev, &sd.log)
				if err != nil {
					sd.log.Warn("Unable to get folders for this device, skipping", "device", scopedDev.DeviceId)
					return nil
				}
				for _, thisFolder := range folders {
					// turn each folder into a network folder if it isn't already
					sd.RWLock.RLock()
					netFolder, OK := sd.NetFoldersById[thisFolder.Id]
					sd.RWLock.RUnlock()
					if !OK {
						netFolder = folder.NewNetworkFolder[device.Device](thisFolder.Id, sd.log)
					}
					foldersToIngest[thisFolder] = netFolder
				}
				outgoingMsg := MsgHistory{}
				outgoingMsg.AddHistory(incomingMsg, context)
				outgoing := NewNetFoldersMsg{
					foldersToIngest: foldersToIngest,
					MsgHistory:      outgoingMsg,
				}
				return outgoing
			})
	}

	sd.RWLock.RUnlock()
	return getFoldersCmds
}

func InitialState(logger log.Logger, createDevicesFunc func(log.Logger) []tea.Msg) Director {
	return Director{
		DevicesById:         map[string]*device.Device{},
		DeviceConnections:   web.NewDeviceWeb[device.Device](logger),
		NetFoldersById:      map[string]*folder.NetworkFolder[device.Device]{},
		log:                 logger,
		readDevicesFromFile: createDevicesFunc,
	}
}

func (sd *Director) Init() tea.Cmd {
	devices := sd.readDevicesFromFile(sd.log)
	returnMsgs := []tea.Cmd{}
	for _, devInfo := range devices {
		thisDevInfo := devInfo.(NewNetworkDeviceMsg)
		context := fmt.Sprintf("new device info added from Init: %s", thisDevInfo.DeviceID)

		cmd := func() tea.Msg {
			thisDevInfo.MsgHistory.History = []string{context}
			return thisDevInfo
		}
		returnMsgs = append(returnMsgs, cmd)
	}
	return tea.Batch(returnMsgs...)
}

func (sd *Director) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	sd.log.Debug("New Message", "msg", msg)
	var outputCmd func(msg MsgHistory) tea.Msg
	var outputMsg MsgHistory

	switch msg := msg.(type) {

	case NewNetworkDeviceMsg:
		// Create device if it doesn't already exist
		// Update existing devices with new info
		var dev *device.Device
		var exists bool
		sd.RWLock.Lock()
		if dev, exists = sd.DevicesById[msg.DeviceID]; !exists {
			dev = device.NewDevice(msg.Nickname, sd.log)
			sd.DevicesById[msg.DeviceID] = dev
		}

		sd.RWLock.Unlock()
		dev.AddNetworkInfo(msg.DeviceID, msg.APIKey, msg.IPAddress)

		sd.populateDeviceList()
		outputMsg.AddHistory(msg.MsgHistory, "refreshing connected devices after new devices found")
		outputCmd = sd.refreshConnectedDevices

	case NewDevicesConnectionMsg:
		// New connections have been found
		newDeviceConnections := msg.DeviceConnections
		for devA, connectedDevices := range newDeviceConnections {
			sd.log.Debug("waiting to lock")
			sd.RWLock.Lock()
			sd.log.Debug("got lock, waiting for loop to end")
			for _, devB := range connectedDevices {
				sd.DevicesById[devB.DeviceId] = devB
				sd.DeviceConnections.NewDeviceConnection(devA, devB)
			}
			sd.RWLock.Unlock()
		}
		sd.populateDeviceList()
		outputMsg.AddHistory(msg.MsgHistory, "refreshing device folders after forming new dev connections")
		outputCmd = sd.refreshDeviceFolders

	case NewNetFoldersMsg:
		// New folders have been found
		newFolders := msg.foldersToIngest
		for folderToIngest, networkFolder := range newFolders {
			sd.RWLock.Lock()
			sd.NetFoldersById[networkFolder.Id] = networkFolder

			sd.RWLock.Unlock()
			networkFolder.IngestFolder(folderToIngest, sd)
		}
		sd.populateFolderList()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return sd, tea.Quit
		case tea.KeyTab:
			sd.ParseInput(SWAPKEY)

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				sd.ParseInput(DEVICELISTDOWNKEY)

			case "w":
				sd.ParseInput(DEVICELISTNETWORKKEY)

			case "e":
				sd.ParseInput(DEVICELISTUPKEY)

			case "a":
				sd.ParseInput(FOLDERLISTDOWNKEY)

			case "s":
				sd.ParseInput(FOLDERLISTNETWORKKEY)

			case "d":
				sd.ParseInput(FOLDERLISTUPKEY)

			case "r":
				outputMsg.History = []string{"refreshing devices from user input"}
				outputCmd = sd.refreshConnectedDevices
			}
		}

	}
	outputWrapper := func() tea.Msg {
		if outputCmd == nil {
			return nil
		}
		return outputCmd(outputMsg)
	}
	return sd, outputWrapper
}

func (sd *Director) View() string {
	output := ""
	status := sd.GetStatus()

	output += status + " " + strconv.FormatInt(int64(sd.DeviceIndex), 10) +
		" " + strconv.FormatInt(int64(sd.FolderIndex), 10) + "\n"
	switch true { // sd.modelPage {
	case true: //titlePage: // just print all the devices and then all the folders
		deviceList := sd.DeviceList
		for _, dev := range deviceList {
			if dev == nil {
				output += "no device\n"
				continue
			}
			output += fmt.Sprintf("%s (%s)- %s \n", dev.DeviceId, dev.Nickname, dev.Status)
		}

		folderList := sd.FolderList
		output += "\n"
		for _, folder := range folderList {
			if folder == nil {
				output += "no folder\n"
				continue
			}
			output += folder.Id + "\n"
		}
	}
	return output
}

func (sd *Director) GetFoldersByDevice(dev *device.Device) map[string]*folder.Folder {
	out := map[string]*folder.Folder{}
	sd.RWLock.RLock()
	for id, netFold := range sd.NetFoldersById {
		fold, exists := netFold.Folders[dev]
		if exists {
			out[id] = fold
		}
	}
	sd.RWLock.RUnlock()
	return out
}

func (sd *Director) GetDeviceById(id string) (*device.Device, bool) {
	sd.RWLock.RLock()
	dev, OK := sd.DevicesById[id]
	sd.RWLock.RUnlock()
	return dev, OK
}

func (sd *Director) SetDeviceById(dev *device.Device) {
	sd.RWLock.Lock()
	sd.DevicesById[dev.DeviceId] = dev

	sd.RWLock.Unlock()
}

func (sd *Director) GetFolderById(id string) (*folder.NetworkFolder[device.Device], bool) {
	sd.RWLock.RLock()
	netFold, OK := sd.NetFoldersById[id]
	sd.RWLock.RUnlock()
	return netFold, OK
}

func (sd *Director) setFolderById(netFolder *folder.NetworkFolder[device.Device]) {
	sd.RWLock.Lock()
	sd.NetFoldersById[netFolder.Id] = netFolder

	sd.RWLock.Unlock()
}
