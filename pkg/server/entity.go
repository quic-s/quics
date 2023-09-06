package server

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Server struct {
	Password string
}

func (server *Server) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(server); err != nil {
		log.Panicf("Error while encoding data: %s", err)
	}

	return buffer.Bytes()
}

func (server *Server) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(server)
}
