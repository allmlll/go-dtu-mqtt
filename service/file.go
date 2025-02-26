package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"ruitong-new-service/global"

	"mime/multipart"
	"slices"
	"strings"
	"time"
)

type fileService struct {
	whitelist []string
}

var File = fileService{
	whitelist: []string{"jpg", "jpeg", "png", "gif", "bmp", "tif", "tiff", "webp", "svg"},
}

func (f *fileService) getType(image *multipart.FileHeader) (string, error) {
	name := image.Filename
	split := strings.Split(name, ".")
	typeString := split[len(split)-1]
	if !slices.Contains(f.whitelist, typeString) {
		return "", errors.New("上传的不是图片")
	}
	return typeString, nil
}

func (f *fileService) Upload(image *multipart.FileHeader) (any, error) {
	var typeString string
	var err error
	if typeString, err = f.getType(image); err != nil {
		return nil, err
	}

	open, err := image.Open()
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	// 生成uuid
	u := uuid.NewV1()
	u = uuid.NewV5(u, image.Filename)
	bucketName := viper.GetString("minio.bucketName")
	address := viper.GetString("minio.address")
	port := viper.GetString("minio.port")
	_, err = global.MinioClient.PutObject(context.TODO(), bucketName, u.String(), open, -1, minio.PutObjectOptions{ContentType: "image/" + typeString})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return gin.H{
		"url": fmt.Sprintf("http://%v:%v/%v/%v", address, port, bucketName, u),
	}, nil
}

func (f *fileService) CheckApk() (any, error) {
	objects := global.MinioClient.ListObjects(context.TODO(), "apk", minio.ListObjectsOptions{})
	latestName := ""
	var latestTime time.Time
	for object := range objects {
		if object.Err != nil {
			global.Log.Error(object.Err.Error())
			continue
		}
		if latestName == "" || object.LastModified.After(latestTime) {
			latestName = object.Key
			latestTime = object.LastModified
		}
	}
	if latestName == "" {
		return nil, errors.New("最新文件不存在")
	}
	name := strings.TrimSuffix(latestName, ".apk")
	name = strings.ReplaceAll(name, ".", "")
	address := viper.GetString("minio.address")
	port := viper.GetString("minio.port")
	return gin.H{
		"version": name,
		"url":     fmt.Sprintf("http://%v:%v/apk/%v", address, port, latestName),
	}, nil
}
