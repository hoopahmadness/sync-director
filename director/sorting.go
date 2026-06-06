package director

import (
	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/folder"
	"github.com/hoopahmadness/sync-director/v2/web"
)

type DeviceByID []*device.Device

func (byID DeviceByID) Len() int {
	return len(byID)
}
func (byID DeviceByID) Less(i, j int) bool {
	if byID[i] == nil {
		return true
	} else if byID[j] == nil {
		return false
	}
	return byID[i].GetId() < byID[j].GetId()
}
func (byID DeviceByID) Swap(i, j int) {
	byID[i], byID[j] = byID[j], byID[i]
}

type deviceByStatus []*device.Device

func (dStat deviceByStatus) Len() int {
	return len(dStat)
}
func (dStat deviceByStatus) Less(i, j int) bool {
	if dStat[i] == nil {
		return true
	} else if dStat[j] == nil {
		return false
	}
	return dStat[i].GetOrderedStatus() < dStat[j].GetOrderedStatus()
}
func (dStat deviceByStatus) Swap(i, j int) {
	dStat[i], dStat[j] = dStat[j], dStat[i]
}

type NetFoldersByID[T any, TPtr web.Pairable[T]] []*folder.NetworkFolder[T, TPtr]

func (nfID NetFoldersByID[T, TPtr]) Len() int {
	return len(nfID)
}
func (nfID NetFoldersByID[T, TPtr]) Less(i, j int) bool {
	if nfID[i] == nil {
		return true
	} else if nfID[j] == nil {
		return false
	}
	return nfID[i].GetId() < nfID[j].GetId()
}
func (nfID NetFoldersByID[T, TPtr]) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}

type netFolderByPairables[T any, TPtr web.Pairable[T]] []*folder.NetworkFolder[T, TPtr]

func (nfID netFolderByPairables[T, TPtr]) Len() int {
	return len(nfID)
}
func (nfID netFolderByPairables[T, TPtr]) Less(i, j int) bool {
	if nfID[i] == nil {
		return true
	} else if nfID[j] == nil {
		return false
	}
	return len(nfID[i].Folders) < len(nfID[j].Folders)
}
func (nfID netFolderByPairables[T, TPtr]) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}
