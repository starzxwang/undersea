package data

import (
	"context"
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

func (r *LoginRepo) GetUserIp(ctx context.Context, uid int) (ip string) {
	return r.rdb.HGet(ctx, UserIpMappingCacheKey, strconv.Itoa(uid)).String()
}
