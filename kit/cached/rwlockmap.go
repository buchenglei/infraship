package cached

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

func (r *RwLockMap[K, V]) Set(key K, value V) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.saved[key] = value
	return nil
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

func (r *RwLockMap[K, V]) Delete(key K) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.saved, key)

	return nil
}

func (r *RwLockMap[K, V]) Range(f func(K, V) error) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var err error
	for k, v := range r.saved {
		if err = f(k, v); err != nil {
			if err == ErrRangeContinue {
				continue
			}
			if err == ErrRangeEnd {
				return nil
			}

			return err
		}
	}

	return nil
}

func (r *RwLockMap[K, V]) Count() int {
	return len(r.saved)
}

func (r *RwLockMap[K, V]) Update(key K, updater func(v V) (V, error)) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	v, ok := r.saved[key]
	if !ok {
		return ErrKeyNotExist
	}
	newV, err := updater(v)
	if err != nil {
		return err
	}
	r.saved[key] = newV
	return nil
}

func (r *RwLockMap[K, V]) Clear() error {
	return nil
}
