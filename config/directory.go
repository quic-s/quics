package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ViperConfig struct {
	Key   string
	Value string
}

// GetDirPath returns the path of the .quics directory
func GetDirPath() string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics")
}

// ReadEnvFile read .qis.env file if it is existed
func ReadEnvFile() []map[string]string {
	envPath := filepath.Join(GetDirPath(), ".qis.env")
	file, err := os.Open(envPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file
	data, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	// make list with key-value format
	dataStr := string(data)

	// separate string with new line (\n)
	lines := strings.Split(dataStr, "\n")

	var kvList []map[string]string
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")

		key := parts[0]
		value := strings.Join(parts[1:], " ")
		kvMap := map[string]string{key: value}

		kvList = append(kvList, kvMap)
	}

	log.Println(kvList)
	return kvList
}

func GetRootDir() []map[string]string {
	rawList := ReadEnvFile()
	var kvList []map[string]string
	for _, kvMap := range rawList {
		for key, value := range kvMap {
			if len(key) > 5 && key[:5] == "ROOT." {
				kvList = append(kvList, map[string]string{key[5:]: value})
			}
		}
	}
	return kvList
}
