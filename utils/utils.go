package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"user-authentication/models"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, err models.Error) {

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)

}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func ComparePasswords(hashedPassword string, password []byte) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	

		if err != nil {
			log.Println(err)
			return false
		}

		return true
}

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name": user.Name,
		"id": user.ID,
		"status": user.Status,
		"iss":   "course",
	})

	tokenString, err := token.SignedString([]byte(secret))

	logFatal(err)
	// spew.Dump(token)

	return tokenString, nil
}