package models

type WatchStructure struct {
	Path             string `json:"path"`
	WatchRecursively bool   `json:"watchRecursively"`
}
