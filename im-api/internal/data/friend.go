package data

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/biz/do"
	"undersea/im-api/internal/data/po"
)

type FriendRepo struct {
	db       *gorm.DB
	userRepo biz.UserRepo
}

func NewFriendRepo(db *gorm.DB, userRepo biz.UserRepo) biz.FriendRepo {
	return &FriendRepo{
		userRepo: userRepo,
		db:       db,
	}
}

// 获取好友列表
func (r *FriendRepo) GetFriends(ctx context.Context, uid int) (friends []*do.User, err error) {
	// 获取该Uid下普通群聊所有记录
	var friendIds []int
	err = r.db.WithContext(ctx).Model(&po.Friend{}).Select("friend_id").Where("uid=? and deleted=?", uid, false).
		Pluck("friend_id", &friendIds).Error
	if err != nil {
		err = fmt.Errorf("GetFriends->get friend_ids err,%v", err)
		return
	}

	// 根据uids，去用户表查询详细信息
	return r.userRepo.GetUsersByIds(ctx, friendIds)
}

func (r *FriendRepo) AddFriend(ctx context.Context, friendName string, uid int) (err error) {
	// 获取好友id
	user, err := r.userRepo.GetUserByName(ctx, friendName)
	if err != nil {
		err = fmt.Errorf("AddFriend->GetUserByName err,%v", err)
		return
	}

	if user == nil {
		err = fmt.Errorf("用户名不存在")
		return
	}

	if user.Id == uid {
		err = fmt.Errorf("无法邀约自己")
		return
	}

	// 先看下两者是否已经是好友
	var item *po.Friend
	err = r.db.WithContext(ctx).Select("id").Where("uid=? and friend_id=? and deleted=?", uid, user.Id, false).Take(&item).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err = fmt.Errorf("AddFriend->get friend err,%v", err)
		return
	}

	if err == nil {
		// 已经是好友
		return
	}

	// 添加好友信息
	return r.db.WithContext(ctx).Create([]*po.Friend{
		{
			FriendId: user.Id,
			Uid:      uid,
		},
		{
			FriendId: uid,
			Uid:      user.Id,
		},
	}).Error
}
