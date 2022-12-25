package connector

import (
	"context"
)

//go:generate mockgen -destination=./connector_mock.go -package=connector -source=./connector.go
type Connector[T any] interface {
	Connect(ctx context.Context) (T, error)
	Close(ctx context.Context) error
}
