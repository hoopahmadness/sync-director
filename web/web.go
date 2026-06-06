package web

import (
	"fmt"

	"github.com/hoopahmadness/sync-director/v2/biofabric"
	log "github.com/inconshreveable/log15"
)

// A device web represents a set of links between pairs of devices in the context of
// a specific folder or all folders
type Web[D any, DPtr Pairable[D]] struct {
	AllPairs  map[*Pair[D, DPtr]]bool
	PerDevice map[DPtr][]*Pair[D, DPtr]
	log       log.Logger // todo remove this and start passing children loggers around instead
}

func (dw *Web[D, Device]) String() string {
	webStr := "Device web:\n"
	for pairing := range dw.AllPairs {
		firstOffering := ""
		secondOffering := ""
		if pairing.OfferPending == pairing.Dev1 {
			firstOffering = "*"
		}
		if pairing.OfferPending == pairing.DevA {
			secondOffering = "*"
		}
		x := pairing.Dev1
		y := pairing.DevA
		pairStr := fmt.Sprintf("%s%s paired with %s%s\n", x.Name(), firstOffering, y.Name(), secondOffering)
		webStr += pairStr
	}
	return webStr
}

// Given a device, find the other devices that are sharing the folder with it
// Does not care if a device is pending or not.
func (dw *Web[D, Device]) getOtherDevices(aDevice Device) []Device {
	dw.initializeDevicePairsListIfNeeded(aDevice)
	pairs := dw.PerDevice[aDevice]
	devices := []Device{}
	for _, pair := range pairs {
		devices = append(devices, pair.Other(aDevice))
	}
	return devices
}

// A folder can add a new pairing between two of its devices.
// In general pairs are bidirectional so it doesn't matter which device is the host;
// adding the inverse pair will be compared to the existing pair and is a no-op.
// If the host device has not accepted syncing for a folder from the synced device, set 'pending'
// to true and the device pair will show that there is a pending folder.
func (dw *Web[D, Device]) NewDevicePairForFolder(hostDevice, syncedDevice Device, pending bool, syncDirections map[Device]biofabric.EdgeMode) {
	dp, err := dw.getDevicePair(hostDevice, syncedDevice)
	if err != nil {
		return
	}
	if pending {
		dp.OfferPending = syncedDevice
	}
	for dev, dir := range syncDirections {
		dp.DeviceDirection[dev] = dir
	}
}

// Create a new connection between devices. All device pairs created this way are pending by default.
// When the inverse pair is created (in other words, when the synced device tries to add a connection
// to this device) then the connection will no longer be pending
func (dw *Web[D, Device]) NewDeviceConnection(hostDevice, connectedDevice Device) {
	dp, err := dw.getDevicePair(hostDevice, connectedDevice)
	if err != nil {
		return
	}
	switch dp.OfferPending {
	case nil:
		dp.OfferPending = hostDevice
	case connectedDevice:
		dp.OfferPending = nil
	}
}

// Adds a new DevicePair to the web.
// Each DevicePair is kept in a map as well as kept in a list for each related device
// (that's three pointers per pair)
// Adding an existing pair again is a no-op
// For now let's assume this is only run internally by the web object
func (dw *Web[D, Device]) addPairing(pair *Pair[D, Device]) {
	if _, OK := dw.AllPairs[pair]; OK {
		return
	}
	aDevice := pair.Dev1
	anotherDevice := pair.DevA
	dw.initializeDevicePairsListIfNeeded(aDevice, anotherDevice)
	listforADevice := dw.PerDevice[aDevice]
	listforADevice = append(listforADevice, pair)
	dw.PerDevice[aDevice] = listforADevice

	listforAnotherDevice := dw.PerDevice[anotherDevice]
	listforAnotherDevice = append(listforAnotherDevice, pair)
	dw.PerDevice[anotherDevice] = listforAnotherDevice

	dw.AllPairs[pair] = true
}

func (dw *Web[D, Device]) initializeDevicePairsListIfNeeded(devices ...Device) {
	for _, aDevice := range devices {
		_, OK := dw.PerDevice[aDevice]
		if !OK {
			dw.PerDevice[aDevice] = []*Pair[D, Device]{}
		}
	}
}

// This func should be used by a DeviceWeb to get a pointer to an existing DevicePair
// If the DevicePair for these devices can't be found it is created.
// This allows us to guarantee that any pair of devices regardless of order will yield the same object
// Returns error when attempting to pair a device to itself
func (dw *Web[D, Device]) getDevicePair(aDevice, anotherDevice Device) (*Pair[D, Device], error) {
	// make sure devices are not the same
	if aDevice == anotherDevice {
		return nil, fmt.Errorf("cannot pair two devices that are the same")
	}

	// see if the pair already exists
	pairList := dw.PerDevice[aDevice]
	for _, pair := range pairList {
		if pair.Other(aDevice) == anotherDevice {
			return pair, nil
		}
	}

	// create it
	dp := &Pair[D, Device]{
		Dev1:         aDevice,
		DevA:         anotherDevice,
		OfferPending: nil,
	}

	dw.addPairing(dp)

	return dp, nil
}

func (dw *Web[D, Device]) GetAllPairs() map[biofabric.Pair[D, Device]]bool {
	iMap := map[biofabric.Pair[D, Device]]bool{}
	for k, v := range dw.AllPairs {
		iMap[k] = v
	}
	return iMap
}
func (dw *Web[D, Device]) GetRelatedPairs(dev Device) []biofabric.Pair[D, Device] {
	iArr := []biofabric.Pair[D, Device]{}
	for _, relatedDev := range dw.PerDevice[dev] {
		iArr = append(iArr, relatedDev)
	}
	return iArr
}

func NewDeviceWeb[Device any, DevicePtr Pairable[Device]](log log.Logger) *Web[Device, DevicePtr] {
	newWeb := &Web[Device, DevicePtr]{}
	newWeb.PerDevice = map[DevicePtr][]*Pair[Device, DevicePtr]{}
	newWeb.AllPairs = map[*Pair[Device, DevicePtr]]bool{}
	newWeb.log = log
	return newWeb
}
