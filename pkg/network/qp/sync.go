package qp

import (
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/types"
)

type SyncHandler struct {
	syncService sync.Service
}

func NewSyncHandler(service sync.Service) *SyncHandler {
	return &SyncHandler{
		syncService: service,
	}
}

// sync root directory
// 1. (client) Open transaction
// 2. (client) Send request data for syncing root directory
// 3. (server) Receive request data
// 4. (server) Sync root directory of client to database
// TODO: 5. (server) Send response data for syncing root directory
func (sh *SyncHandler) SyncRootDir(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	var request *types.SyncRootDirReq
	if err := request.Decode(data); err != nil {
		log.Println("quics: ", err)
		return
	}

	// get root directory path of requested data
	err = sh.syncService.SyncRootDir(request)
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	// TODO: is it necessary to send response data?
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
func (sh *SyncHandler) PleaseSync(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
	log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

	// -> step 2 to 4: metadata

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	pleaseFileMetaReq := &types.PleaseFileMetaReq{}
	if err := pleaseFileMetaReq.Decode(data); err != nil {
		log.Println("quics: ", err)
		return
	}

	// ple
	// sFileMetaRes := &types.PleaseFileMetaRes{}
	// h.pleaseFileMetaRes = sh.syncServ.GetFileMetadata(pleaseFileMetaReq)

	// <- step 2 to 4: metadata

	// -> step 5 to 7: sync infomration

	// <- step 5 to 7: sync infomration

	// -> step 8 to 10: file

	// <- step 8 to 10: file

	// -> step 11: must sync transaction

	// <- step 11: must sync transaction

	// TODO: find and return certain file metadata

	data, fileInfo, fileContent, err := stream.RecvFileBMessage()
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	var pleaseSyncReq types.PleaseSyncReq
	if err := pleaseSyncReq.Decode(data); err != nil {
		log.Println("quics: ", err)
		return
	}

	// FIXME: change the condition from whether the file is exist to whether the request data is empty or not full

	err = sh.syncService.SyncFileToLatestDir(pleaseSyncReq.AfterPath, fileInfo, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	err = sh.syncService.SyncFileToHistoryDir(pleaseSyncReq.AfterPath, pleaseSyncReq.LastUpdateTimestamp, fileInfo, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	// open must sync transaction
	// openMustSyncTransaction(conn)
}
