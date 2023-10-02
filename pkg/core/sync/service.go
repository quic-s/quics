package sync

import "github.com/quic-s/quics/pkg/types"

type MySyncService struct {
	syncRepository Repository
}

func NewSyncService(syncRepository Repository) *MySyncService {
	return &MySyncService{
		syncRepository: syncRepository,
	}
}

// IsExistFile checks if file exists in database
func (service *MySyncService) IsExistFile(path string) int {
	file := service.syncRepository.GetFileByPath(path)
	if file == nil {
		return 0
	} else {
		return 1
	}
}

// SaveFileFromPleaseSync saves file from please sync request
func (service *MySyncService) SaveFileFromPleaseSync(path string, file types.File) error {
	err := service.syncRepository.SaveFileByPath(path, file)
	if err != nil {
		return err
	}

	return nil
}

// GetFilesByRootDir returns files by root directory path
func (service *MySyncService) GetFilesByRootDir(rootDirPath string) []types.File {
	filesByRootDir := make([]types.File, 0)
	files := service.syncRepository.GetAllFiles()

	for _, file := range files {
		if file.RootDir.AfterPath == rootDirPath {
			filesByRootDir = append(files, file)
		}
	}

	return filesByRootDir
}

// GetFiles returns all files in database
func (service *MySyncService) GetFiles() []types.File {
	return service.syncRepository.GetAllFiles()
}

// GetFileByPath returns file entity by path
func (service *MySyncService) GetFileByPath(path string) *types.File {
	return service.syncRepository.GetFileByPath(path)
}
