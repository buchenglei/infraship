package app

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFeatureNotAllowExec = errors.New("feature not allow exectue")
)

type FeatureCondition func() error

type Feature interface {
	Name() string
	TurnOn()
	TurnOff()
	Execute(ctx context.Context, feat func() error, condFuncs ...FeatureCondition) error
}

type simpleFeature struct {
	name       string
	off        bool
	conditions []FeatureCondition
}

func NewSimpleFeature(name string, preCondFuncs ...FeatureCondition) Feature {
	return &simpleFeature{
		name:       name,
		conditions: preCondFuncs,
	}
}

func (f *simpleFeature) Name() string {
	return f.name
}

func (f *simpleFeature) TurnOn() {
	f.off = false
}

func (f *simpleFeature) TurnOff() {
	f.off = true
}

func (f *simpleFeature) Execute(ctx context.Context, feat func() error, condFuncs ...FeatureCondition) error {
	if f.off {
		return fmt.Errorf("%w: forbid by mannual", ErrFeatureNotAllowExec)
	}

	var err error
	if len(f.conditions) > 0 {
		for _, cond := range f.conditions {
			if err = cond(); err != nil {
				return fmt.Errorf("%w: %s", ErrFeatureNotAllowExec, err.Error())
			}
		}
	}
	if len(condFuncs) > 0 {
		for _, cond := range condFuncs {
			if err = cond(); err != nil {
				return fmt.Errorf("%w: %s", ErrFeatureNotAllowExec, err.Error())
			}
		}
	}

	return feat()
}
