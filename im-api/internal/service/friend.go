package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/service/dto"
	"undersea/pkg/api"
	"undersea/pkg/log"
)

type FriendService struct {
	friendUseCase *biz.FriendUseCase
}

func NewFriendService(friendUseCase *biz.FriendUseCase) *FriendService {
	return &FriendService{
		friendUseCase: friendUseCase,
	}
}

// 添加好友
func (s *FriendService) AddFriend(c *gin.Context) {
	var req dto.AddFriendReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	if req.FriendName == "" {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, "参数不能为空", nil))
		return
	}

	err = s.friendUseCase.AddFriend(c, req.FriendName, req.Uid)
	if err != nil {
		log.E(c, err).Msgf("AddFriend->AddFriend err")
		c.JSON(http.StatusOK, api.Failed(api.CodeException, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.Success(nil))

	return
}

// 好友列表
func (s *FriendService) GetFriends(c *gin.Context) {
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

	friends, err := s.friendUseCase.GetFriends(c, req.Uid)
	if err != nil {
		log.E(c, err).Msgf("GetFriends->GetFriends err")
		c.JSON(http.StatusOK, api.Failed(api.CodeException, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.Success(&dto.GetFriendsResp{
		Friends: dto.ConvertUsersDO2DTO(friends),
	}))

	return
}
