package main

import (
	"database/sql"
	"encoding/json"

	//"fmt"

	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

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



func main() {

	r := mux.NewRouter()
	godotenv.Load()
	port := os.Getenv("PORT")
	r.HandleFunc("/login/email/{mail}", Fetch).Methods("GET")
	r.HandleFunc("/login/{email}/{password}", loginCheck).Methods("GET")
	r.HandleFunc("/fdashboard/details/{email}", getfacultydetails).Methods("GET")

	r.HandleFunc("/fdashboard/competencydetails/{speciality}", getcompnames).Methods("GET")
	r.HandleFunc("/fdashboard/competencydetails/speciality/{speciality}/competencyid/{competencyid}", getcompetencyalongwithstudents).Methods("GET")
	r.HandleFunc("/profile/email/{email}", getprofile).Methods("GET")
	r.HandleFunc("/competencyevaluations/competencyid/{competencyid}/studentid/{studentid}", getcompetencyevaluations).Methods("GET")
	r.HandleFunc("/competencyevaluations/competencyid/{competencyid}/studentid/{studentid}/opnum", addroweval).Methods("POST")
	r.HandleFunc("/competencyevaluations/competencyid/{competencyid}/competencyevaluationid/{competencyevaluationid}", getfeedbackform).Methods("GET")
	// r.HandleFunc("/competencyevaluations/competencyid/{competencyid}/studentid/{studentid}/opnum/{opnum}/femail/{facultyemail}", createarowincompetencyevaluationsandsendform).Methods("GET")
	// r.HandleFunc("/competencyevaluationsdetails/competencyid/{competencyid}/studentid/{studentid}", evaluationformdetails).Methods("GET")

	r.HandleFunc("/competencyevaluations/competencyevaluationid/{competencyevaluationid}", postform).Methods("POST")
	r.HandleFunc("/competencyevaluations/facultyview/competencyid/{competencyid}/competencyevaluationid/{competencyevaluationid}", getfeedbackformwithsubmissiondetails).Methods("GET")

	r.HandleFunc("/facultytodo/meet/{email}", getfacultytodomeet).Methods("GET")
r.HandleFunc("/facultytodo/postmeet", postfacultytodomeet).Methods("POST")


	r.HandleFunc("/facultytodo/reference/{email}", facultytodoreference).Methods("GET")
	r.HandleFunc("/studentdashboard/details/studentmail/{email}", studentdashboarddetails).Methods("GET")
	r.HandleFunc("/studentdashboard/specialities", getspecnames).Methods("GET")
	r.HandleFunc("/studentdashboard/email/{email}/speciality/{speciality}", getstudentdashboardspecialitieswithcompetencies).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, r))

}
func postfacultytodomeet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	type Meet struct {
		Comeval_id int `json:"competencyevaluationid"`
		Meet string `json:"meettime"`
	}

	res := new(Meet)
	erro := json.NewDecoder(r.Body).Decode(&res)
	if erro != nil {
		panic(erro.Error())
	}

