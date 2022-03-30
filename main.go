package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Username string
	Email    string
	Password string
}

func dbConn() (db *sql.DB) {
	//conn string to mysql
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "123456"
	dbName := "binsardb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT username , email , password FROM master_user WHERE stsrc = 'A' ORDER BY username DESC")
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	res := []User{}
	for selDB.Next() {
		var username, email, password string
		err = selDB.Scan(&username, &email, &password)
		if err != nil {
			panic(err.Error())
		}
		usr.Username = username
		usr.Email = email
		usr.Password = password
		res = append(res, usr)
	}
	tpl.ExecuteTemplate(w, "index", res)
	//tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	usrId := r.URL.Query().Get("username")
	selDB, err := db.Query("SELECT username , email , password FROM master_user WHERE stsrc = 'A' AND username=?", usrId)
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	for selDB.Next() {
		var username, email, password string
		err = selDB.Scan(&username, &email, &password)
		if err != nil {
			panic(err.Error())
		}
		usr.Username = username
		usr.Email = email
		usr.Password = password
	}
	tpl.ExecuteTemplate(w, "show", usr)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "new", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	usrId := r.URL.Query().Get("username")
	selDB, err := db.Query("SELECT username, email, password FROM master_user WHERE stsrc = 'A' AND username=?", usrId)
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	for selDB.Next() {
		var username, email, password string
		err = selDB.Scan(&username, &email, &password)
		if err != nil {
			panic(err.Error())
		}
		usr.Username = username
		usr.Email = email
		usr.Password = password
	}
	tpl.ExecuteTemplate(w, "edit", usr)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		insForm, err := db.Prepare("INSERT INTO master_user(username, email, password) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(username, email, password)
		log.Println("INSERT: Username: " + username + " | Email: " + email + " | Password: " + password)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		//id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET email=?, password=? WHERE username=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(username, email, password)
		log.Println("UPDATE: Username: " + username + " | Email: " + email + " | Password: " + password)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	usr := r.URL.Query().Get("username")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE username=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(usr)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

var db *gorm.DB
var err error

func main() {
	log.Println("Server started on: http://localhost:9000")
	cssHandler := http.FileServer(http.Dir("./pub/"))

	http.Handle("/pub/", http.StripPrefix("/pub/", cssHandler))
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":9000", nil)

	// db, err = gorm.Open("mysql", "root:123456@/binsardb")
	// // NOTE: See we’re using = to assign the global var
	// // instead of := which would assign it only in this function

	// if err != nil {
	// 	log.Println("Connection Failed to Open")
	// } else {
	// 	log.Println("Connection Established")
	// }
}
