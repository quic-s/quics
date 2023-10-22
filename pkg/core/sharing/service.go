package sharing

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
)

type SharingService struct {
	historyRepository history.Repository
	syncRepository    sync.Repository
	sharingRepository Repository
}

const (
	// FIXME: this link is for testing purposes only
	PrefixLink = "http://localhost:6121/api/v1/download/files"
)

func NewService(historyRepository history.Repository, syncRepository sync.Repository, sharingRepository Repository) *SharingService {
	return &SharingService{
		historyRepository: historyRepository,
		syncRepository:    syncRepository,
		sharingRepository: sharingRepository,
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

	// make link
	paramUUID := "?uuid=" + strings.ToLower(fileHistory.UUID)
	paramFile := "&file=" + strings.ToLower(file.AfterPath)
	link := PrefixLink + paramUUID + paramFile

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

func (ss *SharingService) DownloadFile(uuid string, afterPath string) (*os.File, fs.FileInfo, error) {
	paramUUID := "?uuid=" + strings.ToLower(uuid)
	paramFile := "&file=" + strings.ToLower(afterPath)
	link := PrefixLink + paramUUID + paramFile

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

	// get file
	fileName := sharing.File.BeforePath + sharing.File.AfterPath

	// open original history file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// get file information
	fileInfo, err := file.Stat()
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

	return file, fileInfo, nil
}
