package do

import "time"

type User struct {
	Id        int
	Pwd       string
	Avatar    string
	Name      string
	Deleted   int
	CreatedAt time.Time
}
