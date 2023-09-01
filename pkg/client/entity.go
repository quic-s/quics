package client

import (
	"bytes"
	"encoding/gob"
	"github.com/quic-s/quics/pkg/registration"
	"github.com/quic-s/quics/pkg/sync"
	"log"
)

// Client is used to save connected client information
type Client struct {
	Uuid  string
	Id    uint64
	Ip    string
	Root  []registration.RootDirectory // root directory path information
	Files []sync.File                  // list of synchronized files
}

// RegisterClientRequest is used when registering client
type RegisterClientRequest struct {
	RequestId uint64
	Ip        string
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
