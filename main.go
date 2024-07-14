package main

import (
	"fmt"
)

var devicesById = map[string]*Device{}
var deviceConnections *DeviceWeb
var netFoldersById = map[string]*NetworkFolder{}

// var webManager = &WebManager{}

func main() {
	// start with map of device IDs to Devices
	deviceConnections = newDeviceWeb()
	myDevices := createMyDevices()
	someClient := *myDevices[1].Client
	fmt.Println(someClient)
	for _, dev := range myDevices {
		// populate my device map
		devicesById[dev.deviceId] = dev

		// get connected devices and register each connection in the web
		connectedDevs, err := dev.GetConnectedDevices()
		if err != nil {
			fmt.Println("Get ConnectedDevices crashed")
			panic("Get ConnectedDevices crashed")
		}
		for _, connectedDev := range connectedDevs {
			deviceConnections.NewDeviceConnection(dev, connectedDev)
		}
	}
	// go back through each device and get the folders
	for _, device := range devicesById {
		folders, _ := device.GetFolders()
		for _, folder := range folders {
			// turn each folder into a network folder if it isn't already
			netFolder, OK := netFoldersById[folder.Id]
			if !OK {
				netFolder = newNetworkFolder(folder.Id)
				netFoldersById[folder.Id] = netFolder
			}
			netFolder.IngestFolder(folder)
		}
	}
	fmt.Println(devicesById)
	fmt.Println(deviceConnections)
	fmt.Println(netFoldersById)
}
