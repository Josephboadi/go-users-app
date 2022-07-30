package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"user-authentication/models"
	userRepository "user-authentication/repository/user"
	"user-authentication/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)



var users []models.User



func (c Controller) GetUsers(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error


		users = []models.User{}
		userRepo := userRepository.UserRepository{}

		users, err := userRepo.GetUsers(db, user, users)
		
		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
		}
		
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, users)
	}
}

func (c Controller) GetUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		params := mux.Vars(r)

		users = []models.User{}
		userRepo := userRepository.UserRepository{}
	
		id, _ := strconv.Atoi(params["id"])
		// log.Println(id)
		// typeof var
		// log.Println(reflect.TypeOf(i))

		

		user, err := userRepo.GetUser(db, user, id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not Found"
				utils.RespondWithError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server error"
				utils.RespondWithError(w, http.StatusInternalServerError, error)
				return
			}
		}
	
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)
	
		
	}
}



func (c Controller) AddUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
	
		json.NewDecoder(r.Body).Decode(&user)
	
		if user.Email == "" || user.Name == ""  {
			error.Message = "Enter missing fields."
			utils.RespondWithError(w, http.StatusBadRequest, error)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte("12345"), 10)
		logFatal(err)
	
		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		userID, err := userRepo.AddUser(db, user)

		user.ID = userID


		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "text/plain")
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) UpdateUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		
		json.NewDecoder(r.Body).Decode(&user)
	
		if user.ID == 0 || user.Name == "" || user.Email == ""  {
			error.Message = "All fields are required."
			utils.RespondWithError(w, http.StatusBadRequest, error)
		}
	
		userRepo := userRepository.UserRepository{}
		rowsUpdated, err := userRepo.UpdateUser(db, user)

		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.ResponseJSON(w, rowsUpdated)
	
	}
}


func (c Controller) RemoveUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error

		params := mux.Vars(r)
		userRepo := userRepository.UserRepository{}
	
		id, _ := strconv.Atoi(params["id"])
		
		rowsDeleted, err := userRepo.RemoveUser(db, id)

		if err != nil {
			error.Message = "Server error."
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		
		if rowsDeleted == 0 {
			error.Message = "Not Found"
			utils.RespondWithError(w, http.StatusNotFound, error)
			return
		}
	
		w.Header().Set("Content-Type", "text/plain")
		utils.ResponseJSON(w, rowsDeleted)
	}
}
