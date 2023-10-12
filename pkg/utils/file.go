package utils

import (
	"strconv"
	"strings"
)

// GetNamesByAfterPath extracts root directory name and file name from afterPath
// FIXME: should think file with directories
func GetNamesByAfterPath(afterPath string) (string, string) {
	paths := strings.Split(afterPath, "/")
	rootDirName := paths[1]
	fileName := paths[2]
	return rootDirName, fileName
}

// GetHistoryFileNameByAfterPath returns history file name in history directory extracting from afterPath
func GetHistoryFileNameByAfterPath(afterPath string, timestamp uint64) string {
	rootDirName, fileName := GetNamesByAfterPath(afterPath)
	historyDirPath := GetQuicsHistoryPathByRootDir(rootDirName)
	historyFilePath := historyDirPath + fileName + "_" + strconv.FormatUint(timestamp, 10)
	return historyFilePath
}
