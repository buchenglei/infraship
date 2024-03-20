package config

import (
	"errors"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

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

func (c *Config) GetString(k string, defaults ...string) string {
	if !c.Viper.InConfig(k) && len(defaults) == 1 {
		return defaults[0]
	}

	return c.Viper.GetString(k)
}

func (c *Config) GetInt(k string, defaults ...int) int {
	if !c.Viper.InConfig(k) && len(defaults) == 1 {
		return defaults[0]
	}

	return c.Viper.GetInt(k)
}

func (c *Config) GetBool(k string, defaults ...bool) bool {
	if !c.Viper.InConfig(k) && len(defaults) == 1 {
		return defaults[0]
	}

	return c.Viper.GetBool(k)
}
