package data

import "time"

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

func (self *Thread) CreatedAtDate() string {
	return self.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (self *Thread) NumReplies() int {
	var count int
	rows, err := Db.Query("SELECT COUNT(*) FROM posts where thread_id = $1", self.Id)

	if err != nil {
		return count
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return count
		}
	}

	return count
}

func (self *Thread) Posts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", self.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func (self *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", self.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}
