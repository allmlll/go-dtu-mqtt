package initiallize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"ruitong-new-service/global"
)

func RedisInit() {
	if global.Redis == nil {
		address := viper.GetString("redis.address")
		port := viper.GetString("redis.port")
		password := viper.GetString("redis.password")
		dB := viper.GetInt("redis.DB")
		poolSize := viper.GetInt("redis.PoolSize")

		global.Redis = redis.NewClient(&redis.Options{
			Addr:     address + ":" + port,
			Password: password,
			DB:       dB,
			PoolSize: poolSize,
		})

		_, err := global.Redis.Ping(context.TODO()).Result()
		if err == nil {
			global.LogInfo.Info("redis connect!")
		} else {
			global.Log.Error(err.Error())
		}
	}
}
