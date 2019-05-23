package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"
	"github.com/gorilla/mux"
	"github.com/isaac/model"
	"github.com/isaac/routes"
	_ "github.com/lib/pq"
)

const sql_create = `
CREATE TABLE applicant(
    id_number INT PRIMARY KEY,
    surname VARCHAR(20) NOT NULL,
    other_name VARCHAR(50) NOT NULL,
    nationality VARCHAR(20) DEFAULT 'Kenyan',
    no_of_dependents INT DEFAULT 0,
    mobile_number VARCHAR(20),
    alternative_number VARCHAR(20)
);

CREATE TABLE applicant_address_details(
    postal_address int primary key,
    city VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL,
    physical_address VARCHAR(30) NOT NULL,
    house_no VARCHAR(20),
    road VARCHAR(80),
    town VARCHAR(50) NOT NULL,
    period_at_current_address int , --in years
    rented NUMERIC(1) NOT NULL,
    id_number int REFERENCES applicant(id_number)
);

CREATE TABLE loan(
    loan_id INT PRIMARY KEY,
    loan_amount NUMERIC(6, 2),
    purpose TEXT NOT NULL,
    interest_rate NUMERIC(5, 5) NOT NULL DEFAULT 5,
    repayment_period int not null,
    id_number INT REFERENCES applicant(id_number)
);

CREATE TABLE guarantor(
    id_number INT PRIMARY KEY,
    name text not null,
    mobile_number VARCHAR(13)
);

CREATE TABLE guaranteed_loans(
    applicant_id int REFERENCES applicant(id_number),
    guarantor_id INT REFERENCES guarantor(id_number),
    reationship TEXT NOT NULL,
    PRIMARY KEY(applicant_id, guarantor_id)
);

CREATE TABLE bank(
    bank_name VARCHAR(80), CHECK(bank_name IN ('Equity bank', 'Cooperative Bank')),--to-be-completed 
    PRIMARY KEY(bank_name)   
);

CREATE TABLE account_details(
    id int REFERENCES applicant(id_number),
    bank_name VARCHAR(80)  REFERENCES bank(bank_name),
    PRIMARY KEY (id, bank_name)
);

CREATE TABLE next_of_kin(
    id_number INT PRIMARY KEY,
    surname VARCHAR(50) NOT NULL,
    other_name VARCHAR(80),
    occupation VARCHAR(50),
    place_of_work VARCHAR(50),
    mobile_number VARCHAR(13) NOT NULL
);

CREATE TABLE applicant_relationship(
    applicant_id_number INT REFERENCES applicant(id_number),
    nok_id_number INT REFERENCES next_of_kin(id_number), --next of kin ID number
    relationship VARCHAR(50) NOT NULL,
    PRIMARY KEY(applicant_id_number, nok_id_number)
);

CREATE TABLE loan_repayment(
    loan_id INT REFERENCES loan(loan_id),
    paid_amount NUMERIC(9, 2) NOT NULL,
    date_paid DATE NOT NULL,
    paid_by VARCHAR(80) NOT NULL,
    PRIMARY KEY(loan_id)
);

CREATE TABLE loan_defaulters(
    loan_id INT REFERENCES loan(loan_id),
    amount_due NUMERIC(9, 2),
    days_overdue INT NOT NULL
);
`

const (
	host     = "localhost"
	port     = 5432
	user     = "isaac"
	password = "toor@#()2390"
	dbname   = "rocky_boda"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) InitializeRoutes() {
	var routesAvailable = routes.Routes{
		//Some get methods
		routes.Route{
			"RiderDetails",
			"GET",
			"/rider/{id:[0-9]+}",
			a.getRiderDetails,
		},
		routes.Route{
			"GetAllRiders",
			"GET",
			"/rider",
			a.getAllRiders,
		},
		routes.Route{
			"GetLoanDefaulter",
			"GET",
			"/loanDefaulters",
			a.getLoanDefaulters,
		},
		routes.Route{
			"GetLoanRepayments",
			"GET",
			"/loanRepayments",
			a.getLoanRepayments,
		},
		//Some post methods
		routes.Route{
			"PostRiderDetails",
			"POST",
			"/postRiderDetails",
			a.postRiderDetails,
		},
		routes.Route{
			"PostLoanRepayment",
			"POST",
			"/postLoanRepayment",
			a.postLoanRepayment,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, item := range routesAvailable {
		var handler http.Handler
		handler = item.HandlerFunc
		handler = Logger(handler, item.Name)
		router.Methods(item.Method).Path(item.Pattern).Name(item.Name).Handler(handler)
	}
	a.Router = router
}

func (a *App) Initialize() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	a.DB = db

	err = a.DB.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

//Log all request
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":"+addr, a.Router))
}

func (a *App) postLoanRepayment(w http.ResponseWriter, r *http.Request) {
	//TODO
	respondWithJson(w)
}

func (a *App) postRiderDetails(w http.ResponseWriter, r *http.Request) {
	//TODO
	respondWithJson(w)
}
func (a *App) getLoanRepayments(w http.ResponseWriter, r *http.Request) {
	var riders []model.Rider
	query := "SELECT * FROM loan_repayments"
	rows, err := a.DB.Query(query)
	respondWithJson(w)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var rider model.Rider
		err = rows.Scan() //TODO
		error := map[string]string{}
		if err != nil {
			json.NewEncoder(w).Encode(&error)
			return
		}
		riders = append(riders, rider)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&riders)
}
func (a *App) getRiderDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var rider model.Rider
	query := "SELECT * FROM applicant WHERE id_number=$1"
	respondWithJson(w)
	row := a.DB.QueryRow(query, id)
	err := row.Scan(&rider.IDNumber, &rider.Surname, &rider.OtherName, &rider.Nationality, &rider.NoOfDependents, &rider.MobileNumber, &rider.AlternativeNumber)
	switch err {
	case sql.ErrNoRows:
		error := map[string]string{"error": "Rider doesn't exist"}
		if err != nil {
			json.NewEncoder(w).Encode(&error)
			return
		}
		return
	case nil:
		json.NewEncoder(w).Encode(&rider)
	default:
		panic(err)
	}

}
func (a *App) getLoanDefaulters(w http.ResponseWriter, r *http.Request) {
	var riders []model.Rider
	query := "SELECT * FROM loan_defaulters"
	rows, err := a.DB.Query(query)
	if err != nil {
		panic(err)
	}
	respondWithJson(w)
	defer rows.Close()
	for rows.Next() {
		var rider model.Rider
		err = rows.Scan() //TODO
		error := map[string]string{}
		if err != nil {
			json.NewEncoder(w).Encode(&error)
			return
		}
		riders = append(riders, rider)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&riders)
}
func (a *App) getAllRiders(w http.ResponseWriter, r *http.Request) {
	var riders []model.Rider
	query := "SELECT * FROM applicant"
	rows, err := a.DB.Query(query)
	if err != nil {
		panic(err)
	}
	// w.Header().Add("Content-Type", "application/json")
	respondWithJson(w)
	defer rows.Close()
	for rows.Next() {
		var rider model.Rider
		err = rows.Scan(&rider.IDNumber, &rider.Surname, &rider.OtherName, &rider.Nationality, &rider.NoOfDependents, &rider.MobileNumber, &rider.AlternativeNumber) //TODO
		error := map[string]string{}//For error Handling
		if err != nil {
			json.NewEncoder(w).Encode(&error)
			return
		}
		riders = append(riders, rider)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&riders)
}
func (a *App) Close(){
	defer a.DB.Close()
}

func respondWithJson(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}