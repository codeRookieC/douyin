package oss

import (
	"fmt"
	"os"
	"video_server/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssDao struct {
	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucketName      string
	codedir         string
	client          *oss.Client
	bucket          *oss.Bucket
}

var client *oss.Client
var bucket *oss.Bucket

func init() {
	c, err := oss.New(config.C.Oss.Endpoint, config.C.Oss.AccessKeyId, config.C.Oss.AccessKeySecret)
	if err != nil {
		handleError(err)
	}
	client = c
	// 获取存储空间。
	b, err := client.Bucket(config.C.Oss.BucketName)
	if err != nil {
		handleError(err)
	}
	bucket = b
}
func NewOssDao() *OssDao {
	return &OssDao{
		endpoint:        config.C.Oss.Endpoint,
		accessKeyId:     config.C.Oss.AccessKeyId,
		accessKeySecret: config.C.Oss.AccessKeySecret,
		bucketName:      config.C.Oss.BucketName,
		codedir:         config.C.Oss.Codedir,
		client:          client,
		bucket:          bucket,
	}
}
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
func (o *OssDao) PutObject(objectName string, localFileName string) error {
	// // 创建OSSClient实例。
	// client, err := oss.New(o.endpoint, o.accessKeyId, o.accessKeySecret)
	// if err != nil {
	// 	handleError(err)
	// 	return err
	// }
	// // 获取存储空间。
	// bucket, err := client.Bucket(o.bucketName)
	// if err != nil {
	// 	handleError(err)
	// 	return err
	// }
	// 上传文件。
	err := o.bucket.PutObjectFromFile(o.codedir+objectName, localFileName)
	if err != nil {
		handleError(err)
		return err
	}
	return err
}

func (o *OssDao) GetObject(objectName string, downloadedFileName string) error {

	// 下载文件。
	err := o.bucket.GetObjectToFile(o.codedir+objectName, downloadedFileName)
	if err != nil {
		handleError(err)
		return err
	}
	return err
}
func (o *OssDao) DeleteObject(objectName string) error {
	// 删除文件。
	err := o.bucket.DeleteObject(o.codedir + objectName)

	if err != nil {
		handleError(err)
		return err
	}
	return err
}
func (o *OssDao) GetObjectUrl(objectName string, expiredInSec int64) (string, error) {

	// 生成用于下载的签名URL，并指定签名URL的有效时间为 expiredInSec 秒。
	signedURL, err := o.bucket.SignURL(o.codedir+objectName, oss.HTTPGet, expiredInSec)
	if err != nil {
		handleError(err)
		return signedURL, err
	}
	return signedURL, err
}

// GetNameList()
