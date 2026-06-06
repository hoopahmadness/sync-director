package web

import "github.com/hoopahmadness/sync-director/v2/biofabric"

// Device pairs are bidirectional objects that show a relationship between two devices
// If one of the devices has offerred to share a folder or connection but a second device has not accepted it,
// the *offering* device will be put in the pending slot.
type Pair[I any, IPtr Pairable[I]] struct {
	Dev1            IPtr
	DevA            IPtr
	OfferPending    IPtr
	DeviceDirection map[IPtr]biofabric.EdgeMode
}

func (dp Pair[I, IPtr]) Other(given IPtr) IPtr {
	if dp.Dev1 == given {
		return dp.DevA
	} else if dp.DevA == given {
		return dp.Dev1
	}
	return nil
}

func (dp Pair[I, IPtr]) GetItems() [2]IPtr {
	items := new([2]IPtr)
	items[0] = dp.Dev1
	items[1] = dp.DevA
	return *items
}

func (dp Pair[I, IPtr]) GetNodeAttributes(node IPtr) (attributes biofabric.EdgeAttributes) {
	if node.Offline() {
		attributes.Mode = biofabric.NEGATIVEEdgeMode
		return
	}

	if dp.OfferPending != nil && dp.OfferPending != node {
		attributes.Mode = biofabric.NULLEdgeMode
		return
	}

	if direction := dp.DeviceDirection[node]; direction != biofabric.NULLEdgeMode {
		attributes.Mode = direction
		return
	}

	attributes.Mode = biofabric.NEUTRALEdgeMode
	return

}
