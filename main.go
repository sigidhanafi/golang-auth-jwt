package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my-secret-key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credential struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)

	fmt.Println("Listen and serve :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {

	// set header content type to json
	// client expected return is json
	w.Header().Set("Content-Type", "application/json")

	var cred Credential

	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{
			"status":  "error",
			"message": err,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	expectedPassword, ok := users[cred.Username]

	if !ok || expectedPassword != cred.Password {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Authentication failed",
		}

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(response)
		return
	}

	expiredTime := time.Now().Add(5 * time.Minute)

	claim := &Claim{
		Username: cred.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		response := map[string]interface{}{
			"status":  "error",
			"message": err,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "OK",
		"message": "Sign in success",
		"token":   tokenString,
	}

	json.NewEncoder(w).Encode(response)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// set header content type to json
	// client expected return is json
	w.Header().Set("Content-Type", "application/json")

	requestAuth := r.Header.Get("Authorization")

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
			return jwtKey, nil
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

	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "OK",
		"message": "Authorized",
	}

	json.NewEncoder(w).Encode(response)

}
