package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"undersea/im-api/internal/biz"
	"undersea/im-api/internal/service/dto"
	"undersea/pkg/api"
)

type UserService struct {
	userUseCase *biz.UserUseCase
}

func NewUserService(userUseCase *biz.UserUseCase) *UserService {
	return &UserService{
		userUseCase: userUseCase,
	}
}

func (s *UserService) Login(c *gin.Context) {
	var req dto.LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	user, err := s.userUseCase.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.Success(dto.ConvertUserDO2DTO(user)))
	return
}

func (s *UserService) Register(c *gin.Context) {
	var req dto.RegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	id, err := s.userUseCase.Register(c, req.Username, req.Password, req.Avatar)
	if err != nil {
		c.JSON(http.StatusOK, api.Failed(api.CodeParamError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.Success(&dto.RegisterResp{
		Id: id,
	}))

	return
}
