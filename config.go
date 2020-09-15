package main

import (
	"github.com/jinzhu/configor"
	"github.com/mcuadros/go-defaults"
)

// App config struct
type Config struct {
	Endpoint string `default:"" json:"endpoint" form:"endpoint" query:"endpoint" required:"true"`
	Username string `default:"" json:"username" form:"username" query:"username" required:"true"`
	Password string `default:"" json:"password" form:"password" query:"password" required:"true"`
}

// Create config from configFile
func NewConfig(configFile string) (*Config, error) {

	config := new(Config)
	defaults.SetDefaults(config)

	if err := configor.Load(config, configFile); err != nil {
		return nil, err
	}
	return config, nil
}
