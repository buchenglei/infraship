package cached

import "sync"

type SyncMap struct {
	store sync.Map
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		store: sync.Map{},
	}
}
