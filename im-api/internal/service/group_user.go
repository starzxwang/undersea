package service

import (
	"undersea/im-api/internal/biz"
)

type GroupUserService struct {
	groupUserUseCase *biz.GroupUserUseCase
}

func NewGroupUserService(groupUserUseCase *biz.GroupUserUseCase) *GroupUserService {
	return &GroupUserService{
		groupUserUseCase: groupUserUseCase,
	}
}
