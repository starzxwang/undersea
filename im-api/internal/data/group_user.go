package data

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/biz/do"
)

type GroupUserRepo struct {
	db       *gorm.DB
	userRepo *UserRepo
}

func NewGroupUserRepo(db *gorm.DB, userRepo *UserRepo) biz.GroupUserRepo {
	return &GroupUserRepo{
		db: db,
	}
}

// 获取好友列表
func (r *GroupUserRepo) GetFriends(ctx context.Context, uid int) (friends []*do.User, err error) {
	// 获取该Uid下普通群聊所有记录
	var gids []string
	err = r.db.WithContext(ctx).Select("gid").Where("uid=? and deleted=?", uid, false).Pluck("gid", &gids).Error

	if err != nil {
		err = fmt.Errorf("GetFriends->get gids err,%v", err)
		return
	}

	if len(gids) == 0 {
		return []*do.User{}, nil
	}

	// 通过这些记录的id，获取其好友的uids
	var uids []int
	err = r.db.WithContext(ctx).Select("uid").Where("gid in(?) and deleted=? and uid!=?", gids, false, uid).
		Pluck("uid", &uids).Error

	if err != nil {
		err = fmt.Errorf("GetFriends->get uids err,%v", err)
		return
	}

	if len(uids) == 0 {
		return []*do.User{}, nil
	}

	// 根据uids，去用户表查询详细信息
	return r.userRepo.GetUsersByIds(ctx, uids)
}
