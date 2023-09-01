package sync

import (
	"bytes"
	"encoding/gob"
)

// File defines file sync information
type File struct {
	Path                string // key
	LatestHash          string
	LatestSyncTimestamp uint64
	ClientData          map[string]SyncedClient // key of map: client uuid
}

// SyncedClient saves client's last synced timestamp and hash at that time
type SyncedClient struct {
	LastHash          string
	LastSyncTimestamp uint64
}

// PleaseSyncRequest is used when updating file's changes from client to server
type PleaseSyncRequest struct {
	RequestId uint64
	Uuid      string

	// e.g., /home/ubuntu/rootDir/file
	BeforePath             string // /home/ubuntu
	AfterPath              string // /rootDir/file
	LatestUpdatedTimestamp uint64

	LatestSyncTimestamp uint64 // to save this timestamp data to server database
}

func (pleaseSyncRequest *PleaseSyncRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSyncRequest)
}

// PleaseSyncResponse is used when updating file's changes from client to server (not conflict)
type PleaseSyncResponse struct {
	RequestId           uint64
	LatestSyncTimestamp uint64
}

func (pleaseSyncResponse *PleaseSyncResponse) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSyncResponse)
}

// MustSyncRequest is used when synchronizing file's changes from server to client
// request, but server -> client
type MustSyncRequest struct {
	RequestId             uint64
	LatestHash            string // depends on server
	LatestSyncTimestamp   uint64 // depends on server
	PreviousHash          string // depends on client
	PreviousSyncTimestamp uint64 // depends on client
	BeforePath            string
	AfterPath             string
}

// MustSyncResponse is used when synchronizing file's changes from server to client
// response, but client -> server
type MustSyncResponse struct {
	RequestId              uint64
	Uuid                   string
	BeforePath             string
	AfterPath              string
	LatestUpdatedTimestamp uint64
	LatestSyncTimestamp    uint64
}
