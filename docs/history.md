# History

This document describes the history and how to handle it.

## What is History?

History means that as a file changes in constant synchronization, the server stores the state of each change. So for each logical clock, a file will be stored on the server with a timestamp as version.

## When is History created?

History is saved whenever changes are made to a file and it is passed from the client to the server. 

This means that the files are synchronized normally without problems such as conflict. The server saves information about each file in the database and saves the files in the history directory.

### History Directory

For each version (timestamp), all files for that version are stored in the `<syncRootDir>.history` directory. At this time, the file is saved with the name `filename_<timestamp>` to distinguish between versions.

### History Data

```go
// FileHistory is used to store the file's history
type FileHistory struct {
	AfterPath  string // key
	BeforePath string
	Date       string
	UUID       string
	Timestamp  uint64
	Hash       string
	File       FileMetadata // must have file metadata at the point that client wanted in time
}
```

The history data is a struct that stores in the database. It contains the path of the file and the data of the file that is stored in the history directory.

## History Management

### History Lookup

Clients can get history information about files that are synchronized to the server using the qic history show command. 

In this case, the history information is stored on the server, so it will request the information from the server and output it. 

Clients can also download and view past versions of files stored on the server. 

### History Rollback

Since the history information is stored on the server, the client can request the server to roll back to a specific version.

In this case, the server increments the timestamp of the file and sends the file to the client as a MUSTSYNC transaction. The client receives the file and replaces the file with the received file.

