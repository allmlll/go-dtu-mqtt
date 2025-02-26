package util

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func RedisReceiveMessageTimeOut(c *redis.PubSub, ctx context.Context, timeout time.Duration) (*redis.Message, error) {
	for {
		msg, err := c.ReceiveTimeout(ctx, timeout)
		if err != nil {
			return nil, err
		}

		switch msg := msg.(type) {
		case *redis.Subscription:
			// Ignore.
		case *redis.Pong:
			// Ignore.
		case *redis.Message:
			return msg, nil
		default:
			err := fmt.Errorf("redis: unknown message: %T", msg)
			return nil, err
		}
	}
}
