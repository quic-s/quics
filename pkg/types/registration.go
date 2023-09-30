package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Client is used to save connected client information
type Client struct {
	UUID string
	Id   uint64
	Ip   string
	Root []RootDirectory // root directory path information
}

func (client *Client) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(client); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
	}

	return buffer.Bytes()
}

func (client *Client) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(client)
}

// RootDirectory is used when registering root directory to client
type RootDirectory struct {
	Path     string // key, saved path at server (absolute path)
	Owner    string // the client that registers this root directory
	Password string // if not exist password, then the value is ""
}

func (rootDirectory *RootDirectory) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rootDirectory); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
	}

	return buffer.Bytes()
}

func (rootDirectory *RootDirectory) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rootDirectory)
}

// ClientRegisterReq is used when registering client
type ClientRegisterReq struct {
	UUID           string
	ClientPassword string
}

func (clientRegisterReq *ClientRegisterReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(clientRegisterReq)
}

// ClientDisconnectorReq is used when disconnecting client with server
type ClientDisconnectorReq struct {
	UUID     string // uuid of client
	Password string // password of server
}

func (clientDisconnectorReq *ClientDisconnectorReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(clientDisconnectorReq)
}

// RootDirRegisterReq is used when registering root directory of a client
type RootDirRegisterReq struct {
	UUID            string
	RootDirPassword string // password of the root directory
	BeforePath      string
	AfterPath       string
}

func (rootDirRegisterReq *RootDirRegisterReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rootDirRegisterReq)
}

type SyncRootDirReq struct {
	UUID            string
	RootDirPassword string // password of the root directory
	BeforePath      string
	AfterPath       string
}

func (syncRootDirReq *SyncRootDirReq) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(syncRootDirReq)
}
