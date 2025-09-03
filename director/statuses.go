package director

const (
	// stages
	DEVICELISTNETWORKPAGE = "Device List for Network"
	FOLDERLISTNETWORKPAGE = "Folder List for Network"
	DEVICELISTFOLDERPAGE  = "Device List for Folder"
	FOLDERLISTDEVICEPAGE  = "Folder List for Device"
	DEVICEFOLDERPAGE      = "Device and Folder Page"

	// actions
	DEVICELISTUPKEY      = "Device List Increase"
	DEVICELISTDOWNKEY    = "Device List Decrease"
	DEVICELISTNETWORKKEY = "Device List Return to Zero"
	FOLDERLISTUPKEY      = "Folder List Increase"
	FOLDERLISTDOWNKEY    = "Folder List Decrease"
	FOLDERLISTNETWORKKEY = "Folder List Return to Zero"
	SWAPKEY              = "Swap View"
)

func (sd *Director) GetStatus() string {
	zeroDevice := sd.DeviceIndex == 0

	zeroFolder := sd.FolderIndex == 0

	switch {
	case zeroDevice && zeroFolder:
		if sd.DeviceOverFolder {
			return DEVICELISTNETWORKPAGE
		} else {
			return FOLDERLISTNETWORKPAGE
		}
	case zeroDevice:
		return DEVICELISTFOLDERPAGE
	case zeroFolder:
		return FOLDERLISTDEVICEPAGE
	default:
		return DEVICEFOLDERPAGE
	}
}

func (sd *Director) ParseInput(action string) {
	zeroDevice := sd.DeviceIndex == 0

	zeroFolder := sd.FolderIndex == 0
	oldStatus := sd.GetStatus()

	switch action {
	case FOLDERLISTUPKEY:
		sd.FolderIndex += 1
		if sd.FolderIndex >= len(sd.FolderList) {
			sd.FolderIndex = 0
		}
	case FOLDERLISTDOWNKEY:
		sd.FolderIndex += -1
		if sd.FolderIndex < 0 {
			sd.FolderIndex = len(sd.FolderList) - 1
		}
	case DEVICELISTNETWORKKEY:
		sd.FolderIndex = 0
	case DEVICELISTUPKEY:
		sd.DeviceIndex += 1
		if sd.DeviceIndex >= len(sd.DeviceList) {
			sd.DeviceIndex = 0
		}
	case DEVICELISTDOWNKEY:
		sd.DeviceIndex += -1
		if sd.DeviceIndex < 0 {
			sd.DeviceIndex = len(sd.DeviceList) - 1
		}
	case FOLDERLISTNETWORKKEY:
		sd.DeviceIndex = 0
	case SWAPKEY:
		if zeroDevice && zeroFolder {
			sd.DeviceOverFolder = !sd.DeviceOverFolder
		}
	}
	newStatus := sd.GetStatus()
	sd.log.Debug("Input parsed.", "new status", newStatus, "old status", oldStatus)

}
