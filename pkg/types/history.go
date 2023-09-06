package types

type FileHistory struct {
	Id   uint64
	Date string
	Uuid string
	Path string       // path of stored the file with history
	File FileMetadata // must have file metadata at the point that client wanted in time
}
