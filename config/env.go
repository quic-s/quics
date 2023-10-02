package config

import (
	"github.com/quic-s/quics/pkg/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// GetViperEnvVariables gets env variables in env file using viper
func GetViperEnvVariables(key string) string {
	envPath := filepath.Join(utils.GetQuicsDirPath(), "qis.env")

	_, err := os.Stat(envPath)
	if err != nil {
		sourceEnvPath := ".env"
		sourceViper := viper.New()
		sourceViper.SetConfigFile(sourceEnvPath)
		sourceViper.SetConfigType("env")

		if err := sourceViper.ReadInConfig(); err != nil {
			log.Fatalf("quics: Error while reading source config file: %s", err)
			return ""
		}

		for _, key := range sourceViper.AllKeys() {
			value := sourceViper.Get(key)
			viper.Set(key, value)
		}

		if err := viper.WriteConfigAs(envPath); err != nil {
			log.Fatalf("quics: Error while writing config file: %s", err)
		}

	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("quics: Error while reading config file: ", err)
		return ""
	}

	value := viper.GetString(key)
	return value
}

// WriteViperEnvVariables writes env variables to env file using viper
func WriteViperEnvVariables(key string, value string) {
	envPath := filepath.Join(utils.GetQuicsDirPath(), ".qis.env")
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")

		err = viper.WriteConfigAs(envPath)
		if err != nil {
			log.Fatalf("quics: Error while writing config file: %s", err)
		}
	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("quics: Error while reading config file: ", err)
	}

	viper.Set(key, value)
	err = viper.WriteConfigAs(envPath)
	if err != nil {
		log.Fatalf("quics: Error while writing config file: %s", err)
	}
}
