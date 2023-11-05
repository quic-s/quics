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
func (sh *SyncHandler) RegisterRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	request := &types.RootDirRegisterReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// Register root directory of client to database
	response, err := sh.syncService.RegisterRootDir(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// do fullscan in goroutine
	go func() {
		_, err := sh.syncService.Rescan(&types.RescanReq{
			UUID: request.UUID,
		})
		if err != nil {
			log.Println("quics err: [RESCAN after ", transactionName, "] ", err)
			return
		}
	}()
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// sync root directory
func (sh *SyncHandler) SyncRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.RootDirRegisterReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// get root directory path of requested data
	rootDirRegisterRes, err := sh.syncService.SyncRootDir(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := rootDirRegisterRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// do fullscan in goroutine
	go func() {
		_, err := sh.syncService.Rescan(&types.RescanReq{
			UUID: request.UUID,
		})
		if err != nil {
			log.Println("quics err: [RESCAN after ", transactionName, "] ", err)
			return
		}
	}()
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// get root directory list
func (sh *SyncHandler) GetRemoteDirs(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	request := &types.AskConflictListReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	rootDirs, err := sh.syncService.GetRootDirList()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	res, err := rootDirs.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(res)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

func (sh *SyncHandler) DisconnectRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.DisconnectRootDirReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// get root directory path of requested data
	disconnectRootDirRes, err := sh.syncService.DisconnectRootDir(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := disconnectRootDirRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// please sync transaction
// it is used when client wants to sync file
func (sh *SyncHandler) PleaseSync(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	// -> return file metadata to client

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	pleaseSyncReq := &types.PleaseSyncReq{}
	if err := pleaseSyncReq.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
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
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := pleaseSyncRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// <- update file sync information before update file contents

	// -> update file contents

	data, fileInfo, fileContent, err := stream.RecvFileBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	pleaseTakeReq := &types.PleaseTakeReq{}
	if err := pleaseTakeReq.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
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
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err = pleaseTakeRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// <- update file contents
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// get conflict list transaction
// it is used when client wants to get conflict status list
func (sh *SyncHandler) AskConflictList(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")
	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.AskConflictListReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// get root directory path of requested data
	askConflictListRes, err := sh.syncService.GetConflictList(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := askConflictListRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// choose one transaction
// it is used when client wants to choose one of conflict files
func (sh *SyncHandler) ChooseOne(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.PleaseFileReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
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
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := pleaseFileRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// rescan transaction
// it is used when client wants to rescan (fullscan)
func (sh *SyncHandler) Rescan(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.RescanReq{}
	if err := request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	rescanRes, err := sh.syncService.Rescan(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := rescanRes.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}
	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// rollback transaction
// it is used when client wants to rollback file to specific version
func (sh *SyncHandler) RollbackFileByHistory(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.RollBackReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := sh.syncService.RollbackFileByHistory(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// conflict download transaction
// it is used when client wants to download conflict files
func (sh *SyncHandler) ConflictDownload(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")
	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.AskStagingNumReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, err := sh.syncService.GetStagingNum(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendBMessage(data)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	// if count is zero, then close transaction
	if response.ConflictNum == 0 {
		return nil
	}

	requests, err := sh.syncService.GetConflictFiles(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	for _, request := range requests {
		data, err = request.Encode()
		if err != nil {
			log.Println("quics err: [", transactionName, "] ", err)
			return err
		}

		if request.Candidate == "server" {
			filePath := utils.GetQuicsSyncDirPath() + request.AfterPath
			err = stream.SendFileBMessage(data, filePath)
			if err != nil {
				log.Println("quics err: [", transactionName, "] ", err)
				return err
			}
		} else {
			filePath := utils.GetConflictFileNameByAfterPath(request.AfterPath, request.Candidate)
			err = stream.SendFileBMessage(data, filePath)
			if err != nil {
				log.Println("quics err: [", transactionName, "] ", err)
				return err
			}
		}
	}

	log.Println("quics: [", transactionName, "] transaction finished")
	return nil
}

// download history transaction
// it is used when client wants to download specific history(version) file
func (sh *SyncHandler) DownloadHistory(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: receive ", transactionName, " transaction")
	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	request := &types.DownloadHistoryReq{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	response, filePath, err := sh.syncService.DownloadHistory(request)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	data, err = response.Encode()
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	err = stream.SendFileBMessage(data, filePath)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return err
	}

	log.Println("quics: [", transactionName, "] transaction finished")
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
	transactionName string
	wg              *stdsync.WaitGroup
	stream          *qp.Stream
}

// OpenTransaction opens transaction
// this method is called when server wants to open transaction (server-push)
func (sa *SyncAdapter) OpenTransaction(transactionName string, uuid string) (sync.Transaction, error) {
	log.Println("quics: open ", transactionName, " transaction")
	// get connection from pool by uuid
	conn, err := sa.Pool.GetConnection(uuid)
	if err != nil {
		log.Println("quics err: [", transactionName, "] ", err)
		return nil, err
	}

	transaction := &Transaction{
		transactionName: transactionName,
		wg:              &stdsync.WaitGroup{},
		stream:          nil,
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
		log.Println("quics err: [", transactionName, "] ", err)
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) Close() error {
	t.wg.Done()
	log.Println("quics: close ", t.transactionName, " transaction")
	return nil
}

// send and receive mustsync request and response
// using on CallMustSync method in sync service
func (t *Transaction) RequestMustSync(mustSyncReq *types.MustSyncReq) (*types.MustSyncRes, error) {
	request, err := mustSyncReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		return nil, err
	}

	mustSyncRes := &types.MustSyncRes{}
	if err := mustSyncRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}
	return mustSyncRes, nil
}

// send and receive giveyou request and response
// using on CallMustSync method in sync service
func (t *Transaction) RequestGiveYou(giveYouReq *types.GiveYouReq, historyFilePath string) (*types.GiveYouRes, error) {
	request, err := giveYouReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// send (history file)
	err = t.stream.SendFileBMessage(request, historyFilePath)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	giveYouRes := &types.GiveYouRes{}
	if err := giveYouRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}
	return giveYouRes, nil
}

// send and receive forcesync request and response
// using on CallForceSync method in sync service
func (t *Transaction) RequestForceSync(mustSyncReq *types.MustSyncReq, historyFilePath string) (*types.MustSyncRes, error) {
	request, err := mustSyncReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// send (history file)
	err = t.stream.SendFileBMessage(request, historyFilePath)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	mustSyncRes := &types.MustSyncRes{}
	if err := mustSyncRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}
	return mustSyncRes, nil
}

// send and receive askallmeta request and response
// using on FullScan method in sync service
func (t *Transaction) RequestAskAllMeta(askAllMetaReq *types.AskAllMetaReq) (*types.AskAllMetaRes, error) {
	request, err := askAllMetaReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	askAllMetaRes := &types.AskAllMetaRes{}
	if err := askAllMetaRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}
	return askAllMetaRes, nil
}

// send and receive needsync request and response
// using when server wants to get PleaseSync request from client
func (t *Transaction) RequestNeedSync(needSyncReq *types.NeedSyncReq) (*types.NeedSyncRes, error) {
	request, err := needSyncReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	// receive
	res, err := t.stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}

	needSyncRes := &types.NeedSyncRes{}
	if err := needSyncRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, err
	}
	return needSyncRes, nil
}

// send needcontent request and receive file metadata and file contents
// using on FullScan method in sync service when file contents are needed
func (t *Transaction) RequestNeedContent(needContentReq *types.NeedContentReq) (*types.NeedContentRes, *types.FileMetadata, io.Reader, error) {
	request, err := needContentReq.Encode()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, nil, nil, err
	}

	err = t.stream.SendBMessage(request)
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, nil, nil, err
	}

	// receive
	res, fileInfo, content, err := t.stream.RecvFileBMessage()
	if err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
		return nil, nil, nil, err
	}

	needContentRes := &types.NeedContentRes{}
	if err := needContentRes.Decode(res); err != nil {
		log.Println("quics err: [", t.transactionName, "] ", err)
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
