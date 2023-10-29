package qp

import (
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/types"
)

type RegistrationHandler struct {
	registrationService registration.Service
}

func NewRegistrationHandler(service registration.Service) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: service,
	}
}

// register client
func (rh *RegistrationHandler) RegisterClient(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", types.REGISTERCLIENT, " transaction")
	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}
	request := &types.ClientRegisterReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: ", err)
		return err
	}

	// call registration service
	response, err := rh.registrationService.RegisterClient(request, conn)
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

type RegistrationAdapter struct {
	Pool *connection.Pool
}

func NewRegistrationAdapter(pool *connection.Pool) *RegistrationAdapter {
	return &RegistrationAdapter{
		Pool: pool,
	}
}

func (ra *RegistrationAdapter) UpdateClientConnection(uuid string, conn *qp.Connection) error {
	err := ra.Pool.UpdateConnection(uuid, conn)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}
	return nil
}

func (ra *RegistrationAdapter) DeleteConnection(uuid string) error {
	err := ra.Pool.DeleteConnection(uuid)
	if err != nil {
		return err
	}
	return nil
}
