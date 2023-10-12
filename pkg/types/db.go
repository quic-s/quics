package types

// **************************************************************
// beforePath: except root directory path
// afterPath: /rootDir/..
// **************************************************************

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"
)

type DatabaseDataTypes interface {
	Client | RootDirectory | File | FileHistory | FileMetadata | Sharing
}

type DatabaseData[T DatabaseDataTypes] interface {
	Encode() []byte
	Decode(data []byte) error
}

// Client is used to save connected client information
type Client struct {
	UUID string // key
	Id   uint64
	Ip   string
	Root []RootDirectory
}

// RootDirectory is used when registering root directory to client
type RootDirectory struct {
	AfterPath  string // key
	BeforePath string
	Owner      string
	Password   string
	UUIDs      []string
}

// File is used to store the file's information
type File struct {
	AfterPath           string // key
	BeforePath          string
	RootDirKey          string
	LatestHash          string
	LatestSyncTimestamp uint64
	ContentsExisted     bool
	IsConflict          bool
	Metadata            FileMetadata
}

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

// FileMetadata retains file contents at last sync timestamp
type FileMetadata struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

// Sharing is used to store the file download information
type Sharing struct {
	Id       uint // key
	Count    uint
	MaxCount uint
	Link     string
	Owner    string
	File     FileHistory // to share file at point that client wanted in time
}

func (client *Client) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(client); err != nil {
		log.Println("quics: (Client.Encode) ", err)
	}

	return buffer.Bytes()
}

func (client *Client) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(client)
}

func (rootDirectory *RootDirectory) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(rootDirectory); err != nil {
		log.Println("quics: (RootDirectory.Encode) ", err)
	}

	return buffer.Bytes()
}

func (rootDirectory *RootDirectory) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(rootDirectory)
}

func (file *File) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(file); err != nil {
		log.Println("quics: (File.Encode) ", err)
	}

	return buffer.Bytes()
}

func (file *File) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(file)
}

func (fileHistory *FileHistory) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileHistory); err != nil {
		log.Println("quics: (FileHistory.Encode) ", err)
	}

	return buffer.Bytes()
}

func (fileHistory *FileHistory) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(fileHistory)
}

func (fileMetadata *FileMetadata) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileMetadata); err != nil {
		log.Println("quics: (FileMetadata.Encode) ", err)
	}

	return buffer.Bytes()
}

func (fileMetadata *FileMetadata) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(fileMetadata)
}

func (fileMetadata *FileMetadata) DecodeFromOSFileInfo(fileInfo os.FileInfo) {
	fileMetadata.Name = fileInfo.Name()
	fileMetadata.Size = fileInfo.Size()
	fileMetadata.Mode = fileInfo.Mode()
	fileMetadata.ModTime = fileInfo.ModTime()
	fileMetadata.IsDir = fileInfo.IsDir()
}

func (fileMetadata *FileMetadata) WriteToFile(path string) error {
	// When the file is a directory, create the directory and return.
	if fileMetadata.IsDir {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		// Set file metadata.
		err = file.Chmod(fileMetadata.Mode)
		if err != nil {
			return err
		}
		err = os.Chtimes(path, time.Now(), fileMetadata.ModTime)
		if err != nil {
			return err
		}
		return nil
	}

	// When the file is not a directory, create the file and write the file content.

	// Open file with O_TRUNC flag to overwrite the file when the file already exists.
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMetadata.Mode)
	if err != nil {
		return err
	}
	// Set file metadata.
	err = file.Chmod(fileMetadata.Mode)
	if err != nil {
		return err
	}
	err = os.Chtimes(path, time.Now(), fileMetadata.ModTime)
	if err != nil {
		return err
	}

	return nil
}

func (sharing *Sharing) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(sharing); err != nil {
		log.Println("quics: (Sharing.Encode) ", err)
	}

	return buffer.Bytes()
}

func (sharing *Sharing) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(sharing)
}
