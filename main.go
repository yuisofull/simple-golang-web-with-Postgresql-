package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	DB_NAME     = "godb"
	DB_IP       = "localhost"
)

var db *sql.DB

type MainPageData struct {
	Users []User
}
type User struct {
	ID       int
	Name     string
	Email    string
	Sex      int
	Interest string
}

func addUserPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****addPageHandler running*****")
	tpl, err := template.ParseFiles("./addPage.html")
	if err != nil {
		panic(err)
	}
	if r.Method == "GET" {
		tpl.Execute(w, nil)
		return
	} else if r.Method == "POST" {
		r.ParseForm()
		var data struct {
			Name     string
			Email    string
			Sex      string
			Interest string
		}
		// func (r *Request) FormValue(key string) string
		data.Name = r.FormValue("nameName")
		data.Email = r.FormValue("nameEmail")
		data.Sex = r.FormValue("nameSex")
		data.Interest = r.FormValue("nameInterest")
		if data.Name == "" || data.Email == "" || data.Sex == "" {
			fmt.Println("Error inserting row:", err)
			tpl.Execute(w, "Error inserting data, please check all fields.")
			return
		}
		// if !verification.ValidEmailInputField(r, "nameEmail") {
		// 	fmt.Println("Error inserting row:", err)
		// 	tpl.Execute(w, "Invalid email.")
		// 	return
		// }
		res, err := db.Exec("INSERT INTO users (name, email, sex, interest) VALUES ($1, $2, $3,$4);", data.Name, data.Email, data.Sex, data.Interest)
		if err != nil {
			fmt.Println("Error inserting row:", err)
			tpl.Execute(w, "Error inserting data, please check all fields.")
			return
		}
		lastInserted, _ := res.LastInsertId()
		rowsAffected, _ := res.RowsAffected()
		fmt.Println("ID of last row inserted:", lastInserted)
		fmt.Println("number of rows affected:", rowsAffected)
		tpl.Execute(w, "User Successfully Added")
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****mainPageHandler running*****")
	tpl, err := template.ParseFiles("./mainPage.html")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var data MainPageData
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Sex, &user.Interest); err != nil {
			log.Fatal(err)
		}
		data.Users = append(data.Users, user)
	}
	tpl.Execute(w, data)
}

func updateUserPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updatePageHandler running*****")
	tpl, err := template.ParseFiles("./updatePage.html")
	if err != nil {
		panic(err)
	}
	if r.Method == "GET" {
		r.ParseForm()
		var user User
		id := r.FormValue("userid")
		if err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Sex, &user.Interest); err != nil {
			log.Fatal(err)
		}
		tpl.Execute(w, user)
		return
	}
}

func updateUserResultPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateResultPageHandler running*****")
	if r.Method == "POST" {
		tpl, err := template.ParseFiles("./updatePage.html")
		if err != nil {
			panic(err)
		}
		r.ParseForm()
		sex, err := strconv.Atoi(r.FormValue("nameSex"))
		user := struct {
			User
			ID   string
			Noti string
		}{
			User: User{Name: r.FormValue("nameName"), Email: r.FormValue("nameEmail"), Interest: r.FormValue("nameInterest"), Sex: sex},
			ID:   r.FormValue("userid"),
		}
		user.ID = r.FormValue("userid")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec("UPDATE users SET name = $1, email = $2, sex = $3, interest = $4 WHERE id = $5;", user.Name, user.Email, user.Sex, user.Interest, user.ID); err != nil {
			user.Noti = "Error inserting data, please check all fields."
			tpl.Execute(w, user)
			return
		}
		user.Noti = "Successfully updated"
		tpl.Execute(w, user)
		return
	}
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****deletePageHandler running*****")
	if r.Method == "GET" {
		r.ParseForm()
		ID := r.FormValue("userid")
		tpl, err := template.ParseFiles("./resultPage.html")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec("DELETE FROM users WHERE id = $1;", ID); err != nil {
			tpl.Execute(w, "Error deleting user.")
			return
		}
		tpl.Execute(w, "Successfully deleted user.")
	}
}

func main() {
	var err error
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		DB_USER, DB_PASSWORD, DB_IP, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/add", addUserPageHandler)
	http.HandleFunc("/update/", updateUserPageHandler)
	http.HandleFunc("/updateresult/", updateUserResultPageHandler)
	http.HandleFunc("/delete/", deleteUserPageHandler)
	http.ListenAndServe("localhost:9000", nil)
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			sex INT NOT NULL CHECK (sex IN (1, 2)),
			interest TEXT
		  );
	`); err != nil {
		log.Fatal(err)
	}
}
