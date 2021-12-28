package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)
	type User struct {
		Mail string // <-- CHANGED THIS LINE
		Password string
		Id   string // <-- CHANGED THIS LINE
		Role string
	}
func main() {

	r := mux.NewRouter()
	godotenv.Load()
	port := os.Getenv("PORT")
	r.HandleFunc("/login/email/{mail}", Fetch).Methods("GET")
  r.HandleFunc("/login/{email}/{password}",loginCheck).Methods("GET")
	// if there is an error opening the connection, handle it

	// defer the close till after the main function has finished
	// executing

	//insert, err := db.Query("INSERT INTO PERSON VALUES (1,'vamsi','krishna','vamsijavvadi','Ilovemydad7@','student',9,'vamsijavvadi@gmail.com')")
	// insert, err := db.Query("INSERT INTO USER VALUES ('vamsijavvadi','vamsijavvadi@gmail.com')")
	//  if err != nil {
	//         panic(err.Error())
	//   }
	// if there is an error inserting, handle it

	// be careful deferring Queries if you are using transactions

	log.Fatal(http.ListenAndServe(":"+port, r))

	//   defer insert.Close()
}

func loginCheck(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
  
db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("call getstudents()")
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Mail, &user.Password, &user.Id, &user.Role)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	defer rows.Close()
	type Result struct{
		Status string `json:"status"`
		Role string `json:"role"`
	}
	res:=Result{Status:"False",Role:""}
	for _, item := range users {
		if item.Mail== params["email"]  && item.Password==params["password"]{
			res.Role=item.Role
			res.Status="True"
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	json.NewEncoder(w).Encode(res)

}
func Fetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("call getstudents()")
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Mail, &user.Password, &user.Id, &user.Role)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	defer rows.Close()
	type Result struct {
		Status string `json:"status"`
		Role   string `json:"role"`
	}
	res := Result{Status: "False", Role: ""}

	for _, elem := range users {
		if elem.Mail == params["mail"] {
			res.Status = "True"
			res.Role = elem.Role
			break
		}
	}

	json.NewEncoder(w).Encode(res)
	// return(ToJSON(users)) // <-- CHANGED THIS LINE
}

func ToJSON(obj interface{}) string {

	res, err := json.Marshal(obj)

	if err != nil {
		panic("error with json serialization " + err.Error())
	}
	return string(res)
}
