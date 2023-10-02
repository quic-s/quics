package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

type MessageData interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

// ClientRegisterReq is used when registering client from client to server
type ClientRegisterReq struct {
	UUID           string // client
	ClientPassword string // client
}

// ClientDisconnectorReq is used when disconnecting client with server from client to server
type ClientDisconnectorReq struct {
	UUID           string // client
	ServerPassword string // server
}

// RegisterRootDirReq is used when registering root directory of a client from client to server
type RegisterRootDirReq struct {
	UUID            string
	RootDirPassword string
	BeforePath      string
	AfterPath       string
}

// SyncRootDirReq is used when synchronizing root directory of a client from client to server
type SyncRootDirReq struct {
	UUID            string
	RootDirPassword string
	BeforePath      string
	AfterPath       string
}

// PleaseFileMetaReq is used when client request file's metadata to server
type PleaseFileMetaReq struct {
	UUID      string
	AfterPath string
}

// PleaseFileMetaRes is used when server response the latest file's metadata to client
type PleaseFileMetaRes struct {
	UUID                string // who fixed this file last time
	AfterPath           string
	LatestHash          string
	LatestSyncTimestamp uint64
	ModifiedDate        string
}

// PleaseSyncReq is used when updating file's changes from client to server
type PleaseSyncReq struct {
	UUID                string
	Event               string
	BeforePath          string
	AfterPath           string
	LastUpdateTimestamp uint64
	LastUpdateHash      string
}

// PleaseSyncRes is used to response to client of whether file is updated or not
type PleaseSyncRes struct {
	UUID      string
	AfterPath string
}

// PleaseTakeReq is used when client synchronize file to server
type PleaseTakeReq struct {
	UUID      string
	AfterPath string
}

// PleaseTakeRes is used to response to client of whether file is synchronized or not
type PleaseTakeRes struct {
	UUID      string
	AfterPath string
}

// MustSyncReq is used to inform whether file is updated or not from server to client
type MustSyncReq struct {
	LatestHash          string
	LatestSyncTimestamp uint64
	BeforePath          string
	AfterPath           string
}

// MustSyncRes is used to response to server that client will synchronize file
type MustSyncRes struct {
	UUID                string
	AfterPath           string
	LatestSyncTimestamp uint64
	LatestSyncHash      string
}

// GiveYouReq is used when sending file to client
type GiveYouReq struct {
	UUID      string
	AfterPath string
}

// GiveYouRes is used to response to server that client received file
type GiveYouRes struct {
	UUID              string
	AfterPath         string
	LastSyncTimestamp uint64
	LastHash          string
}

// PleaseFileReq is used when client request file to server (metadata)
type PleaseFileReq struct {
	UUID          string
	SyncTimestamp uint64
	BeforePath    string
	AfterPath     string
}

// PleaseFileRes is used when server response file to client (metadata)
type PleaseFileRes struct {
	UUID      string
	AfterPath string
}

// FileDownloadReq is used when creating file download link
type FileDownloadReq struct {
	UUID       string
	BeforePath string
	AfterPath  string
	MaxCount   uint
}

// FileDownloadRes is used when returning created file download link
type FileDownloadRes struct {
	Link     string
	Count    uint
	MaxCount uint
}

func (clientRegisterReq *ClientRegisterReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(clientRegisterReq); err != nil {
		log.Println("quics: (ClientRegisterReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (clientRegisterReq *ClientRegisterReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(clientRegisterReq)
}

func (clientDisconnectorReq *ClientDisconnectorReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(clientDisconnectorReq); err != nil {
		log.Println("quics: (ClientDisconnectorReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (clientDisconnectorReq *ClientDisconnectorReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(clientDisconnectorReq)
}

func (registerRootDirReq *RegisterRootDirReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerRootDirReq); err != nil {
		log.Println("quics: (RegisterRootDirReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (registerRootDirReq *RegisterRootDirReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerRootDirReq)
}

func (syncRootDirReq *SyncRootDirReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(syncRootDirReq); err != nil {
		log.Println("quics: (SyncRootDirReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (syncRootDirReq *SyncRootDirReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(syncRootDirReq)
}

func (pleaseFileMetaReq *PleaseFileMetaReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseFileMetaReq); err != nil {
		log.Println("quics: (PleaseFileMetaReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseFileMetaReq *PleaseFileMetaReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFileMetaReq)
}

func (pleaseFileMetaRes *PleaseFileMetaRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseFileMetaRes); err != nil {
		log.Println("quics: (PleaseFileMetaRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseFileMetaRes *PleaseFileMetaRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFileMetaRes)
}

func (pleaseSyncReq *PleaseSyncReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseSyncReq); err != nil {
		log.Println("quics: (PleaseSyncReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseSyncReq *PleaseSyncReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSyncReq)
}

func (pleaseSyncRes *PleaseSyncRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseSyncRes); err != nil {
		log.Println("quics: (PleaseSyncRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseSyncRes *PleaseSyncRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSyncRes)
}

func (pleaseTakeReq *PleaseTakeReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseTakeReq); err != nil {
		log.Println("quics: (PleaseTakeReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseTakeReq *PleaseTakeReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseTakeReq)
}

func (pleaseTakeRes *PleaseTakeRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseTakeRes); err != nil {
		log.Println("quics: (PleaseTakeRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseTakeRes *PleaseTakeRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseTakeRes)
}

func (mustSyncReq *MustSyncReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSyncReq); err != nil {
		log.Println("quics: (MustSyncReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (mustSyncReq *MustSyncReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(mustSyncReq)
}

func (mustSyncRes *MustSyncRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSyncRes); err != nil {
		log.Println("quics: (MustSyncRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (mustSyncRes *MustSyncRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(mustSyncRes)
}

func (giveYouReq *GiveYouReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(giveYouReq); err != nil {
		log.Println("quics: (GiveYouReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (giveYouReq *GiveYouReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(giveYouReq)
}

func (giveYouRes *GiveYouRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(giveYouRes); err != nil {
		log.Println("quics: (GiveYouRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (giveYouRes *GiveYouRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(giveYouRes)
}

func (pleaseFileReq *PleaseFileReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseFileReq); err != nil {
		log.Println("quics: (PleaseFileReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseFileReq *PleaseFileReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFileReq)
}

func (pleaseFileRes *PleaseFileRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseFileRes); err != nil {
		log.Println("quics: (PleaseFileRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (pleaseFileRes *PleaseFileRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFileRes)
}

func (fileDownloadReq *FileDownloadReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileDownloadReq); err != nil {
		log.Println("quics: (FileDownloadReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (fileDownloadReq *FileDownloadReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(fileDownloadReq)
}

func (fileDownloadRes *FileDownloadRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileDownloadRes); err != nil {
		log.Println("quics: (FileDownLoadRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (fileDownloadRes *FileDownloadRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(fileDownloadRes)
}
