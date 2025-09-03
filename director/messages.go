package director

import (
	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/folder"
)

type MsgHistory struct {
	History []string
}

func (h *MsgHistory) AddHistory(oldMessage MsgHistory, newHistory string) {
	h.History = append(oldMessage.History, newHistory)
}

// Message with the information to introduce a new device to the manager
type NewNetworkDeviceMsg struct {
	Nickname  string
	DeviceID  string
	APIKey    string
	IPAddress string
	MsgHistory
}

// Message representing new device connections that can be added to the current device connections
// We return a map of multiple connected devices per device
type NewDevicesConnectionMsg struct {
	DeviceConnections map[*device.Device][]*device.Device
	MsgHistory
}

// Message representing new or updated Network Folders
// We return a map of a new folder to the network folder that will ingest it
type NewNetFoldersMsg struct {
	foldersToIngest map[*folder.Folder]*folder.NetworkFolder[device.Device]
	MsgHistory
}
