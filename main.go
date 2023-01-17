package main

import (
	"fmt"
	"golang-auth-jwt/handler"
	"log"
	"net/http"
)

func main() {
	authentication := handler.NewAuth()

	http.HandleFunc("/signin", authentication.Signin)
	http.HandleFunc("/welcome", authentication.Welcome)

	fmt.Println("Listen and serve :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
