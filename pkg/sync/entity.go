package sync

// RootDirectory connected root directory
type RootDirectory struct {
	Path string `json:"path"`
	Date string `json:"date"`
}

// File
// Define file data which has file metadata (saved to database)
type File struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}