a := "update meet set meet_time=\""+res.Meet+"\",need_meet=0 where competency_evaluation_id=\""+strconv.Itoa(res.Comeval_id)+"\" and evaluation_type=\"self\";"
	fd, er := db.Query(a)
	if er != nil {

		panic(er.Error())
	}
	fd.Close()

	

	json.NewEncoder(w).Encode(res)

}
func getstudentdashboardspecialitieswithcompetencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}
	type CompetencyDetails struct {
		Competencynum  int     `json:"competencyid"`
		CompetencyName string  `json:"competencyname"`
		Self           float64 `json:"self"` // <-- CHANGED THIS LINE
		Faculty        float64 `json:"faculty"`
	}

	defer db.Close()

	strow, er := db.Query("select student_id from person p,student s where p.person_id=s.person_id and email=\"" + params["email"] + "\";")
	if er != nil {

		panic(er.Error())

	}
	defer strow.Close()
	student_id := ""
	for strow.Next() {

		err := strow.Scan(&student_id)

		if err != nil {
			panic(err)
		}

	}
	type Competencys struct {
		CompName string `json:"name"` // <-- CHANGED THIS LINE
		Compid   int    `json:"regno"`
	}

	comprow, er := db.Query("call getcompetencies(\"" + params["speciality"] + "\")")
	if er != nil {

		panic(er.Error())

	}
	defer comprow.Close()
	st := make([]*Competencys, 0)
	for comprow.Next() {
		user := new(Competencys)
		err := comprow.Scan(&user.CompName, &user.Compid)

		if err != nil {
			panic(err)
		}
		st = append(st, user)

	}

	compD := make([]*CompetencyDetails, 0)

	typef := "faculty"
	StudentF, er := db.Query("CALL getevalpercentageinstudentpage(\"" + params["speciality"] + "\",\"" + typef + "\",\"" + student_id + "\")")
	if er != nil {

		panic(er.Error())

	}

	type Score struct {
		Competency_Name string
		Competency_id   int
		Self            float64 // <-- CHANGED THIS LINE
		Faculty         float64
	}
	scores := make([]*Score, 0)

	for StudentF.Next() {
		sc := new(Score)
		err := StudentF.Scan(&sc.Faculty, &sc.Competency_id, &sc.Competency_Name)

		if err != nil {
			panic(err)
		}
		scores = append(scores, sc)
	}

	StudentF.Close()
	types := "self"
	StudentS, er := db.Query("CALL getevalpercentageinstudentpage(\"" + params["speciality"] + "\",\"" + types + "\",\"" + student_id + "\")")

	if er != nil {

		panic(er.Error())

	}

	for StudentS.Next() {
		var compid int
		var compname string
		var selfpercentage float64

		err := StudentS.Scan(&selfpercentage, &compid, &compname)

		if err != nil {
			panic(err)
		}
		for index, item := range scores {
			if item.Competency_id == compid {
				scores = append(scores[:index], scores[index+1:]...)
				scores = append(scores, &Score{Competency_Name: item.Competency_Name, Competency_id: item.Competency_id, Self: selfpercentage, Faculty: item.Faculty})

				break
			}

		}
	}

	StudentS.Close()

	for _, sitem := range st {
		fl := 0
		for _, item := range scores {
			if item.Competency_id == sitem.Compid {
				compD = append(compD, &CompetencyDetails{Self: item.Self, Faculty: item.Faculty, Competencynum: item.Competency_id, CompetencyName: item.Competency_Name})
				fl = 1
				break
			}
		}
		if fl == 0 {
			compD = append(compD, &CompetencyDetails{Self: 0, Faculty: 0, Competencynum: sitem.Compid, CompetencyName: sitem.CompName})
		}
	}

	json.NewEncoder(w).Encode(compD)
}

func getspecnames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {

		panic(err.Error())
	}

	rows, err := db.Query("call getspecialitys();")
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()

	type Result struct {
		SpecialityName string `json:"specialityName"`
		SpecialityId   int    `json:"SpecialityId "`
	}
	type Details struct {
		Comp []*Result `json:"details"`
	}

	res := make([]*Result, 0)
	for rows.Next() {
		rt := new(Result)
		err := rows.Scan(&rt.SpecialityName, &rt.SpecialityId)

		if err != nil {
			panic(err)
		}
		res = append(res, rt)
	}
	p := new(Details)
	p.Comp = make([]*Result, 0)
	for _, item := range res {
		p.Comp = append(p.Comp, &Result{SpecialityName: item.SpecialityName, SpecialityId: item.SpecialityId})
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(p)

}

func studentdashboarddetails(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//var competencyids []int=[]int{}

	fd, er := db.Query("select concat(p.first_name,p.last_name) from person p,student s where p.email=\"" + params["email"] + "\"and s.person_id=p.person_id;")
	if er != nil {

		panic(er.Error())
	}
	type Student struct {
		Name  string `json:"name"`
		Batch string `json:"batch"`
	}
	St := new(Student)
	for fd.Next() {

		err := fd.Scan(&St.Name)

		if err != nil {
			panic(err)

		}
	}
	fd.Close()
	ba, er := db.Query("CALL batch(\"" + params["email"] + "\");")
	if er != nil {

		panic(er.Error())
	}

	for ba.Next() {

		err := ba.Scan(&St.Batch)

		if err != nil {
			panic(err)

		}
	}
	ba.Close()

	json.NewEncoder(w).Encode(St)

}

