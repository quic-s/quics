package register

import (
	"bytes"
	"encoding/gob"
	"github.com/quic-s/quics/pkg/sync"
	"log"
)

// Client is used to save connected client information
type Client struct {
	Uuid  string
	Id    uint64
	Ip    string
	Root  []RootDirectory // root directory path information
	Files []sync.File     // list of synchronized files
}

// RegisterClientRequest is used when registering client
type RegisterClientRequest struct {
	Ip string
}

// Decode decodes message data from client through protocol to struct
func (registerClientRequest *RegisterClientRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerClientRequest)
}

// RegisterClientResponse is used when registering client
type RegisterClientResponse struct {
	RequestId uint64
	Uuid      string
}

// Encode encodes struct for sending to client through protocol
func (registerClientResponse *RegisterClientResponse) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerClientResponse); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DisconnectClientRequest is used when disconnecting client with server
type DisconnectClientRequest struct {
	RequestId uint64
	Uuid      string
	Password  string // password of server
}

func (disconnectClientRequest *DisconnectClientRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(disconnectClientRequest)
}

// RootDirectory is used when registering root directory to client
type RootDirectory struct {
	Id       uint64
	Owner    string // the client that registers this root directory
	Password string // if not exist password, then the value is ""
	Path     string
}

// RegisterRootDirRequest is used when registering root directory of a client
type RegisterRootDirRequest struct {
	RequestId uint64
	Uuid      string
	Password  string // password of the root directory

	// e.g., /home/ubuntu/rootDir/*
	BeforePath string // /home/ubuntu
	AfterPath  string // /rootDir/*
}

func (registerRootDirRequest *RegisterRootDirRequest) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(registerRootDirRequest)
}
