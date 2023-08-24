package sync

// RootDirectory connected root directory
type RootDirectory struct {
	Path string
	Date string
}

// File
// Define file data which has file metadata (saved to database)
type File struct {
	Id   string
	Name string
	Path string
}