func facultytodoreference(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	type Studentstomeet struct {
		Name                    string `json:"studentname"`
		Competency_Name         string `json:"competencyname"`
		Student_id              string `json:"studentid"`
		CriteriaQS              string `json:"criteriaqs"`
		CompetencyEvaluation_id int    `json:"CompetencyEvaluation_id"`
	}
	sts := make([]*Studentstomeet, 0)

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	de, er := db.Query("CALL todoreferenceforfaculty(\"" + params["email"] + "\");")
	if er != nil {

		panic(er.Error())

	}

	for de.Next() {
		st := new(Studentstomeet)
		err := de.Scan(&st.Competency_Name, &st.Student_id, &st.CriteriaQS, &st.CompetencyEvaluation_id, &st.Name)

		if err != nil {
			panic(err)

		}
		sts = append(sts, st)
	}
	de.Close()

	json.NewEncoder(w).Encode(sts)

}
func getfacultytodomeet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	type Studentstomeet struct {
		Name                    string `json:"studentname"`
		Competency_Name         string `json:"competencyname"`
		Student_id              string `json:"studentid"`
		CompetencyEvaluation_id int    `json:"CompetencyEvaluation_id"`
	}
	sts := make([]*Studentstomeet, 0)

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	de, er := db.Query("CALL todomeetforfaculty(\"" + params["email"] + "\");")
	if er != nil {

		panic(er.Error())

	}

	for de.Next() {
		st := new(Studentstomeet)
		err := de.Scan(&st.Competency_Name, &st.Student_id, &st.CompetencyEvaluation_id, &st.Name)

		if err != nil {
			panic(err)

		}
		sts = append(sts, st)
	}
	de.Close()

	json.NewEncoder(w).Encode(sts)

}

// func evaluationformdetails(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r) // Gets params

// 	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	defer db.Close()

// 	//var competencyids []int=[]int{}

// 	type Evaluationformdetails struct {
// 		EvaluationId int    `json:"compevaluationid"`
// 		Opnum        string `json:"patientop"`
// 		Date         string `json:"date"`
// 		Time         string `json:"time"`
// 		StudentName  string `json:"studentname"`
// 	}
// 	ev := new(Evaluationformdetails)

// 	de, er := db.Query("select competencyEvaluation_id from competency_evaluation where Student_Student_id=\""+ params["studentid"]+"\" and Competency_id=\""+params["competencyid"]+ "\"order by visit_stamp desc limit 1;")
// 	if er != nil {

// 		panic(er.Error())

// 	}
// 	var comeval_id int
// 	for de.Next() {

// 		err := de.Scan(&comeval_id)

// 		if err != nil {
// 			panic(err)

// 		}
// 	}
// 	de.Close()

// 	op, er := db.Query("call getfacultyfeedbackformdetails(\""+strconv.Itoa(comeval_id)+"\");");
// 	if er != nil {

// 		panic(er.Error())

// 	}

// 	for op.Next() {

// 		err := op.Scan(&ev.StudentName, &ev.Opnum, &ev.Date, &ev.Time)

// 		if err != nil {
// 			panic(err)

// 		}
// 	}
// 	op.Close()

// 	ev.EvaluationId = comeval_id

// 	json.NewEncoder(w).Encode(ev)

// }

// func createarowincompetencyevaluationsandsendform(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r) // Gets params

// 	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	defer db.Close()

// 	//var competencyids []int=[]int{}

// 	type Criteria struct {
// 		CriteriaId int    `json:"criteiaid"`
// 		CriteriaQs string `json:"criteriaqs"`
// 		Option0    string `json:"option0"`
// 		Option1    string `json:"option1"`
// 		Option2    string `json:"option2"`
// 	}
// 	type CriteriaOptions struct {
// 		CriteriaId int    `json:"criteiaid"`
// 		Option     string `json:"option"`
// 		OptVal     int
// 	}
// 	cr := make([]*Criteria, 0)
// 	cri, er := db.Query("call getcriteriasofcompetency(\""+params["competencyid"]+"\")")
// 	if er != nil {

// 		panic(er.Error())

// 	}

// 	for cri.Next() {
// 		cop := new(Criteria)
// 		err := cri.Scan(&cop.CriteriaId, &cop.CriteriaQs)

// 		if err != nil {
// 			panic(err)

// 		}
// 		cr = append(cr, cop)
// 	}
// 	cri.Close()

// 	co := make([]*CriteriaOptions, 0)
// 	opt, er := db.Query("call getcriteriaoptionsofcompetency(\""+params["competencyid"]+"\")")
// 	if er != nil {

// 		panic(er.Error())

// 	}

// 	for opt.Next() {
// 		cop := new(CriteriaOptions)
// 		err := opt.Scan(&cop.CriteriaId, &cop.Option, &cop.OptVal)

