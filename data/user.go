package data

import (
	"time"
	"chitchat/libs/util"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (self *User) CreateThread(topic string) (conv Thread, err error) {
	statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	stmt.Exec()

	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(util.CreateUUID(), topic, self.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

func (self *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(util.CreateUUID(), body, self.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

func (self *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(util.CreateUUID(), self.Email, self.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (self *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", self.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (self *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(util.CreateUUID(), self.Name, self.Email, util.Encrypt(self.Password), time.Now()).Scan(&self.Id, &self.Uuid, &self.CreatedAt)
	return
}

func (self *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(self.Id)
	return
}

func (self *User) Update() (err error) {
	statement := "update users set name = $2, email = $3 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(self.Id, self.Name, self.Email)
	return
}
