package po

import (
	"time"
	"undersea/im-api/internal/biz/do"
	"undersea/pkg/encode"
)

type User struct {
	Id        int
	Pwd       string
	Avatar    string
	Name      string
	Deleted   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*User) TableName() string {
	return "im_user"
}

func ConvertUserPO2DO(user *User) *do.User {
	return &do.User{
		Id:        user.Id,
		Avatar:    user.Avatar,
		Pwd:       user.Pwd,
		Name:      user.Name,
		Deleted:   user.Deleted,
		CreatedAt: user.CreatedAt,
	}
}

func ConvertUserDO2PO(user *do.User) *User {
	return &User{
		Id:        user.Id,
		Avatar:    user.Avatar,
		Pwd:       encode.EncodeMd5(user.Pwd),
		Name:      user.Name,
		Deleted:   user.Deleted,
		CreatedAt: user.CreatedAt,
	}
}

func ConvertUsersPO2DO(users []*User) (ret []*do.User) {
	if len(users) == 0 {
		return []*do.User{}
	}

	ret = make([]*do.User, 0, len(users))
	for _, user := range users {
		ret = append(ret, ConvertUserPO2DO(user))
	}

	return
}