// 		if err != nil {
// 			panic(err)

// 		}
// 		co = append(co, cop)
// 	}
// 	opt.Close()
// 	for _, crit := range cr {
// 		for _, option := range co {
// 			if option.CriteriaId == crit.CriteriaId && option.OptVal == 0 {
// 				crit.Option0 = option.Option
// 			} else if option.CriteriaId == crit.CriteriaId && option.OptVal == 1 {
// 				crit.Option1 = option.Option
// 			} else if option.CriteriaId == crit.CriteriaId && option.OptVal == 2 {
// 				crit.Option2 = option.Option
// 			}
// 		}

// 	}

// 	json.NewEncoder(w).Encode(cr)

// }


func getfeedbackformwithsubmissiondetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	type Criteria struct {
		CriteriaId   int    `json:"criteiaid"`
		CriteriaQs   string `json:"criteriaqs"`
		Optionmatter string `json:"optionmatter"`
		Optval       int    `json:"optionval"`
		Refermatter  string `json:"refermatter"`
	}
	type CriteriaOptions struct {
		CriteriaId int    `json:"criteiaid"`
		Option     string `json:"option"`
		OptVal     int
	}
	type Evaluationformdetails struct {
		EvaluationId int    `json:"compevaluationid"`
		Opnum        string `json:"patientop"`
		Date         string `json:"date"`
		Time         string `json:"time"`
		StudentName  string `json:"studentname"`

		FacultyName string      `json:"facultyname"`
		Crit        []*Criteria `json:"criteriadetails"`
		Meet        string      `json:"meettime"`
	}

	ev := new(Evaluationformdetails)

	op, er := db.Query("call getfacultyfeedbackformdetails(\"" + params["competencyevaluationid"] + "\");")
	if er != nil {

		panic(er.Error())

	}

	for op.Next() {

		err := op.Scan(&ev.StudentName, &ev.Opnum, &ev.Date, &ev.Time)

		if err != nil {
			panic(err)

		}
	}
	op.Close()
	fa, er := db.Query("select concat(p.first_name,p.last_name) from competency_evaluation ce,person p,faculty f where ce.CompetencyEvaluation_id=\"" + params["competencyevaluationid"] + "\" and ce.Faculty_Faculty_id=f.faculty_id and f.person_id=p.person_id;")
	if er != nil {

		panic(er.Error())

	}
	fname := ""
	for fa.Next() {

		err := fa.Scan(&fname)

		if err != nil {
			panic(err)

		}
	}
	fa.Close()

	ev.EvaluationId, err = strconv.Atoi(params["competencyevaluationid"])

	ev.FacultyName = fname

	cr := make([]*Criteria, 0)
	cri, er := db.Query("call getcriteriasofcompetency(\"" + params["competencyid"] + "\")")
	if er != nil {

		panic(er.Error())

	}

	for cri.Next() {
		cop := new(Criteria)
		err := cri.Scan(&cop.CriteriaId, &cop.CriteriaQs)

		if err != nil {
			panic(err)

		}
		cr = append(cr, cop)
	}
	cri.Close()

	co := make([]*CriteriaOptions, 0)
	opt, er := db.Query("call getcriteriaoptionsofcompetency(\"" + params["competencyid"] + "\")")
	if er != nil {

		panic(er.Error())

	}

	for opt.Next() {
		cop := new(CriteriaOptions)
		err := opt.Scan(&cop.CriteriaId, &cop.Option, &cop.OptVal)

		if err != nil {
			panic(err)

		}
		co = append(co, cop)
	}
	opt.Close()
	type CriteriaScore struct {
		CriteriaId int `json:"criteiaid"`
		Optval     int `json:"optionval"`
	}
	csc := make([]*CriteriaScore, 0)
	opo, er := db.Query("select Criteria_id,Score_Type_Value from score where CompetencyEvaluation_id=\"" + params["competencyevaluationid"] + "\"and ScoreType_id=\"faculty\";")
	if er != nil {

		panic(er.Error())

	}

	for opo.Next() {
		cop := new(CriteriaScore)
		err := opo.Scan(&cop.CriteriaId, &cop.Optval)

		if err != nil {
			panic(err)

		}
		csc = append(csc, cop)
	}
	opo.Close()

	for _, crit := range cr {
		for _, option := range csc {
			if option.CriteriaId == crit.CriteriaId {
				crit.Optval = option.Optval
			}
		}

	}
	type CriteriaMatter struct {
		CriteriaId int `json:"criteiaid"`
		Matter     string
	}
	cm := make([]*CriteriaMatter, 0)
	opl, er := db.Query("select criteria_id,reference_matter from reference where evaluation_type=\"faculty\" and competency_evaluation_id=\"" + params["competencyevaluationid"] + "\";")
	if er != nil {

		panic(er.Error())

	}

	for opl.Next() {
		cop := new(CriteriaMatter)
		err := opl.Scan(&cop.CriteriaId, &cop.Matter)

		if err != nil {
			panic(err)

		}
		cm = append(cm, cop)
	}
	opl.Close()

	for _, crit := range cr {
		for _, option := range cm {
			if option.CriteriaId == crit.CriteriaId {
				crit.Refermatter = option.Matter
			}
		}

	}

	for _, crit := range cr {
		for _, option := range co {
			if option.CriteriaId == crit.CriteriaId && option.OptVal == crit.Optval {
				crit.Optionmatter = option.Option
			}
		}

	}

	ev.Crit = make([]*Criteria, 0)
	for _, item := range cr {
		ev.Crit = append(ev.Crit, &Criteria{CriteriaId: item.CriteriaId, CriteriaQs: item.CriteriaQs, Optionmatter: item.Optionmatter, Optval: item.Optval, Refermatter: item.Refermatter})
	}

	opl, er = db.Query("select meet_time from meet where competency_evaluation_id=\"" + params["competencyevaluationid"] + "\"and need_meet=0 and evaluation_type=\"faculty\";")
	if er != nil {

		panic(er.Error())

	}

	for opl.Next() {

		err := opl.Scan(&ev.Meet)

		if err != nil {
			panic(err)

		}

	}
	opl.Close()

	json.NewEncoder(w).Encode(ev)

}

