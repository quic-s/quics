package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Client is used to save connected client information
type Client struct {
	Uuid  string
	Id    uint64
	Ip    string
	Root  []RootDirectory // root directory path information
	Files []File          // list of synchronized files
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
	Path     string // saved path at server
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

// RegisterClientRequest is used when registering client
type RegisterClientRequest struct {
	Uuid           string
	ClientPassword string
}

func (registerClientRequest *RegisterClientRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerClientRequest)
}

// DisconnectClientRequest is used when disconnecting client with server
type DisconnectClientRequest struct {
	Password string // password of server
}

func (disconnectClientRequest *DisconnectClientRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectClientRequest)
}

// RegisterRootDirRequest is used when registering root directory of a client
type RegisterRootDirRequest struct {
	Uuid            string
	RootDirPassword string // password of the root directory
	BeforePath      string
	AfterPath       string
}

func (registerRootDirRequest *RegisterRootDirRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerRootDirRequest)
}

type SyncRootDirRequest struct {
	Uuid            string
	RootDirPassword string // password of the root directory
	BeforePath      string
	AfterPath       string
}

func (syncRootDirRequest *SyncRootDirRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(syncRootDirRequest)
}
