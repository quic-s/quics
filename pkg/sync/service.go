package sync

import "github.com/quic-s/quics/pkg/types"

type Service struct {
	syncRepository *Repository
}

func NewSyncService(syncRepository *Repository) *Service {
	return &Service{syncRepository: syncRepository}
}

func (syncService *Service) IsExistFile(path string) int {
	file := syncService.syncRepository.GetFileByPath(path)
	if file == nil {
		return 0
	} else {
		return 1
	}
}

func (syncService *Service) CheckIsOccurredConflict(path string, request types.PleaseSync) (*types.File, int) {
	file := syncService.syncRepository.GetFileByPath(path)

	if file.LatestSyncTimestamp >= request.LastUpdateTimestamp {
		// conflict
		return file, 1
	} else {
		return file, 0
	}
}

func (syncService *Service) SaveFileFromPleaseSync(path string, file types.File) error {
	err := syncService.syncRepository.SaveFileByPath(path, file)
	if err != nil {
		return err
	}

	return nil
}

func (syncService *Service) GetFilesByRootDir(rootDirPath string) []types.File {
	filesByRootDir := make([]types.File, 0)
	files := syncService.syncRepository.GetAllFiles()

	for _, file := range files {
		if file.RootDir.Path == rootDirPath {
			filesByRootDir = append(files, file)
		}
	}

	return filesByRootDir
}

func (syncService *Service) GetFiles() []types.File {
	return syncService.syncRepository.GetAllFiles()
}

func (syncService *Service) GetFileByPath(path string) *types.File {
	return syncService.syncRepository.GetFileByPath(path)
}
