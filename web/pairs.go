package web

// Device pairs are bidirectional objects that show a relationship between two devices
// If one of the devices has offerred to share a folder or connection but a second device has not accepted it,
// the *offering* device will be put in the pending slot.
type Pair[Device interface{}] struct {
	Dev1         *Device
	DevA         *Device
	OfferPending *Device
}

func (dp *Pair[Device]) Other(given *Device) *Device {
	if dp.Dev1 == given {
		return dp.DevA
	} else if dp.DevA == given {
		return dp.Dev1
	}
	return nil
}

// If one of the devices has not accepted the folder then this returns the
// device *offering* the folder for syncing. If both hosts are sharing then returns nil
func (dp *Pair[Device]) GetPending() *Device {
	return dp.OfferPending
}
