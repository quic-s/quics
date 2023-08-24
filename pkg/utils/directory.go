package utils

import (
	"log"
	"os"
	"path/filepath"
)

// CreateDirIfNotExist
// Create the quics folder if it does not exist
func CreateDirIfNotExist() {
	quicsDir := GetDirPath()
	_, err := os.Stat(quicsDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(quicsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created quics directory: ", quicsDir)
	} else {
		log.Println("Using existing quics directory: ", quicsDir)
	}
}

// GetDirPath
// Return the path of the quics directory
func GetDirPath() string {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(tempDir, "quics")
}
