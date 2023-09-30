package types

import (
	"bytes"
	"encoding/gob"
)

type TwoOptions struct {
	BeforePath              string
	AfterPath               string
	ServerSideHash          string
	ServerSideSyncTimestamp uint64
	ClientSideHash          string
	ClientSideTimestamp     uint64
}

type ChosenOne struct {
	BeforePath          string
	AfterPath           string
	ChosenHash          string
	ChosenTimestamp     uint64
	LastUpdateHash      string //new
	LastUpdateTimestamp uint64 // chosenTimestamp +1
}

func (twoOptions *TwoOptions) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(twoOptions); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (twoOptions *TwoOptions) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(twoOptions); err != nil {
		panic(err)
	}

}

func (chosenOne *ChosenOne) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(chosenOne); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (chosenOne *ChosenOne) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(chosenOne); err != nil {
		panic(err)
	}

}
