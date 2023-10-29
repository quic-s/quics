package qp

import (
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/sharing"
	"github.com/quic-s/quics/pkg/types"
)

type SharingHandler struct {
	sharingService sharing.Service
}

func NewSharingHandler(sharingService sharing.Service) *SharingHandler {
	return &SharingHandler{
		sharingService: sharingService,
	}
}

func (sh *SharingHandler) StartSharing(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	request := &types.ShareReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: ", err)
		return err
	}

	response, err := sh.sharingService.CreateLink(request)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sh *SharingHandler) StopSharing(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	request := &types.StopShareReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: ", err)
		return err
	}

	response, err := sh.sharingService.DeleteLink(request)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}
