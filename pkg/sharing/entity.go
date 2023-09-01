package sharing

import "github.com/quic-s/quics/pkg/history"

type Sharing struct {
	Id       uint
	Count    uint
	MaxCount uint
	Link     string
	Owner    string              // client uuid
	File     history.FileHistory // to share file at point that client wanted in time
}

// FileDownloadRequest is used when creating file download link
type FileDownloadRequest struct {
	Uuid       string
	BeforePath string
	AfterPath  string
	MaxCount   uint
}

// FileDownloadResponse is used when returning created file download link
type FileDownloadResponse struct {
	Link     string
	Count    uint
	MaxCount uint
}
