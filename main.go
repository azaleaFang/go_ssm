package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go_ssm/common"
	"go_ssm/controllers"
	"log"
	"net/http"
	"os"
)

func main() {
	common.Db, common.Err = sql.Open("sqlite3", "./test.db")
	if common.Err != nil {
		panic(common.Err)
	}
	defer common.Db.Close()
	common.Err = common.Db.Ping()
	if common.Err != nil {
		panic(common.Err)
	}
	os.Setenv("PORT", "8898")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// route
	http.HandleFunc("/", controllers.IndexHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/register", controllers.RegisterHandler)
	http.HandleFunc("/list", controllers.ListHandler)
	http.HandleFunc("/create", controllers.CreateHandler)
	http.HandleFunc("/update", controllers.UpdateHandler)
	http.HandleFunc("/delete", controllers.DeleteHandler)
	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":"+port, nil)

}
