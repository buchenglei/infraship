package skeleton

import "errors"

var ErrKeyNotExist = errors.New("key not exist")

type Finder[K, V any] interface {
	Find(K) (V, error)
}

type FindHandler[K, V any] func(K) (V, error)

func (f FindHandler[K, V]) Find(k K) (V, error) {
	return f(k)
}
