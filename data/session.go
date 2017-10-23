package data

import "time"

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func (self *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", self.Uuid).
		Scan(&self.Id, &self.Uuid, &self.Email, &self.UserId, &self.CreatedAt)

	if err != nil {
		valid = false
		return
	}

	if self.Id != 0 {
		valid = true
	}

	return
}

func (self *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(self.Uuid)
	return
}

func (self *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", self.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)

	return
}
