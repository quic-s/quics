package sync

// TODO: combine path or separate path to save?

// RootDirectory connected root directory
type RootDirectory struct {
	Id       uint64 `json:"id"`
	Owner    string `json:"owner"`    // the client that registers this root directory
	Password string `json:"password"` // if not exist password, then the value is ""
	Path     string `json:"path"`
	Date     string `json:"date"`
}

// File defines file sync information
type File struct {
	Id         string        `json:"id"`
	RootDir    RootDirectory `json:"root_dir"`
	Name       string        `json:"name"`
	BeforePath string        `json:"before_path"`
	AfterPath  string        `json:"after_path"`
	LastSyncAt uint32        `json:"last_sync_at"`
}

// RegisterRootDirRequest is used to send from client to server when registering root directory of a client
type RegisterRootDirRequest struct {
	Uuid       string `json:"uuid"`
	Password   string `json:"password"`
	BeforePath string `json:"before_path"`
	AfterPath  string `json:"after_path"`
}

// ChangedFileUpdateRequest is used when updating file's changes from client to server
type ChangedFileUpdateRequest struct {
	Uuid          string `json:"uuid"`
	BeforePath    string `json:"before_path"`
	AfterPath     string `json:"after_path"`
	LastUpdatedAt uint32 `json:"last_updated_at"`
}

// ChangedFileUpdateResponse is used when updating file's changes from client to server
type ChangedFileUpdateResponse struct {
	IsConflict int    `json:"is_conflict"` // 1: conflict, 0: not conflict
	LastSyncAt uint32 `json:"last_sync_at"`
}

// ChangedFileSyncRequest is used when synchronizing file's changes from server to client
type ChangedFileSyncRequest struct {
	Uuid          string `json:"uuid"`
	BeforePath    string `json:"before_path"`
	AfterPath     string `json:"after_path"`
	LastUpdatedAt uint32 `json:"last_updated_at"`
}

// ChanedFileSyncResponse is used when synchronizing file's changes from server to client
// if last updated timestamp from client < last sync timestamp from server, then send changed file to client
type ChanedFileSyncResponse struct {
	IsConflict int    `json:"is_conflict"` // 1: conflict, 0: not conflict
	BeforePath string `json:"before_path"` // may have been changed
	AfterPath  string `json:"after_path"`  // may have been changed
	LastSyncAt uint32 `json:"last_sync_at"`
}
