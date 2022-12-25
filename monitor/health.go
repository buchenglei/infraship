package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/buchenglei/infraship/skeleton"
	"github.com/rs/zerolog/log"
)

type healthOptionFunc func(*Health)

func HealthAddModule(name string, m skeleton.Pinger) healthOptionFunc {
	return func(h *Health) {
		h.modules[name] = m
	}
}

func HealthSetAutomaticWatchInterval(i time.Duration) healthOptionFunc {
	if i <= 0 {
		i = time.Minute * 3
	}
	return func(h *Health) {
		h.interval = i
	}
}

func HealthSetMonitorCap(cap int) healthOptionFunc {
	return func(h *Health) {
		h.cap = cap
	}
}

func HealthSetErrorHandleFunc(f func(t time.Time, errs []error)) healthOptionFunc {
	return func(h *Health) {
		h.errsHandleFunc = f
	}
}

type healthResult struct {
	t    time.Time
	errs []error
}

type Health struct {
	interval       time.Duration
	cap            int
	errsHandleFunc func(t time.Time, errs []error)

	modules    map[string]skeleton.Pinger
	resultChan chan healthResult
}

func NewHealth(opts ...healthOptionFunc) *Health {
	h := &Health{
		interval: 3 * time.Minute,
		cap:      10,
		errsHandleFunc: func(t time.Time, errs []error) {
			log.Warn().Time("report", t).Errs("result", errs).Msg("module health check result")
		},
	}

	for _, opt := range opts {
		opt(h)
	}

	h.modules = make(map[string]skeleton.Pinger, h.cap)
	h.resultChan = make(chan healthResult, h.cap*2)
	go h.handleResult()

	return h
}

func (h *Health) Monitor() {
	// 第一次启动立刻检查一次
	ctx := context.Background()
	errs := h.check(ctx)
	if len(errs) > 0 {
		h.resultChan <- healthResult{
			t:    time.Now(),
			errs: errs,
		}
	}

	go func() {
		ticker := time.NewTicker(h.interval)
		for range ticker.C {
			errs = h.check(ctx)
			if len(errs) > 0 {
				h.resultChan <- healthResult{
					t:    time.Now(),
					errs: errs,
				}
			}
		}
	}()
}

func (h *Health) check(ctx context.Context) []error {
	var errs []error = make([]error, 0, h.cap)
	for name, module := range h.modules {
		err := module.Ping(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("[%s] is bad heath: %w", name, err))
		}
	}

	return errs
}

func (h *Health) handleResult() {
	for res := range h.resultChan {
		h.errsHandleFunc(res.t, res.errs)
	}
}
