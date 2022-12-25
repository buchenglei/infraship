package web

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

type StdHttpConnectorConfig struct{}

type StdHttpConnector struct {
	once   sync.Once
	conf   StdHttpConnectorConfig
	client *http.Client
}

func NewStdHttpConnector(conf StdHttpConnectorConfig) *StdHttpConnector {
	return &StdHttpConnector{conf: conf}
}

func (s *StdHttpConnector) Ping(ctx context.Context) error {
	if s.client == nil {
		return errors.New("http client is empty")
	}

	resp, err := s.client.Head("https://www.baidu.com")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("network bad")
	}

	return nil
}

func (s *StdHttpConnector) Connect(ctx context.Context) (*http.Client, error) {
	var err error
	s.once.Do(func() {
		s.client = &http.Client{}
	})
	if err != nil {
		return nil, err
	}

	return s.client, nil
}

func (s *StdHttpConnector) Close(ctx context.Context) error {
	return nil
}
