package fs

import (
	"crypto/sha1"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
)

type SyncDir struct {
	lockNum uint8
	pathMut map[byte]*sync.Mutex
	SyncDir string
}

func NewSyncDir(syncDir string) *SyncDir {
	lockNum := uint8(32)
	pathMut := make(map[byte]*sync.Mutex)

	for i := uint8(0); i < lockNum; i++ {
		pathMut[i] = &sync.Mutex{}
	}

	return &SyncDir{
		lockNum: uint8(lockNum),
		pathMut: pathMut,
		SyncDir: syncDir,
	}
}

// SyncFileToLatestDir creates/updates sync file to latest directory
func (s *SyncDir) SaveFileToLatestDir(afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error {
	// lock mutex by hash value of file path
	// using hash value is to reduce the number of mutex
	h := sha1.New()
	h.Write([]byte(afterPath))
	hash := h.Sum(nil)

	s.pathMut[uint8(hash[0]%s.lockNum)].Lock()
	defer s.pathMut[uint8(hash[0]%s.lockNum)].Unlock()

	latestFilePath := filepath.Join(s.SyncDir, afterPath)

	err := fileMetadata.WriteFileWithInfo(latestFilePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (s *SyncDir) GetFileFromLatestDir(afterPath string) (*types.FileMetadata, io.Reader, error) {
	latestFilePath := filepath.Join(s.SyncDir, afterPath)

	file, err := os.Open(latestFilePath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	return types.NewFileMetadataFromOSFileInfo(fileInfo), file, nil
}

func (s *SyncDir) DeleteFileFromLatestDir(afterPath string) error {
	// lock mutex by hash value of file path
	// using hash value is to reduce the number of mutex
	h := sha1.New()
	h.Write([]byte(afterPath))
	hash := h.Sum(nil)

	s.pathMut[uint8(hash[0]%s.lockNum)].Lock()
	defer s.pathMut[uint8(hash[0]%s.lockNum)].Unlock()

	latestFilePath := filepath.Join(s.SyncDir, afterPath)

	err := os.Remove(latestFilePath)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	reootToFile, _ := filepath.Split(afterPath)
	latestFileDir := filepath.Join(s.SyncDir, reootToFile)
	dir, err := os.Open(latestFileDir)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// Delete directory when it is empty
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		os.Remove(filepath.Join(latestFileDir))
	}

	return nil
}

func (s *SyncDir) SaveFileToConflictDir(uuid string, afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error {
	err := fileMetadata.WriteFileWithInfo(utils.GetConflictFileNameByAfterPath(afterPath, uuid), fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (s *SyncDir) GetFileFromConflictDir(afterPath string, uuid string) (*types.FileMetadata, io.Reader, error) {
	file, err := os.Open(utils.GetConflictFileNameByAfterPath(afterPath, uuid))
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	return types.NewFileMetadataFromOSFileInfo(fileInfo), file, nil
}

func (s *SyncDir) DeleteFilesFromConflictDir(afterpath string) error {
	rootToFileDir, fileName := filepath.Split(afterpath)
	// 정규식 객체를 생성합니다.
	re, err := regexp.Compile("^" + fileName + "_.*")
	if err != nil {
		return err
	}

	rootDir, fileDir := utils.GetNamesByAfterPath(rootToFileDir)
	// 삭제할 파일이 있는 디렉토리를 엽니다.
	dir, err := os.Open(filepath.Join(s.SyncDir, rootDir+".conflict", fileDir))
	if err != nil {
		return err
	}
	defer dir.Close()

	// 디렉토리의 모든 파일을 읽습니다.
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// 각 파일에 대해 정규식과 일치하는지 확인하고, 일치하면 삭제합니다.
	for _, file := range files {
		if re.MatchString(file.Name()) {
			err = os.Remove(filepath.Join(s.SyncDir, rootDir+".conflict", fileDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	// Delete directory when it is empty
	files, err = dir.Readdir(-1)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		os.Remove(filepath.Join(s.SyncDir, rootDir+".conflict", fileDir))
	}

	return nil
}

// SyncFileToHistoryDir creates/updates sync file to history directory
func (s *SyncDir) SaveFileToHistoryDir(afterPath string, timestamp uint64, fileMetadata *types.FileMetadata, fileContent io.Reader) error {
	// lock mutex by hash value of file path
	// using hash value is to reduce the number of mutex
	h := sha1.New()
	h.Write([]byte(afterPath))
	hash := h.Sum(nil)

	s.pathMut[uint8(hash[0]%s.lockNum)].Lock()
	defer s.pathMut[uint8(hash[0]%s.lockNum)].Unlock()

	// create history directory
	historyFilePath := utils.GetHistoryFileNameByAfterPath(afterPath, timestamp)

	err := fileMetadata.WriteFileWithInfo(historyFilePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (s *SyncDir) GetFileFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error) {
	file, err := os.Open(utils.GetHistoryFileNameByAfterPath(afterPath, timestamp))
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("quics: ", err)
		return nil, nil, err
	}

	return types.NewFileMetadataFromOSFileInfo(fileInfo), file, nil
}
