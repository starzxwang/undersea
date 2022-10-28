package biz

import (
	"context"
	"fmt"
	"time"
	"undersea/im-api/internal/biz/do"
	"undersea/pkg/encode"
)

type UserUseCase struct {
	UserRepo UserRepo
}

func NewUserUseCase(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (uc *UserUseCase) Login(ctx context.Context, username, pwd string) (user *do.User, err error) {
	user, err = uc.UserRepo.GetUserByName(ctx, username)
	if err != nil {
		err = fmt.Errorf("Login->GetUserByName err,%v", err)
		return
	}

	if user == nil {
		err = fmt.Errorf("该用户名不存在")
		return
	}

	if user.Pwd != encode.EncodeMd5(pwd) {
		err = fmt.Errorf("密码不正确")
		return
	}

	return
}

func (uc *UserUseCase) Register(ctx context.Context, username, pwd, avatar string) (id int, err error) {
	// 用户名是否已经存在
	user, err := uc.UserRepo.GetUserByName(ctx, username)
	if err != nil {
		err = fmt.Errorf("Register->GetUserByName err,%v", err)
		return
	}

	if user != nil {
		err = fmt.Errorf("该用户名已经存在")
		return
	}

	return uc.UserRepo.Register(ctx, &do.User{
		Pwd:       pwd,
		Avatar:    avatar,
		Name:      username,
		CreatedAt: time.Now(),
	})
}

type UserRepo interface {
	GetUserByName(ctx context.Context, username string) (user *do.User, err error)
	Register(ctx context.Context, user *do.User) (id int, err error)
}
