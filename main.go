package main

import (
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("/public"))

	serveMux.HandleFunc("/", index)
	serveMux.HandleFunc("/err", err)
	serveMux.HandleFunc("/login", login)
	serveMux.HandleFunc("/logout", logout)
	serveMux.HandleFunc("/signup", signup)
	serveMux.HandleFunc("/signup_account", signupAccount)
	serveMux.HandleFunc("/authenticate", authenticate)
	serveMux.HandleFunc("/thread/new", newThread)
	serveMux.HandleFunc("/thread/create", createThread)
	serveMux.HandleFunc("/thread/post", postThread)
	serveMux.HandleFunc("/thread/read", readThread)
	serveMux.Handle("/static", http.StripPrefix("/static", fileServer))

	fmt.Println("Server running at 127.0.0.1:8000")
	http.ListenAndServe(":8000", serveMux)
}
