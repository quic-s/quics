package common

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Response is used when sending only response success/error message from server to client
type Response struct {
	RequestId uint64
	Message   string
}

func (response *Response) Encode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(response); err != nil {
		log.Panicf("Error while encoding request data: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
