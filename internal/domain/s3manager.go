package domain

type S3Interface interface {
	UploadFile(key string, fileBytes []byte, contentType string) (string, error)
	GetFileURL(key string) (string, error)
	DeleteFile(key string) error
	GetFile(key string) ([]byte, error)
	GetFileMetadata(key string) (string, error)
}
