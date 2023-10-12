package qp

import (
	"crypto/sha1"
	"log"
	stdsync "sync"

	"github.com/quic-s/quics/pkg/network/qp/connection"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
)

type SyncHandler struct {
	psMut       map[byte]*stdsync.Mutex
	syncService sync.Service
}

func NewSyncHandler(service sync.Service) *SyncHandler {
	psMut := make(map[byte]*stdsync.Mutex)

	for i := uint8(0); i < 16; i++ {
		psMut[i] = &stdsync.Mutex{}
	}
	return &SyncHandler{
		psMut:       psMut,
		syncService: service,
	}
}

// sync root directory
// 1. (client) Open transaction
// 2. (client) Send request data for syncing root directory
// 3. (server) Receive request data
// 4. (server) Sync root directory of client to database
// TODO: 5. (server) Send response data for syncing root directory
func (sh *SyncHandler) SyncRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.SyncRootDirReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// get root directory path of requested data
	err = sh.syncService.SyncRootDir(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// TODO: is it necessary to send response data?
	return nil
}

// please sync transaction
// 1. (client) Open transaction
// 2. (client) PleaseFileMetaReq for getting a file metadata
// 3. (server) Find and return certain file metadata
// 4. (server) PleaseFileMetaRes for returning a file metadata
// 5. (client) PleaseSyncReq if file update is available
// 6. (server) Update the history with file metadata and set flag 'ContentsExisted' = false
// 7. (server) PleaseSyncRes
// 8. (client) PleaseTakeReq for sync a file
// 9. (server) Get file contents and set flag 'ContentsExisted' = true
// 10. (server) PleaseTakeRes
// 11. (server) Go to the MustSync transaction
func (sh *SyncHandler) PleaseSync(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received : PleaseSync", conn.Conn.RemoteAddr().String())

	// -> return file metadata to client

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	pleaseSyncReq := &types.PleaseSyncReq{}
	if err := pleaseSyncReq.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// lock mutex by hash value of file path
	// using hash value is to reduce the number of mutex
	h := sha1.New()
	h.Write([]byte(pleaseSyncReq.AfterPath))
	hash := h.Sum(nil)

	sh.psMut[uint8(hash[0]%16)].Lock()
	defer sh.psMut[uint8(hash[0]%16)].Unlock()

	pleaseSyncRes, err := sh.syncService.UpdateFileWithoutContents(pleaseSyncReq)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := pleaseSyncRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// <- update file sync information before update file contents

	// -> update file contents

	data, fileInfo, fileContent, err := stream.RecvFileBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	pleaseTakeReq := &types.PleaseTakeReq{}
	if err := pleaseTakeReq.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	pleaseTakeRes, err := sh.syncService.UpdateFileWithContents(pleaseTakeReq, fileInfo, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err = pleaseTakeRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// <- update file contents
	return nil
}

type SyncAdapter struct {
	Pool *connection.Pool
}

func NewSyncAdapter(pool *connection.Pool) *SyncAdapter {
	return &SyncAdapter{
		Pool: pool,
	}
}

type Transaction struct {
	wg     *stdsync.WaitGroup
	stream *qp.Stream
}

// must sync transaction
// 1. (server) Open transaction
// 2. (server) MustSyncReq with file metadata to all registered clients without where the file come from
// 3. (client) MustSyncRes if file update is available
// 3-1. (server) If all request data are exist, then go to step 4
// 3-2. (server) If not, then this transaction should be closed
// 4. (server) GiveYouReq for giving file contents
// 5. (client) GiveYouRes
func (sa *SyncAdapter) OpenMustSyncTransaction(uuid string) (sync.Transaction, error) {
	// get connection from pool by uuid
	conn, err := sa.Pool.GetConnection(uuid)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		wg:     &stdsync.WaitGroup{},
		stream: nil,
	}

	// make error channel to receive error from goroutine
	errChan := make(chan error)
	go func() {
		err := conn.OpenTransaction(types.MUSTSYNC, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
			// add wait group to wait for closing transaction
			transaction.wg.Add(1)

			// set stream to transaction
			transaction.stream = stream

			// send nil to error channel
			// this would be followed after set stream
			errChan <- nil

			transaction.wg.Wait()
			return nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// wait for setting stream
	err = <-errChan
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) RequestMustSync(mustSyncReq *types.MustSyncReq) (*types.MustSyncRes, error) {

	request, err := mustSyncReq.Encode()
	if err != nil {
		return nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		return nil, err
	}

	mustSyncRes := &types.MustSyncRes{}
	if err := mustSyncRes.Decode(res); err != nil {
		return nil, err
	}
	return mustSyncRes, nil
}

func (t *Transaction) RequestGiveYou(giveYouReq *types.GiveYouReq, historyFilePath string) (*types.GiveYouRes, error) {
	request, err := giveYouReq.Encode()
	if err != nil {
		return nil, err
	}

	// send (history file)
	err = t.stream.SendFileBMessage(request, historyFilePath)
	if err != nil {
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		return nil, err
	}

	giveYouRes := &types.GiveYouRes{}
	if err := giveYouRes.Decode(res); err != nil {
		return nil, err
	}
	return giveYouRes, nil
}

func (t *Transaction) Close() error {
	t.wg.Done()
	return nil
}
