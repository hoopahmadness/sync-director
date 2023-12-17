package main

type Device struct {
	*Client
	id string
	// lastSeen time.Time
	// lastIP   string
	// apiKey   string
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
