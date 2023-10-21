package sharing

import (
	"errors"
	"fmt"
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
	file, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, err
	}

	fileHistory, err := ss.historyRepository.GetFileHistory(file.AfterPath, file.LatestSyncTimestamp)
	if err != nil {
		return nil, err
	}

	// make link
	paramUUID := "?uuid=" + strings.ToLower(request.UUID)
	paramFile := "&file=" + strings.ToLower(strings.ReplaceAll(file.AfterPath, "/", "-"))
	paramVersion := "&version=" + fmt.Sprint(file.LatestSyncTimestamp)
	link := PrefixLink + paramUUID + paramFile + paramVersion

	// save link to database
	sharing := &types.Sharing{
		Link:     link,
		Count:    0,
		MaxCount: uint(request.MaxCnt),
		Owner:    request.UUID,
		File:     *fileHistory,
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
