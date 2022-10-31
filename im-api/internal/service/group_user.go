package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/service/dto"
	"undersea/pkg/api"
)

type GroupUserService struct {
	groupUserUseCase *biz.GroupUserUseCase
}

func NewGroupUserService(groupUserUseCase *biz.GroupUserUseCase) *GroupUserService {
	return &GroupUserService{
		groupUserUseCase: groupUserUseCase,
	}
}

// 好友列表
func (s *GroupUserService) GetFriends(c *gin.Context) {
	var req dto.GetFriendsReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	if req.Uid == 0 {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, "参数不能为空", nil))
		return
	}

	friends, err := s.groupUserUseCase.GetFriends(c, req.Uid)
	if err != nil {
		c.JSON(http.StatusOK, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.Success(dto.ConvertUsersDO2DTO(friends)))

	return
}
