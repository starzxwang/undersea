package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"undersea/im-balance/internal/biz"
)

const (
	UserIpMappingCacheKey    = "im-balance.user_ip_mapping"
	IpUserListCacheKeyPrefix = "im-balance-ip_user_list-"
)

type BalanceRepo struct {
	rdb *redis.Client
}

func NewBalanceRepo(rdb *redis.Client) biz.BalanceRepo {
	return &BalanceRepo{
		rdb: rdb,
	}
}

func (r *BalanceRepo) GetUserIp(ctx context.Context, uid int) (ip string) {
	return r.rdb.HGet(ctx, UserIpMappingCacheKey, strconv.Itoa(uid)).String()
}

func (r *BalanceRepo) SaveIpUser(ctx context.Context, ip string, uid int) (err error) {
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, UserIpMappingCacheKey, strconv.Itoa(uid), ip).Err()
	pipe.HSet(ctx, r.genIpUserListCacheKey(ip), strconv.Itoa(uid), 1)
	_, err = pipe.Exec(ctx)
	return
}

func (r *BalanceRepo) DeleteIpUser(ctx context.Context, uid int) (err error) {
	return r.rdb.HDel(ctx, UserIpMappingCacheKey, strconv.Itoa(uid)).Err()
}

func (r *BalanceRepo) DeleteIp(ctx context.Context, ip string) (err error) {
	pipe := r.rdb.Pipeline()
	ipUserList, err := pipe.HGetAll(ctx, r.genIpUserListCacheKey(ip)).Result()
	if err != nil {
		return
	}

	uids := make([]string, 0, len(ipUserList))
	for uid := range ipUserList {
		uids = append(uids, uid)
	}

	if len(uids) > 0 {
		pipe.HDel(ctx, UserIpMappingCacheKey, uids...)
	}

	pipe.Del(ctx, r.genIpUserListCacheKey(ip))
	_, err = pipe.Exec(ctx)
	return
}

func (r *BalanceRepo) genIpUserListCacheKey(ip string) string {
	return IpUserListCacheKeyPrefix + ip
}
