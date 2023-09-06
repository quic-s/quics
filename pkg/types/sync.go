package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

// File defines file sync information
type File struct {
	Path                string // key
	LatestHash          string
	LatestSyncTimestamp uint64
}

// PleaseSync is used when updating file's changes from client to server
type PleaseSync struct {
	Uuid                 string
	BeforePath           string
	AfterPath            string
	LastUpdatedTimestamp uint64
	LastUpdateHash       uint64
}

func (pleaseSync *PleaseSync) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSync)
}

// PleaseFile is used when client request file to server
type PleaseFile struct {
	SyncTimestamp uint64
}

func (pleaseFile *PleaseFile) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFile)
}

// MustSyncMessage is used to inform whether file is updated or not
type MustSyncMessage struct {
	LatestSyncTimestamp uint64
	BeforePath          string
	AfterPath           string
}

func (mustSyncMessage *MustSyncMessage) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSyncMessage); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

// MustSyncFileWithMessage is used when synchronizing file's changes from server to client
// request, but server -> client
type MustSyncFileWithMessage struct {
	LatestHash          string
	LatestSyncTimestamp uint64
	BeforePath          string
	AfterPath           string
}

func (mustSyncFileWithMessage *MustSyncFileWithMessage) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSyncFileWithMessage); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
