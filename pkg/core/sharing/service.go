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
		err = errors.New("[SharingService.CreateLink] get file by path: " + err.Error())
		return nil, err
	}

	// get file history for UUID to find last edited person
	fileHistory, err := ss.historyRepository.GetFileHistory(request.AfterPath, file.LatestSyncTimestamp)
	if err != nil {
		err = errors.New("[SharingService.CreateLink] get file history: " + err.Error())
		return nil, err
	}
	prefixLink := "https://" + config.GetRestServerAddress() + "/api/v1/download/files"
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
		err = errors.New("[SharingService.CreateLink] save link to repository: " + err.Error())
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
		err = errors.New("[SharingService.DeleteLink] get link from repository: " + err.Error())
		return nil, err
	}

	// check if the user is the owner of the link
	if sharing.Owner != request.UUID {
		return nil, errors.New("[SharingService.DeleteLink] user is not the owner of the link (UUID: " + request.UUID + ")")
	}

	// delete link
	err = ss.sharingRepository.DeleteLink(request.Link)
	if err != nil {
		err = errors.New("[SharingService.DeleteLink] delete link from repository: " + err.Error())
		return nil, err
	}

	return &types.StopShareRes{
		UUID: request.UUID,
	}, nil
}

func (ss *SharingService) DownloadFile(uuid string, afterPath string) (*types.FileMetadata, io.Reader, error) {
	prefixLink := "https://" + config.GetRestServerAddress() + "/api/v1/download/files"
	paramUUID := "?uuid=" + uuid
	paramFile := "&file=" + afterPath
	link := prefixLink + paramUUID + paramFile

	// get sharing data using link
	sharing, err := ss.sharingRepository.GetLink(link)
	if err != nil {
		err = errors.New("[SharingService.DownloadFile] get link from repository: " + err.Error())
		return nil, nil, err
	}

	// check if the link has been used up
	if sharing.Count >= sharing.MaxCount {
		err := ss.sharingRepository.DeleteLink(link)
		if err != nil {
			err = errors.New("[SharingService.DownloadFile] delete link from repository: " + err.Error())
			return nil, nil, err
		}

		return nil, nil, errors.New("[SharingService.DownloadFile] link has been used up")
	}

	fileInfo, fileContent, err := ss.syncDir.GetFileFromHistoryDir(sharing.File.AfterPath, sharing.File.LatestSyncTimestamp)
	if err != nil {
		err = errors.New("[SharingService.DownloadFile] get file from history dir: " + err.Error())
		return nil, nil, err
	}

	// increase count
	sharing.Count++

	// update link
	err = ss.sharingRepository.UpdateLink(sharing)
	if err != nil {
		err = errors.New("[SharingService.DownloadFile] update link from repository: " + err.Error())
		return nil, nil, err
	}

	return fileInfo, fileContent, nil
}
