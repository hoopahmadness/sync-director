package director

type ByID []Identifier

func (byID ByID) Len() int {
	return len(byID)
}
func (byID ByID) Less(i, j int) bool {
	return byID[i].GetId() < byID[j].GetId()
}
func (byID ByID) Swap(i, j int) {
	byID[i], byID[j] = byID[j], byID[i]
}

type deviceByStatus []Device

func (dStat deviceByStatus) Len() int {
	return len(dStat)
}
func (dStat deviceByStatus) Less(i, j int) bool {
	return dStat[i].GetOrderedStatus() < dStat[j].GetOrderedStatus()
}
func (dStat deviceByStatus) Swap(i, j int) {
	dStat[i], dStat[j] = dStat[j], dStat[i]
}

type netFolderByDevices []Folder

func (nfID netFolderByDevices) Len() int {
	return len(nfID)
}
func (nfID netFolderByDevices) Less(i, j int) bool {
	return nfID[i].GetNumSharedDevices() < nfID[j].GetNumSharedDevices()
}
func (nfID netFolderByDevices) Swap(i, j int) {
	nfID[i], nfID[j] = nfID[j], nfID[i]
}
