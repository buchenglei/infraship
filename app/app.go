package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Runner interface {
	Name() string
	Run(context.Context) error
	Reload(context.Context) error
	Exit(context.Context) error
}

type shutdownReason = uint8

const (
	internalError shutdownReason = iota + 1
	allRunnerExit
	runnerReload
	rcvSysStopSignal
)

type Application struct {
	name string
	ver  string

	reload       chan struct{}
	shutdownChan chan shutdownReason

	runners []Runner
}

func New(name, version string, opts ...ApplicationOption) (*Application, error) {
	app := &Application{
		name:         name,
		ver:          version,
		reload:       make(chan struct{}),
		shutdownChan: make(chan shutdownReason, 1),
		runners:      make([]Runner, 0, 5),
	}

	var err error
	if len(opts) > 0 {
		for _, opt := range opts {
			err = opt(app)
			if err != nil {
				return nil, err
			}
		}
	}

	go app.watchSysSignal()

	return app, nil
}

func (c *Application) Name() string {
	return c.name
}

func (c *Application) WithRunner(r ...Runner) *Application {
	c.runners = append(c.runners, r...)
	return c
}

func (a *Application) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, runner := range a.runners {
		wg.Add(1)
		go func(r Runner) {
			defer wg.Done()
			err := r.Run(ctx)
			if err != nil {
				log.Error().Err(err).Msg("runner failed")
				a.shutdown(internalError)
				return
			}
			log.Info().Str("runner", r.Name()).Msg("runner Running...")
		}(runner)
	}

	go func() {
		wg.Wait()
		a.shutdown(allRunnerExit)
	}()

	return nil
}

func (a *Application) Wait(ctx context.Context) {
	reason := <-a.shutdownChan

	log.Warn().Uint8("exit reason", reason).Msg("recv shutdown signal")

	var err error
	for _, r := range a.runners {
		err = r.Exit(ctx)
		if err != nil {
			log.Error().Str("runner", r.Name()).Err(err).Msg("runner exit failed")
		} else {
			log.Info().Str("runner", r.Name()).Msg("runner exit success")
		}
	}

	// final exit & let sys know why
	os.Exit(int(reason))
}

func (a *Application) shutdown(t shutdownReason) {
	a.shutdownChan <- t
}

func (a *Application) watchSysSignal() {
	//创建监听退出chan
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	a.shutdown(rcvSysStopSignal)
}
