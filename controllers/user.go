package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"user-authentication/models"
	"user-authentication/utils"

	userRepository "user-authentication/repository/user"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct{}


func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) SignUp(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
	
		json.NewDecoder(r.Body).Decode(&user)
		// spew.Dump(user)
	
		if user.Email == "" {
			error.Message = "Email is missing."
	
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
	
		if user.Password == "" {
			error.Message = "Password is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
	
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		logFatal(err)
	
		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		userID, err := userRepo.SignUp(db, user)

		user.ID = userID
	
		if err != nil {
			error.Message = "Server error."
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
	
		user.Password = ""
	
		w.Header().Set("Content-Type", "application/json")
	
		utils.ResponseJSON(w, user)
	
	}
}



func (c Controller) Login(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error
		var errorObject models.Error
		
	
		json.NewDecoder(r.Body).Decode(&user)
	
		if user.Email == "" {
			error.Message = "Email is missing."
	
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
	
		if user.Password == "" {
			error.Message = "Password is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
	
		password := user.Password
	
		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(db, user)
	
		hashedPassword := user.Password
		
		logFatal(err)
		
	
		token, err := utils.GenerateToken(user)
	
		logFatal(err)

		isValidPassword := utils.ComparePasswords(hashedPassword, []byte(password))

		

		if isValidPassword {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Authorization", token)

			jwt.Token = token
		utils.ResponseJSON(w, jwt)
		} else {
			errorObject.Message = "Invalid Password"

			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
		}
	
	
	}
} 


func (c Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error :=jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok :=token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)

				return

			}

			// spew.Dump(token)
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token"
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
		}

	})
}

