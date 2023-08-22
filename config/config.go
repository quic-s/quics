package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

var RuntimeConf = &RuntimeConfig{} // RuntimeConfig as pointer

type RuntimeConfig struct {
	Database *Database `yaml:"database"` // Database as pointer
	Server   *Server   `yaml:"server"`   // Server as pointer
}

type Database struct {
	Path string `yaml:"path"`
}

type Server struct {
	Port int `yaml:"port"`
}

func init() {
	profile := initProfile()
	setRuntimeConfig(profile)
}

func initProfile() string {
	var profile string
	profile = os.Getenv("GO_PROFILE")
	if len(profile) <= 0 {
		profile = "local" // FIXME: if you want to change profile, then edit this.
	}
	fmt.Println("GOLANG_PROFILE: " + profile)
	return profile
}

// setRuntimeConfig
// based on profile, read config file and set configuration at runtime
func setRuntimeConfig(profile string) {
	// initialize viper configuration
	viper.AddConfigPath(".")
	viper.SetConfigName(profile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Unmarshal the configuration information to global variables
	// to use everywhere you want
	err = viper.Unmarshal(RuntimeConf)
	if err != nil {
		panic(err)
	}

	// when the configuration file is changed, then re-unmarshal according to changed events
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)

		var err error
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = viper.Unmarshal(RuntimeConf)
		if err != nil {
			fmt.Println(err)
			return
		}

	})
	viper.WatchConfig()
}
