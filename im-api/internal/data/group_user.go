package data

import (
	"gorm.io/gorm"
	"undersea/im-api/internal/biz"
)

type GroupUserRepo struct {
	db       *gorm.DB
	userRepo biz.UserRepo
}

func NewGroupUserRepo(db *gorm.DB, userRepo biz.UserRepo) biz.GroupUserRepo {
	return &GroupUserRepo{
		userRepo: userRepo,
		db:       db,
	}
}
