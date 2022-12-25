package syncx

import "sync"

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) *WaitGroupWrapper {
	w.Add(1)
	go func() {
		defer w.Done()
		cb()
	}()

	return w
}
