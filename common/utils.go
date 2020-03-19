package common

import (
	"database/sql"
	"net/http"
)

var (
	Db  *sql.DB
	Err error
	Authenticated bool
)

func CheckInternalServerError(err error,w http.ResponseWriter)  {
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
}

func IsAuthenticated(w http.ResponseWriter,r *http.Request)  {
	if !Authenticated {
		http.Redirect(w, r, "/login", 301)
	}
}