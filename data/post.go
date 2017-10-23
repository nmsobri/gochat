package data

import (
	"time"
)

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

func (self *Post) CreatedAtDate() string {
	return self.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (self *Post) User() User {
	user := User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", self.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return user
}
