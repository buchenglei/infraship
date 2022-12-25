package config

import (
	"errors"
	"io"
	"strings"

	"github.com/buchenglei/infraship/skeleton"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ConfigFinder skeleton.Finder[string, any]

var _ ConfigFinder = &Config{}

type Config struct {
	Viper *viper.Viper
}

func New(reader io.ReadCloser) (*Config, error) {
	defer reader.Close()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}

	return &Config{Viper: v}, nil
}

func (c *Config) Sub(prefix string) (*Config, error) {
	if !c.Viper.IsSet(prefix) {
		return nil, errors.New("config key " + prefix + " is not set")
	}

	return &Config{
		Viper: c.Viper.Sub(prefix),
	}, nil
}

func (c *Config) Find(k string) (any, error) {
	v := c.Viper.Get(k)
	if v == nil {
		return nil, skeleton.ErrKeyNotExist
	}

	return v, nil
}

func (c *Config) UnmarshalALL(v interface{}) error {
	return c.Unmarshal("", v)
}

func (c *Config) Unmarshal(prefix string, v interface{}) error {
	var _v *viper.Viper
	if prefix != "" {
		if !c.Viper.IsSet(prefix) {
			return errors.New("config key " + prefix + " is not set")
		}
		_v = c.Viper.Sub(prefix)
	} else {
		_v = c.Viper
	}
	if _v == nil {
		return errors.New("get viper error")
	}

	err := _v.Unmarshal(v, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "toml"
		dc.IgnoreUntaggedFields = true
	})
	if err != nil {
		return err
	}

	return nil
}
