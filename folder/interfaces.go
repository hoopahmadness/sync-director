package folder

import "github.com/hoopahmadness/sync-director/v2/web"

type folderManager[Pairable web.Pairable] interface {
	GetFolderById(id string) (*NetworkFolder[Pairable], bool)
	GetDeviceById(id string) (*Pairable, bool)
}

// type device interface {
// 	QueryFolders() ()
// }
