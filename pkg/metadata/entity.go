package metadata

// FileMetadata
// Metadata format of each file (not saved to database)
type FileMetadata struct {
	Id         string
	Version    string
	Name       string
	Format     string
	Size       uint64
	Auth       string
	Owner      string
	LastEditor string
	CreatedAt  string
	ModifiedAt string
}