func postform(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	type Form struct {
		Criteriaid  int    `json:"criteriaid"`
		Score       int    `json:"score"`
		Refermatter string `json:"matter"`
	}
	type Formwithmeet struct {
		CDetails []*Form `json:"criterias"`
		Meet     string  `json:"meettime"`
	}

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	var feedback Formwithmeet
	erro := json.NewDecoder(r.Body).Decode(&feedback)
	if erro != nil {
		panic(erro.Error())
	}

	for _, item := range feedback.CDetails {

		a := "call postform(\"" + strconv.Itoa(item.Criteriaid) + "\",\"" + params["competencyevaluationid"] + "\",\"" + strconv.Itoa(item.Score) + "\",\"" + item.Refermatter + "\");"
		fd, er := db.Query(a)
		if er != nil {

			panic(er.Error())
		}
		fd.Close()
	}
	flty := "faculty"
	a := "call insertmeettime(\"" + feedback.Meet + "\",\"" + params["competencyevaluationid"] + "\",\"" + flty + "\");"
	fd, er := db.Query(a)
	if er != nil {

		panic(er.Error())
	}
	fd.Close()

	json.NewEncoder(w).Encode(feedback)
}

func getfeedbackform(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//var competencyids []int=[]int{}
	type Criteria struct {
		CriteriaId int    `json:"criteiaid"`
		CriteriaQs string `json:"criteriaqs"`
		Option0    string `json:"option0"`
		Option1    string `json:"option1"`
		Option2    string `json:"option2"`
	}
	type CriteriaOptions struct {
		CriteriaId int    `json:"criteiaid"`
		Option     string `json:"option"`
		OptVal     int
	}
	type Evaluationformdetails struct {
		EvaluationId int    `json:"compevaluationid"`
		Opnum        string `json:"patientop"`
		Date         string `json:"date"`
		Time         string `json:"time"`
		StudentName  string `json:"studentname"`

		FacultyName string      `json:"facultyname"`
		Crit        []*Criteria `json:"criteriadetails"`
	}

	ev := new(Evaluationformdetails)

	// de, er := db.Query("select competencyEvaluation_id from competency_evaluation where Student_Student_id=\""+ params["studentid"]+"\" and Competency_id=\""+params["competencyid"]+ "\"order by visit_stamp desc limit 1;")
	// if er != nil {

	// 	panic(er.Error())

	// }
	// var comeval_id int
	// for de.Next() {

	// 	err := de.Scan(&comeval_id)

	// 	if err != nil {
	// 		panic(err)

	// 	}
	// }
	// de.Close()

	op, er := db.Query("call getfacultyfeedbackformdetails(\"" + params["competencyevaluationid"] + "\");")
	if er != nil {

		panic(er.Error())

	}

	for op.Next() {

		err := op.Scan(&ev.StudentName, &ev.Opnum, &ev.Date, &ev.Time)

		if err != nil {
			panic(err)

		}
	}
	op.Close()
	fa, er := db.Query("select concat(p.first_name,p.last_name) from competency_evaluation ce,person p,faculty f where ce.CompetencyEvaluation_id=\"" + params["competencyevaluationid"] + "\" and ce.Faculty_Faculty_id=f.faculty_id and f.person_id=p.person_id;")
	if er != nil {

		panic(er.Error())

	}
	fname := ""
	for fa.Next() {

		err := fa.Scan(&fname)

		if err != nil {
			panic(err)

		}
	}
	fa.Close()

	ev.EvaluationId, err = strconv.Atoi(params["competencyevaluationid"])

	ev.FacultyName = fname

	cr := make([]*Criteria, 0)
	cri, er := db.Query("call getcriteriasofcompetency(\"" + params["competencyid"] + "\")")
	if er != nil {

		panic(er.Error())

	}

	for cri.Next() {
		cop := new(Criteria)
		err := cri.Scan(&cop.CriteriaId, &cop.CriteriaQs)

		if err != nil {
			panic(err)

		}
		cr = append(cr, cop)
	}
	cri.Close()

	co := make([]*CriteriaOptions, 0)
	opt, er := db.Query("call getcriteriaoptionsofcompetency(\"" + params["competencyid"] + "\")")
	if er != nil {

		panic(er.Error())

	}

	for opt.Next() {
		cop := new(CriteriaOptions)
		err := opt.Scan(&cop.CriteriaId, &cop.Option, &cop.OptVal)

		if err != nil {
			panic(err)

		}
		co = append(co, cop)
	}
	opt.Close()
	for _, crit := range cr {
		for _, option := range co {
			if option.CriteriaId == crit.CriteriaId && option.OptVal == 0 {
				crit.Option0 = option.Option
			} else if option.CriteriaId == crit.CriteriaId && option.OptVal == 1 {
				crit.Option1 = option.Option
			} else if option.CriteriaId == crit.CriteriaId && option.OptVal == 2 {
				crit.Option2 = option.Option
			}
		}

	}

	ev.Crit = make([]*Criteria, 0)
	for _, item := range cr {
		ev.Crit = append(ev.Crit, &Criteria{CriteriaId: item.CriteriaId, CriteriaQs: item.CriteriaQs, Option0: item.Option0, Option1: item.Option1, Option2: item.Option2})
	}
	json.NewEncoder(w).Encode(ev)

}

