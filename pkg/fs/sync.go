package fs

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-protocol/pkg/types/fileinfo"
	"github.com/quic-s/quics/pkg/utils"
)

type SyncDir struct {
	SyncDir string
}

func NewSyncDir(syncDir string) *SyncDir {
	return &SyncDir{
		SyncDir: syncDir,
	}
}

// SyncFileToLatestDir creates/updates sync file to latest directory
func (s *SyncDir) CopyHistoryFileToLatestDir(afterPath string, timestamp uint64, fileInfo *fileinfo.FileInfo) error {
	historyFilePath := utils.GetHistoryFileNameByAfterPath(afterPath, timestamp)

	latestFilePath := filepath.Join(s.SyncDir, afterPath)

	if fileInfo.IsDir {
		err := os.MkdirAll(latestFilePath, fileInfo.Mode)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		file, err := os.Open(latestFilePath)
		if err != nil {
			return err
		}

		// Set file metadata.
		err = file.Chmod(fileInfo.Mode)
		if err != nil {
			return err
		}
		err = os.Chtimes(latestFilePath, fileInfo.ModTime, fileInfo.ModTime)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		return nil
	}

	// copy history file to latest file
	historyFile, err := os.Open(historyFilePath)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	defer historyFile.Close()

	// Open file with O_TRUNC flag to overwrite the file when the file already exists.
	latestFile, err := os.OpenFile(latestFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode)
	if err != nil {
		// If the file does not exist, create the file.
		if os.IsNotExist(err) {
			dir, _ := filepath.Split(latestFilePath)
			if dir != "" {
				err := os.MkdirAll(dir, 0700)
				if err != nil {
					return err
				}
			}
			latestFile, err = os.Create(latestFilePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer latestFile.Close()

	n, err := io.Copy(latestFile, historyFile)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	if n != fileInfo.Size {
		return errors.New("quics: copied file size is not equal to original file size")
	}

	// Set file metadata.
	err = latestFile.Chmod(fileInfo.Mode)
	if err != nil {
		return err
	}
	err = os.Chtimes(latestFilePath, fileInfo.ModTime, fileInfo.ModTime)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (s *SyncDir) DeleteFileFromLatestDir(afterPath string) error {
	latestFilePath := filepath.Join(s.SyncDir, afterPath)

	err := os.Remove(latestFilePath)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (s *SyncDir) SaveFileToConflictDir(afterPath string, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
	rootDirName, fileName := utils.GetNamesByAfterPath(afterPath)
	filePath := utils.GetQuicsConflictPathByRootDir(rootDirName)

	err := fileInfo.WriteFileWithInfo(filepath.Join(filePath, fileName), fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

// SyncFileToHistoryDir creates/updates sync file to history directory
func (s *SyncDir) SaveFileToHistoryDir(afterPath string, timestamp uint64, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
	// create history directory
	historyFilePath := utils.GetHistoryFileNameByAfterPath(afterPath, timestamp)

	err := fileInfo.WriteFileWithInfo(historyFilePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}
