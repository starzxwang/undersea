package po

import "time"

type SmallGroupUser struct {
	Id        int
	Gid       string
	uid       int
	Deleted   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*SmallGroupUser) TableName() string {
	return "im_small_group_user"
}
