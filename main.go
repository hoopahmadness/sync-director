package main

import (
	"fmt"
)

var devicesById = map[string]*Device{}
var netFoldersById = map[string]*NetworkFolder{}
var webManager = &WebManager{}

func main() {
	webManager = newWebManager()
	// start with map of device IDs to Devices
	myDevices := createMyDevices()
	for _, dev := range myDevices {
		devicesById[dev.deviceId] = dev
	}
	folders, _ := devicesById["6COWOOR-PAATSRB-I36VNI7-SUM2RR5-ONENLNB-P7QHO4X-P2QIC22-RGII2AV"].GetFolders()
	for _, folder := range folders {
		netFolder, OK := netFoldersById[folder.Id]
		if !OK {
			netFolder = newNetworkFolder(folder.Id)
			netFoldersById[folder.Id] = netFolder
		}
		netFolder.IngestFolder(folder)
	}
	// ask each device for its folders. Create a device web for each folder.
	fmt.Println(netFoldersById["default"].)
	// fmt.Println(webManager.master.allPairs)
}
