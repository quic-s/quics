package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetQuicsDirPath $HOME/.quics
func GetQuicsDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(homeDir, ".quics")
}

// GetQuicsSyncDirPath $HOME/.quics/sync
func GetQuicsSyncDirPath() string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync")
}

// GetQuicsRootDirPath $HOME/.quics/sync/{rootDir}
func GetQuicsRootDirPath(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir)
}

// GetQuicsHistoryPathByRootDir $HOME/.quics/sync/{rootDir}.history
func GetQuicsHistoryPathByRootDir(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir+".history")
}

// GetQuicsConflictPathByRootDir $HOME/.quics/sync/{rootDir}.conflict
func GetQuicsConflictPathByRootDir(rootDir string) string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, ".quics", "sync", rootDir+".conflict")
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
