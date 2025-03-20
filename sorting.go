package main

type deviceByID []*Device

func (dID deviceByID) Len() int {
	return len(dID)
}
func (dID deviceByID) Less(i, j int) bool {
	return dID[i].DeviceId < dID[j].DeviceId
}
func (dID deviceByID) Swap(i, j int) {
	dID[i], dID[j] = dID[j], dID[i]
}

type deviceByStatus []*Device

func (dStat deviceByStatus) Len() int {
	return len(dStat)
}
func (dStat deviceByStatus) Less(i, j int) bool {
	return orderedStatuses[dStat[i].Status] < orderedStatuses[dStat[j].Status]
}
func (dStat deviceByStatus) Swap(i, j int) {
	dStat[i], dStat[j] = dStat[j], dStat[i]
}

type netFolderById []*NetworkFolder

func (nfID netFolderById) Len() int {
	return len(nfID)
}
func (nfID netFolderById) Less(i, j int) bool {
	return nfID[i].Id < nfID[j].Id
}
func (nfID netFolderById) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}

type netFolderByDevices []*NetworkFolder

func (nfID netFolderByDevices) Len() int {
	return len(nfID)
}
func (nfID netFolderByDevices) Less(i, j int) bool {
	return len(nfID[i].folders) < len(nfID[j].folders)
}
func (nfID netFolderByDevices) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}
