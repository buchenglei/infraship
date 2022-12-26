package kit

import (
	"errors"
	"sync"
)

type RwLockMap[K comparable, V any] struct {
	_placeholder V
	lock         sync.RWMutex
	saved        map[K]V
}

func NewRwLockMap[K comparable, V any](cap int) *RwLockMap[K, V] {
	return &RwLockMap[K, V]{
		saved: make(map[K]V, cap),
	}
}

func (r *RwLockMap[K, V]) Set(key K, value V) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.saved[key] = value
}

func (r *RwLockMap[K, V]) Get(key K) (V, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	v, ok := r.saved[key]
	if !ok {
		return r._placeholder, errors.New("key not exist")
	}

	return v, nil
}

func (r *RwLockMap[K, V]) Delete(key K) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.saved, key)
}

func (r *RwLockMap[K, V]) Range(f func(K, V) error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for k, v := range r.saved {
		f(k, v)
	}
}

func (r *RwLockMap[K, V]) Count() int {
	return len(r.saved)
}

func (r *RwLockMap[K, V]) Update(key K, updater func(v V) (V, error)) {
	r.lock.Lock()
	defer r.lock.Unlock()

	v, ok := r.saved[key]
	if !ok {
		return
	}
	newV, err := updater(v)
	if err != nil {
		return
	}
	r.saved[key] = newV
}
