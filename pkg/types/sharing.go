package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Sharing struct {
	Id       uint
	Count    uint
	MaxCount uint
	Link     string
	Owner    string      // client uuid
	File     FileHistory // to share file at point that client wanted in time
}

// FileDownloadRequest is used when creating file download link
type FileDownloadRequest struct {
	Uuid       string
	BeforePath string
	AfterPath  string
	MaxCount   uint
}

func (fileDownloadRequest *FileDownloadRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(fileDownloadRequest)
}

// FileDownloadResponse is used when returning created file download link
type FileDownloadResponse struct {
	Link     string
	Count    uint
	MaxCount uint
}

func (fileDownloadResponse *FileDownloadResponse) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileDownloadResponse); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
