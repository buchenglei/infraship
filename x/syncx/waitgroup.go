package syncx

import "sync"

type WaitGroupWrapper struct {
	wg sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) *WaitGroupWrapper {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		cb()
	}()

	return w
}

func (w *WaitGroupWrapper) Wait() {
	w.wg.Wait()
}
