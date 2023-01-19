package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

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

func Signin(w http.ResponseWriter, r *http.Request) {

	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error load env file")
	}
	// set header content type to json
	// client expected return is json
	w.Header().Set("Content-Type", "application/json")

	var cred Credential

	err = json.NewDecoder(r.Body).Decode(&cred)

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

	tokenString, err := token.SignedString([]byte(env["JWT_SECRET_KEY"]))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		response := map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
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

	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "OK",
		"message": "Authorized",
	}

	json.NewEncoder(w).Encode(response)

}
