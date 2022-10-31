package biz

import (
	"context"
	"undersea/im-api/internal/biz/do"
)

type FriendUseCase struct {
	friendRepo FriendRepo
}

func NewFriendUseCase(friendRepo FriendRepo) *FriendUseCase {
	return &FriendUseCase{
		friendRepo: friendRepo,
	}
}

func (uc *FriendUseCase) GetFriends(ctx context.Context, uid int) (friends []*do.User, err error) {
	return uc.friendRepo.GetFriends(ctx, uid)
}

func (uc *FriendUseCase) AddFriend(ctx context.Context, friendName string, uid int) (err error) {
	return uc.friendRepo.AddFriend(ctx, friendName, uid)
}

type FriendRepo interface {
	GetFriends(ctx context.Context, uid int) (friends []*do.User, err error)
	AddFriend(ctx context.Context, friendName string, uid int) (err error)
}
