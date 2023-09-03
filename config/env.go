package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GetViperEnvVariables get env variables in env file using viper
func GetViperEnvVariables(key string) string {
	envPath := filepath.Join(GetDirPath(), ".qis.env")

	_, err := os.Stat(envPath)
	if err != nil {
		sourceEnvPath := ".env"
		sourceViper := viper.New()
		sourceViper.SetConfigFile(sourceEnvPath)
		sourceViper.SetConfigType("env")

		if err := sourceViper.ReadInConfig(); err != nil {
			log.Fatalf("Error while reading source config file: %s", err)
			return ""
		}

		for _, key := range sourceViper.AllKeys() {
			value := sourceViper.Get(key)
			viper.Set(key, value)
		}

		if err := viper.WriteConfigAs(envPath); err != nil {
			log.Fatalf("Error while writing config file: %s", err)
		}

	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file: ", err)
		return ""
	}

	value := viper.GetString(key)
	return value
}

func WriteViperEnvVariables(key string, value string) {
	envPath := filepath.Join(GetDirPath(), ".qis.env")
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")

		err = viper.WriteConfigAs(envPath)
		if err != nil {
			log.Fatalf("Error while writing config file: %s", err)
		}
	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file: ", err)
	}

	viper.Set(key, value)
	err = viper.WriteConfigAs(envPath)
	if err != nil {
		log.Fatalf("Error while writing config file: %s", err)
	}
}