package director

// import (flow "github.com/hoopahmadness/flow")

const (
	// stages
	DEVICELISTNETWORKPAGE = "Device List for Network"
	FOLDERLISTNETWORKPAGE = "Folder List for Network"
	DEVICELISTFOLDERPAGE  = "Device List for Folder"
	FOLDERLISTDEVICEPAGE  = "Folder List for Device"
	DEVICEFOLDERPAGE      = "Device and Folder Page"

	DEVICELISTUPKEY      = "Device List Increase"
	DEVICELISTDOWNKEY    = "Device List Decrease"
	DEVICELISTNETWORKKEY = "Device List Return to Network"
	FOLDERLISTUPKEY      = "Folder List Increase"
	FOLDERLISTDOWNKEY    = "Folder List Decrease"
	FOLDERLISTNETWORKKEY = "Folder List Return to Network"
	SPACEBAR             = "Space bar"
)

func (m *model) GetStatus() (string, error) {

}
