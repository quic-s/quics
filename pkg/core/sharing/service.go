package sharing

import (
	"errors"
	"io"

	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
)

type SharingService struct {
	historyRepository history.Repository
	syncRepository    sync.Repository
	sharingRepository Repository
	syncDir           SyncDirAdapter
}

var prefixLink = "https://" + config.GetRestServerAddress() + "/api/v1/download/files"

func NewService(historyRepository history.Repository, syncRepository sync.Repository, sharingRepository Repository, syncDir SyncDirAdapter) *SharingService {
	return &SharingService{
		historyRepository: historyRepository,
		syncRepository:    syncRepository,
		sharingRepository: sharingRepository,
		syncDir:           syncDir,
	}
}

func (ss *SharingService) CreateLink(request *types.ShareReq) (*types.ShareRes, error) {
	// get file for creating link
	file, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, err
	}

	// get file history for UUID to find last edited person
	fileHistory, err := ss.historyRepository.GetFileHistory(request.AfterPath, file.LatestSyncTimestamp)
	if err != nil {
		return nil, err
	}
	prefixLink = "https://" + config.GetRestServerAddress() + "/api/v1/download/files"
	// make link
	paramUUID := "?uuid=" + fileHistory.UUID
	paramFile := "&file=" + file.AfterPath
	link := prefixLink + paramUUID + paramFile

	// save link to database
	sharing := &types.Sharing{
		Link:     link,
		Count:    0,
		MaxCount: uint(request.MaxCnt),
		Owner:    request.UUID,
		File:     *file,
	}

	err = ss.sharingRepository.SaveLink(sharing)
	if err != nil {
		return nil, err
	}

	return &types.ShareRes{
		Link: link,
	}, nil
}

func (ss *SharingService) DeleteLink(request *types.StopShareReq) (*types.StopShareRes, error) {
	// get sharing data using link
	sharing, err := ss.sharingRepository.GetLink(request.Link)
	if err != nil {
		return nil, err
	}

	// check if the user is the owner of the link
	if sharing.Owner != request.UUID {
		return nil, errors.New("not authorized")
	}

	// delete link
	err = ss.sharingRepository.DeleteLink(request.Link)
	if err != nil {
		return nil, err
	}

	return &types.StopShareRes{
		UUID: request.UUID,
	}, nil
}

func (ss *SharingService) DownloadFile(uuid string, afterPath string) (*types.FileMetadata, io.Reader, error) {
	prefixLink = "https://" + config.GetRestServerAddress() + "/api/v1/download/files"
	paramUUID := "?uuid=" + uuid
	paramFile := "&file=" + afterPath
	link := prefixLink + paramUUID + paramFile

	// get sharing data using link
	sharing, err := ss.sharingRepository.GetLink(link)
	if err != nil {
		return nil, nil, err
	}

	// check if the link has been used up
	if sharing.Count >= sharing.MaxCount {
		err := ss.sharingRepository.DeleteLink(link)
		if err != nil {
			return nil, nil, err
		}

		return nil, nil, errors.New("link has been expired")
	}

	fileInfo, fileContent, err := ss.syncDir.GetFileFromHistoryDir(sharing.File.AfterPath, sharing.File.LatestSyncTimestamp)
	if err != nil {
		return nil, nil, err
	}

	// increase count
	sharing.Count++

	// update link
	err = ss.sharingRepository.UpdateLink(sharing)
	if err != nil {
		return nil, nil, err
	}

	return fileInfo, fileContent, nil
}
