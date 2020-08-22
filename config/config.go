package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/jinzhu/configor"
	"github.com/mcuadros/go-defaults"
)

// App config struct
type Config struct {
	Endpoint string `default:"" json:"endpoint" form:"endpoint" query:"endpoint"`
	Username string `default:"" json:"username" form:"username" query:"username"`
	Password string `default:"" json:"password" form:"password" query:"password"`
}

// Create config from configFile
func NewConfig(configFile string) (*Config, error) {

	config := new(Config)
	defaults.SetDefaults(config)

	configBytes, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			return nil, err
		}
	}

	if err := configor.Load(config); err != nil {
		return nil, err
	}
	return config, nil
}
