package biz

import (
	"context"
	"undersea/im-api/internal/biz/do"
)

type GroupUserUseCase struct {
	groupUserRepo GroupUserRepo
}

func NewGroupUseCase(groupUserRepo GroupUserRepo) *GroupUserUseCase {
	return &GroupUserUseCase{
		groupUserRepo: groupUserRepo,
	}
}

func (uc *GroupUserUseCase) GetFriends(ctx context.Context, uid int) (friends []*do.User, err error) {
	return uc.groupUserRepo.GetFriends(ctx, uid)
}

type GroupUserRepo interface {
	GetFriends(ctx context.Context, uid int) (friends []*do.User, err error)
}
