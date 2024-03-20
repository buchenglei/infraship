package cached

import "errors"

var (
	ErrKeyNotExist   = errors.New("key is not exist")
	ErrRangeContinue = errors.New("continue")
	ErrRangeEnd      = errors.New("end")
)

type Cached[K comparable, V any] interface {
	Get(K) (V, error)
	Set(K, V) error
	Update(key K, updater func(oldV V) (V, error)) error
	Delete(K) error
	Range(func(K, V) error) error
	Count() int
	Clear() error
}

var (
	_ Cached[uint8, uint8] = &RwLockMap[uint8, uint8]{}
)
