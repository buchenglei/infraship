package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Name     string
	Host     string
	Port     int
	DB       int
	PoolSize int
	Password string
}

type RedisConnector struct {
	conf RedisConfig

	client *redis.Client
}

func NewRedisConnector(conf RedisConfig) (*RedisConnector, error) {
	return &RedisConnector{conf: conf}, nil
}

func (r *RedisConnector) Ping(ctx context.Context) error {
	if r.client.Ping(ctx).String() != "Pong" {
		return fmt.Errorf("ping redis client error: %s", r.client.Ping(ctx).String())
	}

	return nil
}

func (r *RedisConnector) Connect(ctx context.Context) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.conf.Host, r.conf.Port),
		Password: r.conf.Password,
		DB:       r.conf.DB,
		PoolSize: r.conf.PoolSize,
	})

	// 检查redis连接状态是否正常
	if pong := cli.Ping(context.Background()); pong != nil && (pong.Err() != nil || pong.Val() != "PONG") {
		return nil, pong.Err()
	}

	r.client = cli

	return cli, nil
}

func (r *RedisConnector) Close(ctx context.Context) error {
	return r.client.Close()
}
