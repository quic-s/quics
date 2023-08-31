package client

import (
	"github.com/quic-s/quics/pkg/sync"
)

type Client struct {
	Uuid  string               `json:"uuid"`
	Id    uint64               `json:"id"`
	Ip    string               `json:"ip"`
	Root  []sync.RootDirectory `json:"root_directory"` // root directory path information
	Files []sync.File          `json:"files"`          // list of synchronized files
}

// RegisterClientRequest is used to send from client to server when registering client
type RegisterClientRequest struct {
	Ip string `json:"ip"`
}

// RegisterClientResponse is used to send from server to client when registering client
type RegisterClientResponse struct {
	Uuid string `json:"uuid"`
}

// DisconnectClientRequest is used when disconnecting client with server
type DisconnectClientRequest struct {
	Uuid     string `json:"uuid"`
	Password string `json:"password"`
}
