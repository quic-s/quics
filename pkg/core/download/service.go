package download

type MyDownloadService struct {
	downloadService Repository
}

func NewDownloadService(downloadService Repository) *MyDownloadService {
	return &MyDownloadService{
		downloadService: downloadService,
	}
}
