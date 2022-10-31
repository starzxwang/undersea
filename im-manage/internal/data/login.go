package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"undersea/im-manage/internal/biz"
)

const (
	UserIpMappingCacheKey = "im-balance.user_ip_mapping"
)

type LoginRepo struct {
	rdb *redis.Client
}

func NewLoginRepo(rdb *redis.Client) biz.LoginRepo {
	return &LoginRepo{
		rdb: rdb,
	}
}

func (r *LoginRepo) GetUserIp(ctx context.Context, uid int) (ip string, err error) {
	ip, err = r.rdb.HGet(ctx, UserIpMappingCacheKey, strconv.Itoa(uid)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}

	if err != nil {
		err = fmt.Errorf("GetUserIp->hget err,%v", err)
		return
	}

	return
}
