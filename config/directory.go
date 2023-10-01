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

// GetQuicsDirPath returns the path of the .quics directory
func GetQuicsDirPath() string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics") // $HOME/.quics
}

// GetQuicsSyncDirPath returns the path of the .quics/sync directory
func GetQuicsSyncDirPath() string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync") // $HOME/.quics/sync
}

// GetQuicsRootDirPath returns the path of the ./quics/sync/{rootDir} directory
func GetQuicsRootDirPath(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir[1:]) // $HOME/.quics/sync/{rootDir}
}

// getQuicsHistoryPathByRootDir returns the path of the ./quics/sync/{rootDir}/history directory
func getQuicsHistoryPathByRootDir(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir[1:], "history") // $HOME/.quics/sync/{rootDir}/history
}

// getQuicsLatestPathByRootDir returns the path of the ./quics/sync/{rootDir}/latest directory
func getQuicsLatestPathByRootDir(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir[1:], "latest") // $HOME/.quics/sync/{rootDir}/latest
}

// ReadEnvFile reads .qis.env file if it is existed
func ReadEnvFile() []map[string]string {
	envPath := filepath.Join(GetQuicsDirPath(), "qis.env")
	file, err := os.Open(envPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read .qis.env file
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
