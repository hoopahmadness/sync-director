package main

import "time"

type Device struct { // Do I even need this?
	id       string
	lastSeen time.Time
	lastIP   string
	apiKey   string
}

type DevicePair struct {
	dev1 *Device
	devA *Device
}

func (dp *DevicePair) Other(given *Device) *Device {
	if dp.dev1 == given {
		return dp.devA
	} else if dp.devA == given {
		return dp.dev1
	}
	return nil
}

// A device web represents a set of links between pairs of devices.
type DeviceWeb struct {
	// allPairs  []*DevicePair
	perDevice map[*Device][]*DevicePair
}

func newDeviceWeb() DeviceWeb {
	dw := DeviceWeb{}
	// dw.allPairs = []*DevicePair{}
	dw.perDevice = map[*Device][]*DevicePair{}
	return dw
}

func (dw *DeviceWeb) getOtherDevices(aDevice *Device) []*Device {
	dw.initializeDevicePairsListIfNeeded(aDevice)
	pairs := dw.perDevice[aDevice]
	devices := []*Device{}
	for _, pair := range pairs {
		devices = append(devices, pair.Other(aDevice))
	}
	return devices
}

func (dw *DeviceWeb) getPairing(aDevice, anotherDevice *Device) *DevicePair {
	pairList := dw.perDevice[aDevice]
	for _, pair := range pairList {
		if pair.Other(aDevice) == anotherDevice {
			return pair
		}
	}
	return nil
}

func (dw *DeviceWeb) addPair(aDevice, anotherDevice *Device) {
	dp := &DevicePair{aDevice, anotherDevice}
	// dw.allPairs = append(dw.allPairs, dp)
	dw.initializeDevicePairsListIfNeeded(aDevice, anotherDevice)
	listforADevice := dw.perDevice[aDevice]
	listforADevice = append(listforADevice, dp)
	dw.perDevice[aDevice] = listforADevice

	listforAnotherDevice := dw.perDevice[anotherDevice]
	listforAnotherDevice = append(listforAnotherDevice, dp)
	dw.perDevice[anotherDevice] = listforAnotherDevice
}

func (dw *DeviceWeb) initializeDevicePairsListIfNeeded(devices ...*Device) {
	for _, aDevice := range devices {
		_, OK := dw.perDevice[aDevice]
		if !OK {
			dw.perDevice[aDevice] = []*DevicePair{}
		}
	}
}
