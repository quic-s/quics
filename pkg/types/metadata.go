package types

// FileMetadata retains file contents at last sync timestamp
type FileMetadata struct {
	Id         string
	Hash       string // TODO: does it need?
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
