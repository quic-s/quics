package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics/pkg/utils"
	"github.com/spf13/viper"
)

// GetViperEnvVariables gets env variables in env file using viper
func GetViperEnvVariables(key string) string {
	_, err := os.Stat(utils.GetQuicsDirPath())
	if os.IsNotExist(err) {
		err := os.Mkdir(utils.GetQuicsDirPath(), 0755)
		if err != nil {
			log.Fatalf("quics: Error while creating .quics directory: %s", err)
		}
	}

	envPath := filepath.Join(utils.GetQuicsDirPath(), "qis.env")

	_, err = os.Stat(envPath)
	if os.IsNotExist(err) {
		sourceEnvPath := "./.env"
		sourceViper := viper.New()
		sourceViper.SetConfigFile(sourceEnvPath)
		sourceViper.SetConfigType("env")

		if err := sourceViper.ReadInConfig(); err != nil {
			log.Println("quics: ", err)
			return ""
		}

		_, err := os.Create(envPath)
		if err != nil {
			log.Println("quics: ", err)
			return ""
		}
		if err := sourceViper.WriteConfigAs(envPath); err != nil {
			log.Println("quics: ", err)
			return ""
		}
	} else if err != nil {
		log.Panicf("quics: Error while reading config file: %s", err)
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("env")
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
