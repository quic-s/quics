package sharing

import "github.com/quic-s/quics/pkg/history"

type Sharing struct {
	Id    uint
	Count uint
	Max   uint
	Link  string
	File  history.FileHistory // to share file at point that client wanted in time
}

// FileDownloadRequest is used when creating file download link
type FileDownloadRequest struct {
	Uuid       string `json:"uuid"`
	BeforePath string `json:"before_path"`
	AfterPath  string `json:"after_path"`
	Count      uint32 `json:"count"`
	MaxCount   uint32 `json:"max_count"`
}

// FileDownloadResponse is used when returning created file download link
type FileDownloadResponse struct {
	Link  string `json:"link"`
	Count uint32 `json:"count"`
}
