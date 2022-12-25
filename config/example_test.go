package config

import (
	"testing"

	"github.com/buchenglei/infraship/skeleton"
	"github.com/stretchr/testify/assert"
)

type app struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Maintainers []string `toml:"maintainers"`
	Env         string   `toml:"env"`
}

func TestConfigNew(t *testing.T) {
	assert := assert.New(t)

	confReader, err := NewFileReader("./config.example.toml")
	if !assert.NoError(err) {
		return
	}
	conf, err := New(confReader)
	if !assert.NoError(err) {
		return
	}

	assert.Equal(8890, conf.Viper.GetInt("server.grpc.port"))

	app1 := &app{}
	err = conf.Unmarshal("app", app1)
	if !assert.NoError(err) {
		return
	}
	assert.Equal("appname", app1.Name)
	assert.Equal([]string{"xxx1", "xxx2"}, app1.Maintainers)

	global := struct {
		App app `toml:"app"`
	}{}
	err = conf.UnmarshalALL(&global)
	if !assert.NoError(err) {
		return
	}
	assert.Equal("appname", global.App.Name)
	assert.Equal([]string{"xxx1", "xxx2"}, global.App.Maintainers)

	subConf, err := conf.Sub("domain.entity.entity1")
	if !assert.NoError(err) {
		return
	}
	assert.Equal("Tom", subConf.Viper.GetString("name"))
	assert.Equal(100, subConf.Viper.GetInt("age"))

	var finder skeleton.Finder[string, any] = conf
	v, err := finder.Find("domain.public_2")
	if !assert.NoError(err) {
		return
	}
	assert.Equal("domain公共配置2", v.(string))
	v2, err := finder.Find("server.http.port")
	if !assert.NoError(err) {
		return
	}
	assert.EqualValues(7890, v2.(int64))
}
