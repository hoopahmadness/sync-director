package folder

import "github.com/hoopahmadness/sync-director/v2/web"

type folderManager[P any, PPtr web.Pairable[P]] interface {
	GetFolderById(id string) (*NetworkFolder[P, PPtr], bool)
	GetDeviceById(id string) (PPtr, bool)
}
