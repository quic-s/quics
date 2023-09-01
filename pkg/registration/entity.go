package registration

import (
	"bytes"
	"encoding/gob"
)

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
