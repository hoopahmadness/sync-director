package main

import (
	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/folder"
	"github.com/hoopahmadness/sync-director/v2/web"
)

type deviceByID []*device.Device

func (dID deviceByID) Len() int {
	return len(dID)
}
func (dID deviceByID) Less(i, j int) bool {
	return dID[i].DeviceId < dID[j].DeviceId
}
func (dID deviceByID) Swap(i, j int) {
	dID[i], dID[j] = dID[j], dID[i]
}

type deviceByStatus []*device.Device

func (dStat deviceByStatus) Len() int {
	return len(dStat)
}
func (dStat deviceByStatus) Less(i, j int) bool {
	return device.OrderedStatuses[dStat[i].Status] < device.OrderedStatuses[dStat[j].Status]
}
func (dStat deviceByStatus) Swap(i, j int) {
	dStat[i], dStat[j] = dStat[j], dStat[i]
}

type netFolderById[T web.Pairable] []*folder.NetworkFolder[T]

func (nfID netFolderById[T]) Len() int {
	return len(nfID)
}
func (nfID netFolderById[T]) Less(i, j int) bool {
	return nfID[i].Id < nfID[j].Id
}
func (nfID netFolderById[T]) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}

type netFolderByDevices[T web.Pairable] []*folder.NetworkFolder[T]

func (nfID netFolderByDevices[T]) Len() int {
	return len(nfID)
}
func (nfID netFolderByDevices[T]) Less(i, j int) bool {
	return len(nfID[i].Folders) < len(nfID[j].Folders)
}
func (nfID netFolderByDevices[T]) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}
