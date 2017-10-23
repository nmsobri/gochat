package util

import (
	"log"
	"fmt"
	"errors"
	"net/http"
	"crypto/sha1"
	"crypto/rand"
	"chitchat/data"
	"html/template"
	"strings"
)

var logger *log.Logger

func GenerateHtml(res http.ResponseWriter, data interface{}, tplFiles ...string) {
	var files []string

	for _, file := range tplFiles {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	tpl := template.Must(template.ParseFiles(files...))
	tpl.ExecuteTemplate(res, "layout", data)
}

func CreateUUID() string {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return uuid
}

func Encrypt(plaintext string) string {
	cryptext := fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext;
}

func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = data.Db.Exec(statement)
	return
}

func Users() (users []data.User, err error) {
	rows, err := data.Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := data.User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

func UserByEmail(email string) (user data.User, err error) {
	user = data.User{}
	err = data.Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByUUID(uuid string) (user data.User, err error) {
	user = data.User{}
	err = data.Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func Threads() (threads []data.Thread, err error) {
	rows, err := data.Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")

	if err != nil {
		return
	}

	for rows.Next() {
		conv := data.Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		threads = append(threads, conv)
	}

	rows.Close()
	return
}

func ThreadByUUID(uuid string) (conv data.Thread, err error) {
	conv = data.Thread{}
	err = data.Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

func Session(writer http.ResponseWriter, request *http.Request) (data.Session, error) {
	var sess data.Session
	cookie, err := request.Cookie("GOLANGSESSION")

	if err == nil {
		sess = data.Session{Uuid: cookie.Value}

		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}

	return sess, err
}

func Info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func Danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func Warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

func ErrorMessage(res http.ResponseWriter, req *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(res, req, strings.Join(url, ""), 302)
}
