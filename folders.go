package main

type NetworkFolder struct {
	id        string
	folders   map[*Device]Folder
	deviceWeb DeviceWeb
}

type Folder struct {
	id        string
	localPath string
	device    *Device
}
