// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package setting

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	GlobalConfig *Config
)

type Loader struct {
	vp *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader(configPath, configName, configType string) (*Loader, error) {
	vp := viper.New()
	vp.AddConfigPath(configPath)
	vp.SetConfigName(configName)
	vp.SetConfigType(configType)

	// Environment variables support
	// APP_SERVER_PORT will map to Server.Port
	vp.SetEnvPrefix("APP")
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.AutomaticEnv()

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired?
			// For now, let's treat it as error unless we want strictly env vars
			return nil, err
		} else {
			return nil, err
		}
	}

	return &Loader{vp: vp}, nil
}

// Load loads the configuration into the global struct
func (l *Loader) Load() (*Config, error) {
	var cfg Config
	if err := l.vp.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	GlobalConfig = &cfg
	return &cfg, nil
}

// Watch starts watching the config file for changes (hot reload)
func (l *Loader) Watch(callback func(*Config)) {
	l.vp.OnConfigChange(func(e fsnotify.Event) {
		var cfg Config
		if err := l.vp.Unmarshal(&cfg); err == nil {
			GlobalConfig = &cfg
			if callback != nil {
				callback(&cfg)
			}
		}
	})
	l.vp.WatchConfig()
}
