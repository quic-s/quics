package conflict

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TwoOptions is used when resolving file conflicts at client side
type TwoOptions struct {
	ServerSideHash          string
	ServerSideSyncTimestamp uint64
	ClientSideHash          string
	ClientSideTimestamp     uint64
}

func (twoOptions *TwoOptions) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(twoOptions); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
