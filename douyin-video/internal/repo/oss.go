package repo

type OssRepo interface {
	PutObject(objectName string, localFileName string) error
	GetObject(objectName string, downloadedFileName string) error
	DeleteObject(objectName string) error
	// GetNameList()
	GetObjectUrl(objectName string, expiredInSec int64) (string, error)
}
