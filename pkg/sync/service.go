package sync

import "github.com/quic-s/quics/pkg/types"

type Service struct {
	syncRepository *Repository
}

func NewSyncService(syncRepository *Repository) *Service {
	return &Service{syncRepository: syncRepository}
}

func (syncService *Service) CheckIsOccurredConflict(path string, request types.PleaseSync) (*types.File, int) {
	file := syncService.syncRepository.GetFileByPath(path)

	if file.LatestSyncTimestamp < request.LastUpdatedTimestamp {
		// conflict
		return file, 1
	} else {
		return file, 0
	}
}

//func (syncService *Service) SaveFileFromPleaseSync(path string)
