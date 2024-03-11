package models


type WatchStructure struct {
	Path string `json:"path"`
	IsFolder bool  `json:"isFolder"`
	WatchRecursively bool `json:"watchRecursively"`
}