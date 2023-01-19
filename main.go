package main

import (
	"encoding/json"
	"golang-auth-jwt/handler"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Credential struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// exclude /signin from middleware auth
		if r.URL.Path == "/signin" {
			next.ServeHTTP(w, r)
		}

		env, err := godotenv.Read(".env")
		if err != nil {
			log.Fatal("Error load env file")
			return
		}

		// get http header auth
		requestAuth := r.Header.Get("Authorization")

		// if no auth on header, then No Authorized
		if len(requestAuth) <= 0 {
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]interface{}{
				"status":  "error",
				"message": "No Authorized",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		// get token by slicing string
		requestAuthSlice := strings.Split(requestAuth, "Bearer ")
		requestToken := requestAuthSlice[1]

		claim := &Claim{}

		token, err := jwt.ParseWithClaims(
			requestToken,
			claim,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(env["JWT_SECRET_KEY"]), nil
			},
		)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]interface{}{
				"status":  "error",
				"message": "No Authorized",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	mux := http.DefaultServeMux
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/welcome", handler.Welcome)

	var handler http.Handler = mux
	// add middleware
	handler = MiddlewareAuth(handler)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = handler

	log.Println("Server running on port :8080")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
