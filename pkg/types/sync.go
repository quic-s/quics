package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

// File defines file sync information
type File struct {
	Path                string // key
	RootDir             RootDirectory
	LatestHash          string
	LatestSyncTimestamp uint64
}

func (file *File) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(file); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
	}

	return buffer.Bytes()
}

func (file *File) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(file)
}

type PleaseFileMetaReq struct {
	UUID      string
	AfterPath string
}

type PleaseFileMetaRes struct {
	AfterPath           string
	LatestHash          string
	LatestSyncTimestamp uint64
	ModDate             string
	UUID                string // Who fixed this file last time
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

func (pleaseSyncReq *PleaseSyncReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseSyncReq)
}

type PleaseSyncRes struct {
	UUID      string
	AfterPath string
}

// PleaseFileReq is used when client request file to server
type PleaseFileReq struct {
	UUID          string
	SyncTimestamp uint64
	BeforePath    string
	AfterPath     string
}

func (pleaseFileReq *PleaseFileReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(pleaseFileReq)
}

type PleaseTakeReq struct {
	UUID      string
	AfterPath string
}

type PleaseTakeRes struct {
	UUID      string
	AfterPath string
}

// MustSyncReq is used to inform whether file is updated or not
type MustSyncReq struct {
	LatestHash          string
	LatestSyncTimestamp uint64
	BeforePath          string
	AfterPath           string
}

func (mustSyncReq *MustSyncReq) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSyncReq); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

type MustSyncRes struct {
	UUID              string
	AfterPath         string
	LastSyncTimestamp string
	LastSyncHash      string
}

type GiveYouReq struct {
	UUID      string
	AfterPath string
}
type GiveYouRes struct {
	UUID                 string
	AfterPath            string
	LastestSyncTimestamp uint64
	LatestHash           string
}
