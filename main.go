package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
"sort"
"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	
)

type User struct {
	Mail     string // <-- CHANGED THIS LINE
	Password string
	Id       string // <-- CHANGED THIS LINE
	Role     string
}
type StudentDetails struct {
	Name          string  `json:"name"` // <-- CHANGED THIS LINE
	Adno          string  `json:"regno"`
	Self          float64 `json:"self"` // <-- CHANGED THIS LINE
	Faculty       float64 `json:"faculty"`
	Competencynum int     `json:"competencynum"`
}

type Competen struct {
	StudentDetails `json:"competency"`
}
type CompetencyReturn struct {
	C Competen
}

var speciality_for_faculty string = ""



func main() {

	r := mux.NewRouter()
	godotenv.Load()
	port := os.Getenv("PORT")
	r.HandleFunc("/login/email/{mail}", Fetch).Methods("GET")
	r.HandleFunc("/login/{email}/{password}", loginCheck).Methods("GET")
	r.HandleFunc("/fdashboard/details/{email}", getfacultydetails).Methods("GET")

	r.HandleFunc("/fdashboard/competencydetails/{speciality}", getcompnames).Methods("GET")
	r.HandleFunc("/fdashboard/competencydetails/speciality/{speciality}", getcompetencyalongwithstudents).Methods("GET")
	r.HandleFunc("/profile/email/{email}", getprofile).Methods("GET")
	r.HandleFunc("/competencyevaluations/competencyid/{competencyid}/studentid/{studentid}", getcompetencyevaluations).Methods("GET")
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

func getcompetencyevaluations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	
	//var competencyids []int=[]int{}

	type Evaluation struct {
		CompEvaId int `json:"compentencyevaluationid"` 
		Opnum string `json:"patientop"` // <-- CHANGED THIS LINE
		Date string `json:"date"`
		Time string `json:"time"`
		Self float64 `json:"self"`
		Faculty float64 `json:"faculty"`
		Timest string `json:"-"`

	}

	evalrow, er := db.Query("call getallevalofacompetency(?,?)",params["competencyid"],params["studentid"])
	if er != nil {

		panic(err.Error())

	}
	defer evalrow.Close()
	et := []Evaluation{}
	
	for evalrow.Next() {
		user := new(Evaluation)
		err := evalrow.Scan(&user.CompEvaId,&user.Opnum, &user.Date,&user.Time)

		if err != nil {
			panic(err)
		
		}
		datab, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	

StudentF, er := datab.Query("CALL getpercentageforeacheval(?,?,?)", "faculty",params["competencyid"],user.CompEvaId)
	if er != nil {

		panic(err.Error())

	}

	
	
	for StudentF.Next() {
		
		err := StudentF.Scan(&user.Faculty)

		if err != nil {
			panic(err)
		}
		
	}

	StudentF.Close()

	StudentS, er := datab.Query("CALL getpercentageforeacheval(?,?,?)", "self",params["competencyid"],user.CompEvaId)

	if er != nil {

		panic(err.Error())

	}

	for StudentS.Next() {
		
		err := StudentS.Scan(&user.Self)

		if err != nil {
			panic(err)
		}
		
		}
	
	StudentS.Close()
	 datab.Close()
	 user.Timest=user.Date+" "+user.Time
et = append(et,*user)

	}




	
	sort.Slice(et, func(i, j int) bool {
    return et[i].Timest < et[j].Timest
})


	json.NewEncoder(w).Encode(et)

}
func getprofile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("call getprofile(?)", params["email"])
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()

	type Persondetails struct {
		Name  string `json:"name"`
		Role  string `json:"role"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Batch string `json:"batch"`
	}

	pd := make([]*Persondetails, 0)
	for rows.Next() {
		person := new(Persondetails)
		err := rows.Scan(&person.Name, &person.Phone, &person.Email, &person.Role)

		if err != nil {
			panic(err)
		}

		if person.Role == "student" {

			row, err := db.Query("call getbatch(?)", params["email"])
			if err != nil {

				panic(err.Error())

			}

			for row.Next() {

				err := row.Scan(&person.Batch)
				if err != nil {
					panic(err)
				}

			}
			row.Close()

		}
		pd = append(pd, person)
	}
	defer rows.Close()

	json.NewEncoder(w).Encode(pd)

}

func getcompetencyalongwithstudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	speciality_for_faculty = params["speciality"]
	//var competencyids []int=[]int{}

	type Students struct {
		Name string `json:"name"` // <-- CHANGED THIS LINE
		Adno string `json:"regno"`
	}

	studentrow, er := db.Query("call getstudents()")
	if er != nil {

		panic(err.Error())

	}
	defer studentrow.Close()
	st := make([]*Students, 0)
	for studentrow.Next() {
		user := new(Students)
		err := studentrow.Scan(&user.Name, &user.Adno)

		if err != nil {
			panic(err)
		}
		st = append(st, user)

	}

	studentD := make([]*StudentDetails, 0)

	// compnamelist := []string{}
	// for row.Next() {
	// 	var str string
	// 	err := row.Scan(&str)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	compnamelist = append(compnamelist, str)
	// }

	StudentF, er := db.Query("CALL getevalpercentage(?,?)", params["speciality"], "faculty")
	if er != nil {

		panic(err.Error())

	}

	type Score struct {
		Adno          string
		Competency_id int
		Self          float64 `json:"self"` // <-- CHANGED THIS LINE
		Faculty       float64 `json:"faculty"`
	}
	scores := make([]*Score, 0)

	for StudentF.Next() {
		sc := new(Score)
		err := StudentF.Scan(&sc.Faculty, &sc.Adno, &sc.Competency_id)

		if err != nil {
			panic(err)
		}
		scores = append(scores, sc)
	}

	StudentF.Close()

	StudentS, er := db.Query("CALL getevalpercentage(?,?)", params["speciality"], "self")

	if er != nil {

		panic(err.Error())

	}

	for StudentS.Next() {
		var cid int
		var sid string
		var selfpercentage float64

		err := StudentS.Scan(&selfpercentage, &sid, &cid)

		if err != nil {
			panic(err)
		}
		for index, item := range scores {
			if item.Adno == sid {
				scores = append(scores[:index], scores[index+1:]...)
				scores = append(scores, &Score{Adno: item.Adno, Competency_id: item.Competency_id, Self: selfpercentage, Faculty: item.Faculty})

				break
			}

		}
	}
	StudentS.Close()
	for _, item := range scores {
		for _, sitem := range st {
			if item.Adno == sitem.Adno {
				studentD = append(studentD, &StudentDetails{Name: sitem.Name, Adno: item.Adno, Self: item.Self, Faculty: item.Faculty, Competencynum: item.Competency_id})
			} else {
				studentD = append(studentD, &StudentDetails{Name: sitem.Name, Adno: item.Adno, Self: 0, Faculty: 0, Competencynum: item.Competency_id})
			}

		}
	}

	Compre := make([]*CompetencyReturn, 0)
	for _, stud := range studentD {
		stude := StudentDetails{Name: stud.Name, Adno: stud.Adno, Self: stud.Self, Faculty: stud.Faculty, Competencynum: stud.Competencynum}

		Compre = append(Compre, &CompetencyReturn{C: Competen{StudentDetails: stude}})
	}

	json.NewEncoder(w).Encode(Compre)

}
func (c CompetencyReturn) MarshalJSON() ([]byte, error) {
	// encode the original
	m, _ := json.Marshal(c.C)

	// decode it back to get a map
	var a interface{}
	json.Unmarshal(m, &a)
	b := a.(map[string]interface{})

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select Competency_Name,competency_id from competency where Speciality_id in ( select Speciality_id from speciality where Speciality_Name=?);", speciality_for_faculty)
	if err != nil {

		panic(err.Error())

	}

	defer db.Close()

	//var competencyids []int=[]int{}

	type Competency struct {
		Name string  `json:"name"` // <-- CHANGED THIS LINE
		Cid  float64 `json:"cid"`
	}
	comp := make([]*Competency, 0)
	for rows.Next() {
		onec := new(Competency)
		err := rows.Scan(&onec.Name, &onec.Cid)

		if err != nil {
			panic(err)
		}
		comp = append(comp, onec)

	}
	defer rows.Close()

	for i, si := range b {
		var f interface{}
		n, _ := json.Marshal(b[i])
		json.Unmarshal(n, &f)
		c := f.(map[string]interface{})
		//idx := string(c["id"])

		idx := c["competencynum"].(float64)
		for _, co := range comp {

			if co.Cid == idx {
				b[co.Name] = si

				delete(b, "competency")
			}
		}

	}

	return json.Marshal(b)

}

func getcompnames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("call getcompetencies(?)", params["speciality"])
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()

	type Result struct {
		Competency []string `json:"competency"`
	}

	res := Result{Competency: []string{}}
	for rows.Next() {
		var str string
		err := rows.Scan(&str)

		if err != nil {
			panic(err)
		}
		res.Competency = append(res.Competency, str)
		res = Result{Competency: res.Competency}

	}
	defer rows.Close()

	json.NewEncoder(w).Encode(res)

}
func getfacultydetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select concat(p.first_name,p.last_name) as name,f.speciality from person p,faculty f where p.person_id=f.person_id and p.email=?;", params["email"])
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()
	type Result struct {
		Name       string `json:"name"`
		Speciality string `json:"speciality"`
	}

	res := make([]*Result, 0)
	for rows.Next() {
		user := new(Result)
		err := rows.Scan(&user.Name, &user.Speciality)

		if err != nil {
			panic(err)
		}
		res = append(res, user)
	}
	defer rows.Close()

	json.NewEncoder(w).Encode(res)

}
func loginCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("call getpersons()")
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
	for _, item := range users {
		if item.Mail == params["email"] && item.Password == params["password"] {
			res.Role = item.Role
			res.Status = "True"
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

	rows, err := db.Query("call getpersons()")
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
