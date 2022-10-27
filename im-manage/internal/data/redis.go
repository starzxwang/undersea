package data

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"undersea/im-balance/conf"
)

func NewRedisClient(ctx context.Context, conf conf.Conf) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		Username: conf.Redis.UserName,
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		err = fmt.Errorf("NewRedisClient->ping err,%v", err)
		return
	}

	return
}
