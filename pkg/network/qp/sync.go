package qp

import (
	"crypto/sha1"
	"io"
	"log"
	stdsync "sync"

	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/utils"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
)

type SyncHandler struct {
	lockNum     uint8
	pathMut     map[byte]*stdsync.Mutex
	syncService sync.Service
}

func NewSyncHandler(service sync.Service) *SyncHandler {
	lockNum := uint8(32)
	pathMut := map[byte]*stdsync.Mutex{}

	for i := uint8(0); i < lockNum; i++ {
		pathMut[i] = &stdsync.Mutex{}
	}
	return &SyncHandler{
		lockNum:     lockNum,
		pathMut:     pathMut,
		syncService: service,
	}
}

// register root directory
// 1. (client) Open transaction
// 2. (client) Send request data for registering root directory
// 3. (server) Receive request data
// 4. (server) Register root directory of client to database
// TODO: 5. (server) Send response data for registering root directory
func (sh *SyncHandler) RegisterRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	request := &types.RootDirRegisterReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// Register root directory of client to database
	response, err := sh.syncService.RegisterRootDir(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// do fullscan in goroutine
	go func() {
		_, err := sh.syncService.Rescan(&types.RescanReq{
			UUID: request.UUID,
		})
		if err != nil {
			log.Println("quics: ", err)
			return
		}
	}()
	return nil
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

	request := &types.RootDirRegisterReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// get root directory path of requested data
	rootDirRegisterRes, err := sh.syncService.SyncRootDir(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := rootDirRegisterRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// do fullscan in goroutine
	go func() {
		_, err := sh.syncService.Rescan(&types.RescanReq{
			UUID: request.UUID,
		})
		if err != nil {
			log.Println("quics: ", err)
			return
		}
	}()
	return nil
}

// get root directory list
// 1. (client) Open transaction
// 2. (client) Send request for getting root directory list
// 3. (server) Receive request data
// 4. (server) Get root directory list from database
// 5. (server) Send response data for getting root directory list
func (sh *SyncHandler) GetRemoteDirs(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	request := &types.AskConflictListReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	rootDirs, err := sh.syncService.GetRootDirList()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	res, err := rootDirs.Encode()
	if err != nil {
		return err
	}

	err = stream.SendBMessage(res)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
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

	sh.pathMut[uint8(hash[0]%sh.lockNum)].Lock()
	defer sh.pathMut[uint8(hash[0]%sh.lockNum)].Unlock()

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
	fileMetedata := &types.FileMetadata{
		Name:    fileInfo.Name,
		Size:    fileInfo.Size,
		Mode:    fileInfo.Mode,
		ModTime: fileInfo.ModTime,
		IsDir:   fileInfo.IsDir,
	}

	pleaseTakeRes, err := sh.syncService.UpdateFileWithContents(pleaseTakeReq, fileMetedata, fileContent)
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

func (sh *SyncHandler) AskConflictList(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: AskConflicList received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.AskConflictListReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// get root directory path of requested data
	askConflictListRes, err := sh.syncService.GetConflictList(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := askConflictListRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	return nil
}

func (sh *SyncHandler) ChooseOne(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: ChooseOne received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.PleaseFileReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	// lock mutex by hash value of file path
	// using hash value is to reduce the number of mutex
	h := sha1.New()
	h.Write([]byte(request.AfterPath))
	hash := h.Sum(nil)

	sh.pathMut[uint8(hash[0]%sh.lockNum)].Lock()
	defer sh.pathMut[uint8(hash[0]%sh.lockNum)].Unlock()

	// get root directory path of requested data
	pleaseFileRes, err := sh.syncService.ChooseOne(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := pleaseFileRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	return nil
}

func (sh *SyncHandler) Rescan(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: Rescan received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.RescanReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	rescanRes, err := sh.syncService.Rescan(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := rescanRes.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
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

func (sa *SyncAdapter) OpenTransaction(transactionName string, uuid string) (sync.Transaction, error) {
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
		err := conn.OpenTransaction(transactionName, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
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

func (t *Transaction) RequestAskAllMeta(askAllMetaReq *types.AskAllMetaReq) (*types.AskAllMetaRes, error) {
	request, err := askAllMetaReq.Encode()
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

	askAllMetaRes := &types.AskAllMetaRes{}
	if err := askAllMetaRes.Decode(res); err != nil {
		return nil, err
	}
	return askAllMetaRes, nil
}

func (t *Transaction) RequestNeedSync(needSyncReq *types.NeedSyncReq) (*types.NeedSyncRes, error) {
	request, err := needSyncReq.Encode()
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

	needSyncRes := &types.NeedSyncRes{}
	if err := needSyncRes.Decode(res); err != nil {
		return nil, err
	}
	return needSyncRes, nil
}

func (t *Transaction) Close() error {
	t.wg.Done()
	return nil
}

func (t *Transaction) RequestNeedContent(needContentReq *types.NeedContentReq) (*types.NeedContentRes, *types.FileMetadata, io.Reader, error) {
	request, err := needContentReq.Encode()
	if err != nil {
		return nil, nil, nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		return nil, nil, nil, err
	}

	// receive
	res, fileInfo, content, err := t.stream.RecvFileBMessage()
	if err != nil {
		return nil, nil, nil, err
	}

	needContentRes := &types.NeedContentRes{}
	if err := needContentRes.Decode(res); err != nil {
		return nil, nil, nil, err
	}

	fileMetadata := &types.FileMetadata{
		Name:    fileInfo.Name,
		Size:    fileInfo.Size,
		Mode:    fileInfo.Mode,
		ModTime: fileInfo.ModTime,
		IsDir:   fileInfo.IsDir,
	}
	return needContentRes, fileMetadata, content, nil
}

func (sh *SyncHandler) RollbackFileByHistory(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received: ", conn.Conn.RemoteAddr())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.RollBackReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := sh.syncService.RollbackFileByHistory(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (sh *SyncHandler) ConflictDownload(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received: ", conn.Conn.RemoteAddr())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.AskStagingNumReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, err := sh.syncService.GetStagingNum(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	// if count is zero, then close transaction
	if response.ConflictNum == 0 {
		return nil
	}

	requests, err := sh.syncService.GetConflictFiles(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	for _, request := range requests {
		data, err = request.Encode()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		filePath := utils.GetConflictFileNameByAfterPath(request.AfterPath, request.UUID)
		err = stream.SendFileBMessage(data, filePath)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
	}

	return nil
}

func (sh *SyncHandler) DownloadHistory(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: message received: ", conn.Conn.RemoteAddr())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	request := &types.DownloadHistoryReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return err
	}

	response, filePath, err := sh.syncService.DownloadHistory(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	err = stream.SendFileBMessage(data, filePath)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}
