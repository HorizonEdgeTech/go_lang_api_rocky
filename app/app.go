package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/isaac/model"
	"github.com/isaac/routes"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "isaac"
	password = "your-password"
	dbname   = "rocky_db"
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
	defer a.DB.Close()

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
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) postLoanRepayment(w http.ResponseWriter, r *http.Request) {

}

func (a *App) postRiderDetails(w http.ResponseWriter, r *http.Request) {

}
func (a *App) getLoanRepayments(w http.ResponseWriter, r *http.Request) {
	var riders []model.Rider
	query := "SELECT * FROM loan_repayments"
	rows, err := a.DB.Query(query)
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
	query := "SELECT * FROM applicants WHERE id_number=$1"
	row := a.DB.QueryRow(query, id)
	err := row.Scan(&rider.IDNumber, &rider.Surname, &rider.OtherName, &rider.MobileNumber, &rider.DOB)
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
	query := "SELECT * FROM applicants"
	rows, err := a.DB.Query(query)
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
