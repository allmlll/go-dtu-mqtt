package initiallize

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"ruitong-new-service/global"
)

func MinioInit() {
	address := viper.GetString("minio.address")
	port := viper.GetString("minio.port")
	account := viper.GetString("minio.account")
	password := viper.GetString("minio.password")
	endpoint := address + ":" + port
	accessKeyID := account
	secretAccessKey := password

	var err error
	// Initialize minio client object.
	global.MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		global.Log.Error(err.Error())
	}
	global.LogInfo.Info("minio connect!")
}
