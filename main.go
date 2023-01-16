package main

import (
	"fmt"
	"golang-auth-jwt/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signin", handler.Signin)
	http.HandleFunc("/welcome", handler.Welcome)

	fmt.Println("Listen and serve :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
