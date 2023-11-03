package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics/pkg/utils"
	"github.com/spf13/viper"
)

const (
	DefaultRestServerAddr   = "localhost"
	DefaultRestServerPort   = "6120"
	DefaultRestServerH3Port = "6121"

	DefaultQuicsPort = "6122"

	DefaultPassword = "quics"

	DefaultQuicsCertName = "cert-quics.pem"
	DefaultQuicsKeyName  = "key-quics.pem"
)

func init() {
	_, err := os.Stat(utils.GetQuicsDirPath())
	if os.IsNotExist(err) {
		err := os.Mkdir(utils.GetQuicsDirPath(), 0755)
		if err != nil {
			log.Fatalln("quics err: while creating .quics directory: ", err)
		}
	}

	envPath := filepath.Join(utils.GetQuicsDirPath(), "qis.env")

	_, err = os.Stat(envPath)
	if os.IsNotExist(err) {
		// Create default .env file
		sourceViper := viper.New()
		sourceViper.SetConfigType("env")

		_, err = os.Create(envPath)
		if err != nil {
			log.Fatalln("quics err: while create qis.env", err)
			return
		}

		// Set default values for env variables
		if restServerAddr := os.Getenv("REST_SERVER_ADDR"); restServerAddr != "" {
			sourceViper.Set("REST_SERVER_ADDR", restServerAddr)
		} else {
			sourceViper.Set("REST_SERVER_ADDR", DefaultRestServerAddr)
		}
		if restServerPort := os.Getenv("REST_SERVER_PORT"); restServerPort != "" {
			sourceViper.Set("REST_SERVER_PORT", restServerPort)
		} else {
			sourceViper.Set("REST_SERVER_PORT", DefaultRestServerPort)
		}
		if restServerH3Port := os.Getenv("REST_SERVER_H3_PORT"); restServerH3Port != "" {
			sourceViper.Set("REST_SERVER_H3_PORT", restServerH3Port)
		} else {
			sourceViper.Set("REST_SERVER_H3_PORT", DefaultRestServerH3Port)
		}
		if quicsPort := os.Getenv("QUICS_PORT"); quicsPort != "" {
			sourceViper.Set("QUICS_PORT", quicsPort)
		} else {
			sourceViper.Set("QUICS_PORT", DefaultQuicsPort)
		}
		if password := os.Getenv("PASSWORD"); password != "" {
			sourceViper.Set("PASSWORD", password)
		} else {
			sourceViper.Set("PASSWORD", DefaultPassword)
		}
		if quicsCertName := os.Getenv("QUICS_CERT_NAME"); quicsCertName != "" {
			sourceViper.Set("QUICS_CERT_NAME", quicsCertName)
		} else {
			sourceViper.Set("QUICS_CERT_NAME", DefaultQuicsCertName)
		}
		if quicsKeyName := os.Getenv("QUICS_KEY_NAME"); quicsKeyName != "" {
			sourceViper.Set("QUICS_KEY_NAME", quicsKeyName)
		} else {
			sourceViper.Set("QUICS_KEY_NAME", DefaultQuicsKeyName)
		}

		if err := sourceViper.WriteConfigAs(envPath); err != nil {
			log.Fatalln("quics err: ", err)
			return
		}
	} else if err != nil {
		log.Panicf("quics err: while reading config file: %s", err)
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("env")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("quics err: while reading config file: ", err)
		return
	}
}

// GetViperEnvVariables gets env variables in env file using viper
func GetViperEnvVariables(key string) string {
	value := viper.GetString(key)
	return value
}

// WriteViperEnvVariables writes env variables to env file using viper
func WriteViperEnvVariables(key string, value string) error {
	envPath := filepath.Join(utils.GetQuicsDirPath(), "qis.env")

	viper.Set(key, value)
	err := viper.WriteConfigAs(envPath)
	if err != nil {
		err = errors.New("while writing config file: " + err.Error())
		return err
	}
	return nil
}
