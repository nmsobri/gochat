package main

import (
	"chitchat/data"
	"chitchat/libs/util"
	"fmt"
	"net/http"
)

func index(res http.ResponseWriter, req *http.Request) {
	threads, err := util.Threads()

	if err != nil {
		util.ErrorMessage(res, req, "Error getting thread list")
	} else {
		if _, err := util.Session(res, req); err != nil {
			util.GenerateHtml(res, threads, "layout", "public.navbar", "index")
		} else {
			util.GenerateHtml(res, threads, "layout", "private.navbar", "index")
		}
	}
}

func err(res http.ResponseWriter, req *http.Request) {
	_get := req.URL.Query()
	_, err := util.Session(res, req)

	if err != nil {
		util.GenerateHtml(res, _get.Get("msg"), "layout", "public.navbar", "error")
	} else {
		util.GenerateHtml(res, _get.Get("msg"), "layout", "private.navbar", "error")
	}
}

func login(res http.ResponseWriter, req *http.Request) {
	util.GenerateHtml(res, nil, "login.layout", "public.navbar", "login")
}

func logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("GOLANGSESSION")

	if err != http.ErrNoCookie {
		util.Warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}

	http.Redirect(res, req, "/", 302)
}

func signup(res http.ResponseWriter, req *http.Request) {
	util.GenerateHtml(res, nil, "login.layout", "public.navbar", "signup")
}

func signupAccount(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		util.Danger(err, "Cannot parse form")
	}

	user := data.User{
		Name:     req.PostFormValue("name"),
		Email:    req.PostFormValue("email"),
		Password: req.PostFormValue("password"),
	}

	if err := user.Create(); err != nil {
		util.Danger(err, "Cannot create user")
	}

	http.Redirect(res, req, "/login", 302)
}

func authenticate(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		util.Danger(err, "Cannot parse form")
	}

	email := req.PostFormValue("email")
	password := req.PostFormValue("password")

	user, err := util.UserByEmail(email)

	if err != nil {
		util.Danger("Cannot find user")
	}

	if user.Password == util.Encrypt(password) {
		sess, err := user.CreateSession()

		if err != nil {
			util.Danger("Cannot create session")
		}

		cookie := http.Cookie{
			Name:     "GOLANGSESSION",
			Value:    sess.Uuid,
			HttpOnly: true,
		}

		http.SetCookie(res, &cookie)
		http.Redirect(res, req, "/", 302)
	} else {
		http.Redirect(res, req, "/login", 302)
	}
}

func newThread(res http.ResponseWriter, req *http.Request) {
	_, err := util.Session(res, req)

	if err != nil {
		http.Redirect(res, req, "/login", 302)
	} else {
		util.GenerateHtml(res, nil, "layout", "private.navbar", "new.thread")
	}
}

func createThread(res http.ResponseWriter, req *http.Request) {
	sess, err := util.Session(res, req)

	if err != nil {
		http.Redirect(res, req, "/login", 302)
	} else {
		if err := req.ParseForm(); err != nil {
			util.Danger(err, "Cannot parse form")
		}

		user, err := sess.User()

		if err != nil {
			util.Danger(err, "Cannot get user from session")
		}

		topic := req.PostFormValue("topic")

		if _, err := user.CreateThread(topic); err != nil {
			util.Danger(err, "Cannot create thread")
		}

		http.Redirect(res, req, "/", 302)
	}
}

func postThread(res http.ResponseWriter, req *http.Request) {
	sess, err := util.Session(res, req)

	if err != nil {
		http.Redirect(res, req, "/login", 302)
	} else {
		if err := req.ParseForm(); err != nil {
			util.Danger(err, "Cannot parse form")
		}

		user, err := sess.User()

		if err != nil {
			util.Danger(err, "Cannot get user from the session")
		}

		body := req.PostFormValue("body")
		uuid := req.PostFormValue("uuid")

		thread, err := util.ThreadByUUID(uuid)

		if err != nil {
			util.ErrorMessage(res, req, "Cannot read thread")
		}

		if _, err := user.CreatePost(thread, body); err != nil {
			util.Danger(err, "Cannot create post")
		}

		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(res, req, url, 302)
	}
}

func readThread(res http.ResponseWriter, req *http.Request) {
	_get := req.URL.Query()
	uuid := _get.Get("id")
	thread, err := util.ThreadByUUID(uuid)

	if err != nil {
		util.ErrorMessage(res, req, "Cannot read thread")
	} else {
		_, err := util.Session(res, req)

		if err != nil {
			util.GenerateHtml(res, thread, "layout", "public.navbar", "public.thread")
		} else {
			util.GenerateHtml(res, thread, "layout", "private.navbar", "private.thread")
		}
	}
}
