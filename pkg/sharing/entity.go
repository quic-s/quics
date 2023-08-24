package sharing

import "github.com/quic-s/quics/pkg/history"

type Sharing struct {
	Id   uint
	Link string
	File history.FileHistory // to share file at point that client wanted in time
}
