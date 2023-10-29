package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

const (
	REGISTERCLIENT    = "REGISTERCLIENT"
	DISCONNECTCLIENT  = "DISCONNECTCLIENT"
	REGISTERROOTDIR   = "REGISTERROOTDIR"
	SYNCROOTDIR       = "SYNCROOTDIR"
	GETROOTDIRS       = "GETROOTDIRS"
	DISCONNECTROOTDIR = "DISCONNECTROOTDIR"
	PLEASESYNC        = "PLEASESYNC"
	MUSTSYNC          = "MUSTSYNC"
	FORCESYNC         = "FORCESYNC"
	CONFLICT          = "CONFLICT"
	CONFLICTLIST      = "CONFLICTLIST"
	CONFLICTDOWNLOAD  = "CONFLICTDOWNLOAD"
	CHOOSEONE         = "CHOOSEONE"
	FULLSCAN          = "FULLSCAN"
	RESCAN            = "RESCAN"
	NEEDCONTENT       = "NEEDCONTENT"
	PING              = "PING"
	ROLLBACK          = "ROLLBACK"
	HISTORYSHOW       = "HISTORYSHOW"
	HISTORYDOWNLOAD   = "HISTORYDOWNLOAD"
	DOWNLOAD          = "DOWNLOAD"
	STARTSHARING      = "STARTSHARING"
	STOPSHARING       = "STOPSHARING"
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

type ClientRegisterRes struct {
	UUID string // client
}

// ClientDisconnectorReq is used when disconnecting client with server from client to server
type DisconnectClientReq struct {
	UUID           string // client
	ServerPassword string // server
}

type DisconnectClientRes struct {
	UUID string // client
}

type AskRootDirReq struct {
	UUID string
}

type AskRootDirRes struct {
	RootDirList []string
}

type AskConflictListReq struct {
	UUID string
}

type AskConflictListRes struct {
	UUID      string
	Conflicts []Conflict
}

// RootDirReqRegister is used when registering root directory of a client from client to server
type RootDirRegisterReq struct {
	UUID            string
	RootDirPassword string
	BeforePath      string
	AfterPath       string
}

type RootDirRegisterRes struct {
	UUID string
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
	AfterPath           string
	LastUpdateTimestamp uint64
	LastUpdateHash      string
	LastSyncHash        string
	Metadata            FileMetadata
}

// PleaseSyncRes is used to response to client of whether file is updated or not
type PleaseSyncRes struct {
	UUID      string
	AfterPath string
	Status    string
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
	UUID      string
	AfterPath string
	Side      string
}

// PleaseFileRes is used when server response file to client (metadata)
type PleaseFileRes struct {
	UUID      string
	AfterPath string
}

type AskAllMetaReq struct {
	UUID string
}

type AskAllMetaRes struct {
	UUID         string
	SyncMetaList []SyncMetadata
}

type SyncMetadata struct { // Per file
	BeforePath          string
	AfterPath           string
	LastUpdateTimestamp uint64 // Local File changed time
	LastUpdateHash      string
	LastSyncTimestamp   uint64 // Sync Success Time
	LastSyncHash        string
}

type RescanReq struct {
	UUID          string
	RootAfterPath []string
}

type RescanRes struct {
	UUID string
}

type NeedSyncReq struct {
	UUID        string
	FileNeedPSs []FileNeedPS
}

type FileNeedPS struct {
	AfterPath string
	Event     string
}

type NeedSyncRes struct {
	UUID string
}

type NeedContentReq struct {
	UUID                string
	AfterPath           string
	LastUpdateTimestamp uint64
	LastUpdateHash      string
}

type NeedContentRes struct {
	UUID                string
	AfterPath           string
	LastUpdateTimestamp uint64
	LastUpdateHash      string
}

type Ping struct {
	UUID string
}

type RollBackReq struct {
	UUID      string
	AfterPath string
	Version   uint64
}

type RollBackRes struct {
	UUID string
}

type ShowHistoryReq struct {
	UUID        string
	AfterPath   string
	CntFromHead uint64
}

type ShowHistoryRes struct {
	History []FileHistory
}

type DownloadHistoryReq struct {
	UUID      string
	AfterPath string
	Version   uint64
}

type DownloadHistoryRes struct {
	UUID string
}

type ShareReq struct {
	UUID      string
	AfterPath string
	MaxCnt    uint64
}

type ShareRes struct {
	Link string
}

type StopShareReq struct {
	UUID string
	Link string
}

type StopShareRes struct {
	UUID string
}

type AskStagingNumReq struct {
	UUID      string
	AfterPath string
}

type AskStagingNumRes struct {
	UUID        string
	ConflictNum uint64
}

type ConflictDownloadReq struct {
	UUID      string // client UUID who want to download
	Candidate string // coflict file's UUID from FileHistory (stagingFile)
	AfterPath string // conflict file's AfterPath
}

type DisconnectRootDirReq struct {
	UUID      string
	AfterPath string
}

type DisconnectRootDirRes struct {
	UUID      string
	AfterPath string
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

func (clientRegisterRes *ClientRegisterRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(clientRegisterRes); err != nil {
		log.Println("quics: (ClientRegisterRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (clientRegisterRes *ClientRegisterRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(clientRegisterRes)
}

func (disconnectClientReq *DisconnectClientReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(disconnectClientReq); err != nil {
		log.Println("quics: (ClientDisconnectorReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (disconnectClientReq *DisconnectClientReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectClientReq)
}

func (disconnectClientRes *DisconnectClientRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(disconnectClientRes); err != nil {
		log.Println("quics: (ClientDisconnectorReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (disconnectClientRes *DisconnectClientRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectClientRes)
}

func (askRootDirReq *AskRootDirReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askRootDirReq); err != nil {
		log.Println("quics: (AskRootDirReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askRootDirReq *AskRootDirReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askRootDirReq)
}

func (askRootDirRes *AskRootDirRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askRootDirRes); err != nil {
		log.Println("quics: (AskRootDirRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askRootDirRes *AskRootDirRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askRootDirRes)
}

func (askConflictListReq *AskConflictListReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askConflictListReq); err != nil {
		log.Println("quics: (AskRootDirRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askConflictListReq *AskConflictListReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askConflictListReq)
}

func (askConflictListRes *AskConflictListRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askConflictListRes); err != nil {
		log.Println("quics: (AskRootDirRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askConflictListRes *AskConflictListRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askConflictListRes)
}

func (registerRootDirReq *RootDirRegisterReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerRootDirReq); err != nil {
		log.Println("quics: (RegisterRootDirReq.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (registerRootDirReq *RootDirRegisterReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerRootDirReq)
}

func (registerRootDirRes *RootDirRegisterRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerRootDirRes); err != nil {
		log.Println("quics: (RegisterRootDirRes.Encode) ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (registerRootDirRes *RootDirRegisterRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerRootDirRes)
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

func (askAllMetaReq *AskAllMetaReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askAllMetaReq); err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askAllMetaReq *AskAllMetaReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askAllMetaReq)
}

func (askAllMetaRes *AskAllMetaRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askAllMetaRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askAllMetaRes *AskAllMetaRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askAllMetaRes)
}

func (rescanReq *RescanReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rescanReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (rescanReq *RescanReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rescanReq)
}

func (rescanRes *RescanRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rescanRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (rescanRes *RescanRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rescanRes)
}

func (needSyncReq *NeedSyncReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(needSyncReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (needSyncReq *NeedSyncReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(needSyncReq)
}

func (needSyncRes *NeedSyncRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(needSyncRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (needSyncRes *NeedSyncRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(needSyncRes)
}

func (needContentReq *NeedContentReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(needContentReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (needContentReq *NeedContentReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(needContentReq)
}

func (needContentRes *NeedContentRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(needContentRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (needContentRes *NeedContentRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(needContentRes)
}

func (ping *Ping) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(ping); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (ping *Ping) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(ping)
}

func (rollBackReq *RollBackReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rollBackReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (rollBackReq *RollBackReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rollBackReq)
}

func (rollBackRes *RollBackRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rollBackRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (rollBackRes *RollBackRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rollBackRes)
}

func (showHistoryReq *ShowHistoryReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(showHistoryReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (showHistoryReq *ShowHistoryReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(showHistoryReq)
}

func (showHistoryRes *ShowHistoryRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(showHistoryRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (showHistoryRes *ShowHistoryRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(showHistoryRes)
}

func (downloadHistoryReq *DownloadHistoryReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(downloadHistoryReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (downloadHistoryReq *DownloadHistoryReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(downloadHistoryReq)
}

func (downloadHistoryRes *DownloadHistoryRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(downloadHistoryRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (downloadHistoryRes *DownloadHistoryRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(downloadHistoryRes)
}

func (shareReq *ShareReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(shareReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (shareReq *ShareReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(shareReq)
}

func (shareRes *ShareRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(shareRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (shareRes *ShareRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(shareRes)
}

func (stopShareReq *StopShareReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(stopShareReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (stopShareReq *StopShareReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(stopShareReq)
}

func (stopShareRes *StopShareRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(stopShareRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (stopShareRes *StopShareRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(stopShareRes)
}

func (askStagingNumReq *AskStagingNumReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askStagingNumReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askStagingNumReq *AskStagingNumReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askStagingNumReq)
}

func (askStagingNumRes *AskStagingNumRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(askStagingNumRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (askStagingNumRes *AskStagingNumRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(askStagingNumRes)
}

func (conflictDownloadReq *ConflictDownloadReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(conflictDownloadReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (conflictDownloadReq *ConflictDownloadReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(conflictDownloadReq)
}

func (disconnectRootDirReq *DisconnectRootDirReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(disconnectRootDirReq); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (disconnectRootDirReq *DisconnectRootDirReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectRootDirReq)
}

func (disconnectRootDirRes *DisconnectRootDirRes) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(disconnectRootDirRes); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (disconnectRootDirRes *DisconnectRootDirRes) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectRootDirRes)
}
