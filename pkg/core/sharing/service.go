package sharing

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
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
	fileHistory, err := ss.historyRepository.GetFileHistory(request.AfterPath, request.Version)
	if err != nil {
		return nil, err
	}

	// make link
	paramUUID := "?uuid=" + strings.ToLower(fileHistory.UUID)
	paramFile := "&file=" + strings.ToLower(utils.ExtractFileNameFromHistoryFile(fileHistory.AfterPath))
	paramVersion := "&version=" + fmt.Sprint(fileHistory.Timestamp)
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

func (ss *SharingService) DownloadFile(uuid string, afterPath string, timestamp string) (*os.File, fs.FileInfo, error) {
	paramUUID := "?uuid=" + strings.ToLower(uuid)
	paramFile := "&file=" + strings.ToLower(afterPath)
	paramVersion := "&version=" + strings.ToLower(timestamp)
	link := PrefixLink + paramUUID + paramFile + paramVersion

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
	historyFileName := sharing.File.BeforePath + utils.GetHistoryFileNameByAfterPath(sharing.File.AfterPath, sharing.File.Timestamp)
	downloadFileName := sharing.File.BeforePath + strings.ReplaceAll(afterPath, "-", "/")

	// open original history file
	historyFile, err := os.Open(historyFileName)
	if err != nil {
		return nil, nil, err
	}
	defer historyFile.Close()

	// create file for download (not history file)
	downloadFile, err := os.Create(downloadFileName)
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(downloadFileName)

	// copy history file to download file
	_, err = io.Copy(downloadFile, historyFile)
	if err != nil {
		return nil, nil, err
	}

	// get file information
	downloadFileInfo, err := downloadFile.Stat()
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

	return downloadFile, downloadFileInfo, nil
}
