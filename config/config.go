package config

import (
	"github.com/spf13/viper"
	"log"
)

func main() {
	// Initialize viper configuration
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
