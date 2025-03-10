package global

import "github.com/redis/go-redis/v9"

var (
	Redis *redis.Client
)

const (
	CemsKeyPrefix       = "cems:"
	LatestDataKeyPrefix = "latestData:"
)
