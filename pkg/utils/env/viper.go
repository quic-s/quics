package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetViperEnvVariables(key string) string {
	envPath := filepath.Join(GetDirPath(), ".qis.env")
	_, err := os.Stat(envPath)
	if err != nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file: ", err)
		return ""
	}

	value := viper.Get(key).(string)

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
