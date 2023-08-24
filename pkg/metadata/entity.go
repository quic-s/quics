package metadata

// File file metadata
// TODO: Must to add packet information of each file for transportation
type File struct {
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
