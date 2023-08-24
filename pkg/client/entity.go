package client

import (
	"github.com/quic-s/quics/pkg/metadata"
	"github.com/quic-s/quics/pkg/sync"
)

// Client
// Information of connected client
type Client struct {
	Id    string
	Ip    string
	Root  sync.RootDirectory // root directory path information
	Files []metadata.File    // list of synchronized files
}
