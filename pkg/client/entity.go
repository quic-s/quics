package client

import (
	"github.com/quic-s/quics/pkg/sync"
)

type Client struct {
	Uuid  string             `json:"uuid"`
	Id    uint64             `json:"id"`
	Ip    string             `json:"ip"`
	Root  sync.RootDirectory `json:"root_directory"` // root directory path information
	Files []sync.File        `json:"files"`          // list of synchronized files
}

type CreateClientRequest struct {
	Ip string `json:"ip"`
}

type CreateClientResponse struct {
	Uuid string `json:"uuid"`
}
