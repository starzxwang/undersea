package po

import "time"

type Friend struct {
	Id        int
	FriendId  int
	Uid       int
	Deleted   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Friend) TableName() string {
	return "im_friend"
}