func addroweval(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	type Result struct {
		Opnum string `json:"opnum"`
		Fmail string `json:"fmail"`
	}

	res := new(Result)
	erro := json.NewDecoder(r.Body).Decode(&res)
	if erro != nil {
		panic(erro.Error())
	}

	fd, er := db.Query("select p.person_id,faculty_id from  faculty f,person p where p.person_id=f.person_id and p.email=\"" + res.Fmail + "\";")
	if er != nil {

		panic(er.Error())

	}
	var faculty_id string
	var person_id int

	for fd.Next() {

		err := fd.Scan(&person_id, &faculty_id)

		if err != nil {
			panic(err)

		}
	}
	fd.Close()

	insert, er := db.Query("call createevaluationrow(\"" + params["competencyid"] + "\",\"" + params["studentid"] + "\",\"" + strconv.Itoa(person_id) + "\",\"" + faculty_id + "\",\"" + res.Opnum + "\");")
	if er != nil {

		panic(er.Error())

	}
	insert.Close()

	json.NewEncoder(w).Encode(res)

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
		CompEvaId int     `json:"compentencyevaluationid"`
		Opnum     string  `json:"patientop"` // <-- CHANGED THIS LINE
		Date      string  `json:"date"`
		Time      string  `json:"time"`
		Self      float64 `json:"self"`
		Faculty   float64 `json:"faculty"`
		Timest    string  `json:"-"`
	}

	evalrow, er := db.Query("call getallevalofacompetency(\"" + params["competencyid"] + "\",\"" + params["studentid"] + "\");")
	if er != nil {

		panic(er.Error())

	}
	defer evalrow.Close()
	et := []Evaluation{}

	for evalrow.Next() {
		user := new(Evaluation)
		err := evalrow.Scan(&user.CompEvaId, &user.Opnum, &user.Date, &user.Time)

		if err != nil {
			panic(err)

		}
		datab, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
		if err != nil {
			panic(err.Error())
		}
		typef := "faculty"
		StudentF, er := datab.Query("CALL getpercentageforeacheval(\"" + typef + "\",\"" + params["competencyid"] + "\",\"" + strconv.Itoa(user.CompEvaId) + "\");")
		if er != nil {

			panic(er.Error())

		}

		for StudentF.Next() {

			err := StudentF.Scan(&user.Faculty)

			if err != nil {
				panic(err)
			}

		}

		StudentF.Close()
		types := "self"
		StudentS, er := datab.Query("CALL getpercentageforeacheval(\"" + types + "\",\"" + params["competencyid"] + "\",\"" + strconv.Itoa(user.CompEvaId) + "\");")

		if er != nil {

			panic(er.Error())

		}

		for StudentS.Next() {

			err := StudentS.Scan(&user.Self)

			if err != nil {
				panic(err)
			}

		}

		StudentS.Close()
		datab.Close()
		user.Timest = user.Date + " " + user.Time
		et = append(et, *user)

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

	rows, err := db.Query("call getprofile(\"" + params["email"] + "\")")
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

			row, err := db.Query("call getbatch(\"" + params["email"] + "\")")
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
	type StudentDetails struct {
		Name           string  `json:"name"` // <-- CHANGED THIS LINE
		Adno           string  `json:"regno"`
		Self           float64 `json:"self"` // <-- CHANGED THIS LINE
		Faculty        float64 `json:"faculty"`
		Competencynum  int     `json:"competencynum"`
		CompetencyName string  `json:"competencyname"`
	}

	defer db.Close()

	//var competencyids []int=[]int{}

	type Students struct {
		Name string `json:"name"` // <-- CHANGED THIS LINE
		Adno string `json:"regno"`
	}

	studentrow, er := db.Query("call getstudents()")
	if er != nil {

		panic(er.Error())

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
	typef := "faculty"
	StudentF, er := db.Query("CALL getevalpercentage(\"" + params["speciality"] + "\",\"" + typef + "\",\"" + params["competencyid"] + "\")")
	if er != nil {

		panic(er.Error())

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
	types := "self"
	StudentS, er := db.Query("CALL getevalpercentage(\"" + params["speciality"] + "\",\"" + types + "\",\"" + params["competencyid"] + "\")")

	if er != nil {

		panic(er.Error())

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
	/*
	    for _ , item := range st {
	   	fmt.Printf("%+v\n",item)
	   	}
	*/
	compid, err := strconv.Atoi(params["competencyid"])
	if err != nil {
		panic(err.Error())
	}
	for _, sitem := range st {
		fl := 0
		for _, item := range scores {
			if item.Adno == sitem.Adno {
				studentD = append(studentD, &StudentDetails{Name: sitem.Name, Adno: item.Adno, Self: item.Self, Faculty: item.Faculty, Competencynum: item.Competency_id})
				fl = 1
				break
			}
		}
		if fl == 0 {
			studentD = append(studentD, &StudentDetails{Name: sitem.Name, Adno: sitem.Adno, Self: 0, Faculty: 0, Competencynum: compid})
		}
	}

	rows, err := db.Query("select Competency_Name,competency_id from competency where Speciality_id in ( select Speciality_id from speciality where Speciality_Name=\"" + params["speciality"] + "\");")
	if err != nil {

		panic(err.Error())

	}

	//var competencyids []int=[]int{}

	type Competency struct {
		Name string `json:"name"` // <-- CHANGED THIS LINE
		Cid  int    `json:"cid"`
	}
	comp := make([]*Competency, 0)
	for rows.Next() {
		onec := new(Competency)
		err := rows.Scan(&onec.Name, &onec.Cid)

		if err != nil {
			panic(err)
		}
		comp = append(comp, onec)

		for _, item := range comp {
			for _, sitem := range studentD {
				if item.Cid == sitem.Competencynum {
					sitem.CompetencyName = item.Name
				}

			}
		}

	}
	rows.Close()

	// Compre := make([]*CompetencyReturn, 0)
	// for _, stud := range studentD {
	// 	stude := StudentDetails{Name: stud.Name, Adno: stud.Adno, Self: stud.Self, Faculty: stud.Faculty, Competencynum: stud.Competencynum}

	// 	Compre = append(Compre, &CompetencyReturn{C: Competen{StudentDetails: stude}})
	// }

	json.NewEncoder(w).Encode(studentD)
}

// func (c CompetencyReturn) MarshalJSON() ([]byte, error) {
// 	// encode the original
// 	m, _ := json.Marshal(c.C)

// 	// decode it back to get a map
// 	var a interface{}
// 	json.Unmarshal(m, &a)
// 	b := a.(map[string]interface{})

// 	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	rows, err := db.Query("select Competency_Name,competency_id from competency where Speciality_id in ( select Speciality_id from speciality where Speciality_Name=?);", speciality_for_faculty)
// 	if err != nil {

// 		panic(err.Error())

// 	}

// 	defer db.Close()

// 	//var competencyids []int=[]int{}

// 	type Competency struct {
// 		Name string  `json:"name"` // <-- CHANGED THIS LINE
// 		Cid  float64 `json:"cid"`
// 	}
// 	comp := make([]*Competency, 0)
// 	for rows.Next() {
// 		onec := new(Competency)
// 		err := rows.Scan(&onec.Name, &onec.Cid)

// 		if err != nil {
// 			panic(err)
// 		}
// 		comp = append(comp, onec)

// 	}
// 	defer rows.Close()

// 	for i, si := range b {
// 		var f interface{}
// 		n, _ := json.Marshal(b[i])
// 		json.Unmarshal(n, &f)
// 		c := f.(map[string]interface{})
// 		//idx := string(c["id"])

// 		idx := c["competencynum"].(float64)
// 		for _, co := range comp {

// 			if co.Cid == idx {
// 				b[co.Name] = si

// 				delete(b, "competency")
// 			}
// 		}

// 	}

// 	return json.Marshal(b)

// }
func getcompnames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {

		panic(err.Error())
	}

	rows, err := db.Query("call getcompetencies(\"" + params["speciality"] + "\");")
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()

	type Result struct {
		CompetencyName string `json:"competencyname"`
		CompetencyId   int    `json:"competencyid"`
	}
	type Details struct {
		Comp []*Result `json:"details"`
	}

	res := make([]*Result, 0)
	for rows.Next() {
		rt := new(Result)
		err := rows.Scan(&rt.CompetencyName, &rt.CompetencyId)

		if err != nil {
			panic(err)
		}
		res = append(res, rt)
	}
	p := new(Details)
	p.Comp = make([]*Result, 0)
	for _, item := range res {
		p.Comp = append(p.Comp, &Result{CompetencyName: item.CompetencyName, CompetencyId: item.CompetencyId})
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(p)

}

// func getcompnames(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r) // Gets params

// 	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	rows, err := db.Query("call getcompetencies(?)", params["speciality"])
// 	if err != nil {

// 		panic(err.Error())

// 	}
// 	defer db.Close()

// 	type Result struct {
// 		Competency []string `json:"competency"`

// 	}

// 	res := Result{Competency: []string{}}
// 	for rows.Next() {
// 		var str string
// 		err := rows.Scan(&str)

// 		if err != nil {
// 			panic(err)
// 		}
// 		res.Competency = append(res.Competency, str)
// 		res = Result{Competency: res.Competency}

// 	}
// 	defer rows.Close()

// 	json.NewEncoder(w).Encode(res)

// }
func getfacultydetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	db, err := sql.Open("mysql", "b43dbfed48dc1d:395f6a59@tcp(us-cdbr-east-05.cleardb.net)/heroku_ae8d9f2c5bc1ed0")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select concat(p.first_name,p.last_name) as name,f.speciality from person p,faculty f where p.person_id=f.person_id and p.email=\"" + params["email"] + "\";")
	if err != nil {

		panic(err.Error())

	}
	defer db.Close()
	type Result struct {
		Name       string `json:"name"`
		Speciality string `json:"speciality"`
	}
	res := new(Result)
	for rows.Next() {
		user := new(Result)
		err := rows.Scan(&user.Name, &user.Speciality)

		if err != nil {
			panic(err)
		}
		res = user
	}
	defer rows.Close()
	type Details struct {
		Rest Result `json:"details"`
	}

	json.NewEncoder(w).Encode(Details{Rest: *res})

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

}
