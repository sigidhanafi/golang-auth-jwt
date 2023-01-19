package main

import (
	"golang-auth-jwt/handler"
	"log"
	"net/http"
)

func main() {

	mux := http.DefaultServeMux
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/welcome", handler.Welcome)

	var httpHandler http.Handler = mux
	// add middleware
	httpHandler = handler.MiddlewareAuth(httpHandler)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = httpHandler

	log.Println("Server running on port :8080")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
