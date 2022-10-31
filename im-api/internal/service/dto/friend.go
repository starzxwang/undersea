package dto

type GetFriendsReq struct {
	Uid int `uri:"uid"`
}

type AddFriendReq struct {
	FriendName string `form:"friend_name"`
	Uid        int    `form:"uid"`
}

type GetFriendsResp struct {
	Friends []*User `json:"friends"`
}
