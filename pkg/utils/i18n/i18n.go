package i18n

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	bundle *Bundle
	once   sync.Once
)

// Bundle holds translations
type Bundle struct {
	locales map[string]*viper.Viper
	defaultLang string
}

// Init initializes the i18n bundle
func Init(configPath string, defaultLang string) error {
	var err error
	once.Do(func() {
		b := &Bundle{
			locales:     make(map[string]*viper.Viper),
			defaultLang: defaultLang,
		}

		// List yaml files in configPath/i18n or configPath
		// Assuming files are like en-US.yaml, zh-CN.yaml in configPath
		
		// For simplicity, let's look for specific patterns or list directory if possible
		// Since we can't easily glob with standard library inside a unknown struct without reading dir...
		// Let's assume the user manually loads or we load standard ones.
		// BETTER: Use filepath.Walk or ReadDir
		
		files, e := filepath.Glob(filepath.Join(configPath, "*.yaml"))
		if e != nil {
			err = e
			return
		}
		
		for _, f := range files {
			filename := filepath.Base(f)
			ext := filepath.Ext(filename)
			lang := strings.TrimSuffix(filename, ext) // en-US

			v := viper.New()
			v.SetConfigFile(f)
			v.SetConfigType("yaml")
			if e := v.ReadInConfig(); e != nil {
				// Warn?
				continue
			}
			b.locales[lang] = v
		}
		
		bundle = b
	})
	return err
}

// T translates a key
func T(lang, key string, args ...interface{}) string {
	if bundle == nil {
		return key
	}
	
	// Get locale or default
	v, ok := bundle.locales[lang]
	if !ok {
		// Try fallback
		v, ok = bundle.locales[bundle.defaultLang]
		if !ok {
			return key
		}
	}

	val := v.GetString(key)
	if val == "" {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(val, args...)
	}
	return val
}
