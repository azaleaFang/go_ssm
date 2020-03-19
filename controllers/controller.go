package controllers

import (
	"database/sql"
	"fmt"
	"go_ssm/common"
	"go_ssm/models"
	"html/template"
	"net/http"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "views/register.html")
		return
	}
	// grab user info
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")
	// Check existence of user
	var user models.User
	err := common.Db.QueryRow("SELECT username, password, role FROM users WHERE username=?",
		username).Scan(&user.Username, &user.Password, &user.Role)
	switch {
	// user is available
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		common.CheckInternalServerError(err, w)
		// insert to database
		_, err = common.Db.Exec(`INSERT INTO users(username, password, role) VALUES(?, ?, ?)`,
			username, hashedPassword, role)
		fmt.Println("Created user: ", username)
		common.CheckInternalServerError(err, w)
	case err != nil:
		http.Error(w, "loi: "+err.Error(), http.StatusBadRequest)
		return
	default:
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "views/login.html")
		return
	}
	// grab user info from the submitted form
	username := r.FormValue("usrname")
	password := r.FormValue("psw")
	// query database to get match username
	var user models.User
	err := common.Db.QueryRow("SELECT username, password FROM users WHERE username=?",
		username).Scan(&user.Username, &user.Password)
	common.CheckInternalServerError(err, w)
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login", 301)
	}
	common.Authenticated = true
	http.Redirect(w, r, "/list", 301)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	common.Authenticated = false
	common.IsAuthenticated(w, r)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	common.IsAuthenticated(w, r)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := common.Db.Query("SELECT * FROM cost")
	common.CheckInternalServerError(err, w)
	var funcMap = template.FuncMap{
		"multiplication": func(n float64, f float64) float64 {
			return n * f
		},
		"addOne": func(n int) int {
			return n + 1
		},
	}
	var costs []models.Cost
	var cost models.Cost
	for rows.Next() {
		err = rows.Scan(&cost.Id, &cost.ElectricAmount,
			&cost.ElectricPrice, &cost.WaterAmount, &cost.WaterPrice, &cost.CheckedDate)
		common.CheckInternalServerError(err, w)
		costs = append(costs, cost)
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("views/list.html")
	common.CheckInternalServerError(err, w)
	err = t.Execute(w, costs)
	common.CheckInternalServerError(err, w)

}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	common.IsAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var cost models.Cost
	cost.ElectricAmount, _ = strconv.ParseInt(r.FormValue("ElectricAmount"), 10, 64)
	cost.ElectricPrice, _ = strconv.ParseFloat(r.FormValue("ElectricPrice"), 64)
	cost.WaterAmount, _ = strconv.ParseInt(r.FormValue("WaterAmount"), 10, 64)
	cost.WaterPrice, _ = strconv.ParseFloat(r.FormValue("WaterPrice"), 64)
	cost.CheckedDate = r.FormValue("CheckedDate")
	fmt.Println(cost)

	// Save to database
	stmt, err := common.Db.Prepare(`
		INSERT INTO cost(electric_amount, electric_price, water_amount, water_price, checked_date)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Prepare query error")
		panic(err)
	}
	_, err = stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
		cost.WaterAmount, cost.WaterPrice, cost.CheckedDate)
	if err != nil {
		fmt.Println("Execute query error")
		panic(err)
	}
	http.Redirect(w, r, "/", 301)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	common.IsAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var cost models.Cost
	cost.Id, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	cost.ElectricAmount, _ = strconv.ParseInt(r.FormValue("ElectricAmount"), 10, 64)
	cost.ElectricPrice, _ = strconv.ParseFloat(r.FormValue("ElectricPrice"), 64)
	cost.WaterAmount, _ = strconv.ParseInt(r.FormValue("WaterAmount"), 10, 64)
	cost.WaterPrice, _ = strconv.ParseFloat(r.FormValue("WaterPrice"), 64)
	cost.CheckedDate = r.FormValue("CheckedDate")
	fmt.Println(cost)
	stmt, err := common.Db.Prepare(`
		UPDATE cost SET electric_amount=?, electric_price=?, water_amount=?, water_price=?, checked_date=?
		WHERE id=?
	`)
	common.CheckInternalServerError(err, w)
	res, err := stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
		cost.WaterAmount, cost.WaterPrice, cost.CheckedDate, cost.Id)
	common.CheckInternalServerError(err, w)
	_, err = res.RowsAffected()
	common.CheckInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	common.IsAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var costId, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	stmt, err := common.Db.Prepare("DELETE FROM cost WHERE id=?")
	common.CheckInternalServerError(err, w)
	res, err := stmt.Exec(costId)
	common.CheckInternalServerError(err, w)
	_, err = res.RowsAffected()
	common.CheckInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	common.IsAuthenticated(w, r)
	http.Redirect(w, r, "/list", 301)
}