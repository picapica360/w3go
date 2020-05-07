package config

import (
	"fmt"
	"path"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/picapica360/w3go/config/env"
	"github.com/picapica360/w3go/utils/file"
)

const (
	prefix    = "app"
	extension = ".toml"
)

var (
	conf *AppConfig = &AppConfig{} // global

	once sync.Once
)

// Init the config
func Init() {
	once.Do(func() {
		decodeToml(configFilename(env.Env()), conf)
	})
}

// decodeToml decodes the content in toml file to struct.
// filename is the file name in root directory.
// v is pointer of struct.
func decodeToml(filename string, v interface{}) {
	fpath := path.Join(env.Root(), filename)
	if ok, _ := file.Exists(fpath); !ok {
		panic(fmt.Sprintf(`[config] the config file "%s" not found in root directory.`, fpath))
	}

	if _, err := toml.DecodeFile(fpath, v); err != nil {
		panic(fmt.Sprintf(`[config] the app config initialize failure, filepath:"%s"; err:%v`, fpath, err))
	}
}

// Conf get the config from the 'app.[env].conf' file in root.
// note: the config would be built when app init, and store singleton.
func Conf() *AppConfig {
	return conf
}

func configFilename(env string) string {
	if env == "" {
		return prefix + extension
	}

	return prefix + "." + env + extension
}
