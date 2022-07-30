package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"user-authentication/controllers"
	"user-authentication/driver"

	// "github.com/davecgh/go-spew/spew"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)





var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main()  {
	db = driver.ConnectDB()

	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/api/signup", controller.SignUp(db)).Methods("POST")
	router.HandleFunc("/api/login", controller.Login(db)).Methods("POST")

	

	router.HandleFunc("/api/users", controller.TokenVerifyMiddleWare(controller.GetUsers(db))).Methods("GET")
	router.HandleFunc("/api/users/{id}", controller.TokenVerifyMiddleWare(controller.GetUser(db))).Methods("GET")
	router.HandleFunc("/api/users", controller.TokenVerifyMiddleWare(controller.AddUser(db))).Methods("POST")
	router.HandleFunc("/api/users", controller.TokenVerifyMiddleWare(controller.UpdateUser(db))).Methods("PUT")
	router.HandleFunc("/api/users/{id}", controller.TokenVerifyMiddleWare(controller.RemoveUser(db))).Methods("DELETE")

	log.Println(`Listen on port 8000...`)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))
}




