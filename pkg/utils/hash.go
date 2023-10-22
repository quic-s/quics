package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/quic-s/quics/pkg/types"
)

func MakeHashFromFileMetadata(afterPath string, info *types.FileMetadata) string {
	h := sha512.New()
	h.Write([]byte(afterPath)) // /root/*
	h.Write([]byte(info.ModTime.UTC().String()))
	h.Write([]byte(info.Mode.String()))
	h.Write([]byte(fmt.Sprint(info.Size)))
	return hex.EncodeToString(h.Sum(nil))
}
