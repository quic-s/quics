package types

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
	BeforePath string
	AfterPath  string // key
	Owner      string
	Password   string
}

// File is used to store the file's information
type File struct {
	BeforePath          string
	AfterPath           string // key
	RootDir             RootDirectory
	LatestHash          string
	LatestSyncTimestamp uint64
	ContentsExisted     bool
	Metadata            FileMetadata
}

// FileHistory is used to store the file's history
type FileHistory struct {
	Id         uint64
	Date       string
	UUID       string
	BeforePath string
	AfterPath  string
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
