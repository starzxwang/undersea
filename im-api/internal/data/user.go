package data

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/biz/do"
	"undersea/im-api/internal/data/po"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) biz.UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetUserByName(ctx context.Context, username string) (user *do.User, err error) {
	var userPO *po.User
	err = r.db.WithContext(ctx).Where("`name`=? and deleted=?", username, false).Take(&userPO).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		err = fmt.Errorf("GetUserByName->take err,%v", err)
		return
	}

	return po.ConvertUserPO2DO(userPO), nil
}

func (r *UserRepo) Register(ctx context.Context, user *do.User) (id int, err error) {
	userPO := po.ConvertUserDO2PO(user)
	err = r.db.WithContext(ctx).Create(&userPO).Error

	if err != nil {
		err = fmt.Errorf("Register->create err,%v", err)
		return
	}

	return userPO.Id, nil
}
