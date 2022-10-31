package dto

import (
	"time"
	"undersea/im-api/internal/biz/do"
)

type LoginReq struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterReq struct {
	Avatar   string `form:"avatar"`
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterResp struct {
	Id int `json:"id"`
}

type User struct {
	Id        int       `json:"id"`
	Avatar    string    `json:"avatar"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ConvertUserDO2DTO(userDO *do.User) *User {
	return &User{
		Id:        userDO.Id,
		Avatar:    userDO.Avatar,
		Name:      userDO.Name,
		CreatedAt: userDO.CreatedAt,
	}
}

func ConvertUsersDO2DTO(usersDO []*do.User) []*User {
	ret := make([]*User, 0, len(usersDO))
	for _, userDO := range usersDO {
		ret = append(ret, ConvertUserDO2DTO(userDO))
	}

	return ret
}
