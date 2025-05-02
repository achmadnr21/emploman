package domain

type S3Repository interface {
	UploadFile(filePath string, folderPath string) (string, error)
	DownloadFile(objectKey string, folderPath string) (string, error)
	DeleteFile(objectKey string) error
	ListFiles(folderPath string) ([]string, error)
	GetFileURL(objectKey string) (string, error)
	GetFileMetadata(objectKey string) (map[string]string, error)
}
