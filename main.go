package main

import (
	"golang-auth-jwt/handler"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	// to handle middleware, there are many ways
	// refers to https://www.alexedwards.net/blog/making-and-using-middleware
	// refers to https://dasarpemrogramangolang.novalagung.com/B-middleware-using-http-handler.html
	// choose which one that appropriate to our case

	mux.HandleFunc("/signin", handler.Signin)
	mux.Handle("/welcome", handler.MiddlewareAuth(http.HandlerFunc(handler.Welcome)))

	log.Println("Server running on port :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
