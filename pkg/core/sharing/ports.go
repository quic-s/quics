package sharing

type Repository interface {
	SaveNewDownloadLink(link string) error
}

type Service interface {
	CreateLink(UUID string, afterPath string, count uint64) (string, error)
}
