package qp

import (
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/sharing"
	"github.com/quic-s/quics/pkg/types"
)

type HistoryHandler struct {
	historyService history.Service
	sharingService sharing.Service
}

func NewHistoryHandler(service history.Service, sharingService sharing.Service) *HistoryHandler {
	return &HistoryHandler{
		historyService: service,
		sharingService: sharingService,
	}
}

func (hh *HistoryHandler) ShowHistory(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] receive bmessage: ", err)
		return err
	}

	request := &types.ShowHistoryReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] decode request: ", err)
		return err
	}

	response, err := hh.historyService.ShowHistory(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] while historyService: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] encode response: ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: [", transactionName, "] send bmessage: ", err)
		return err
	}

	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}
