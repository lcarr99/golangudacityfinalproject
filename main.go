package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	errorResponses "github.com/lcarr99/cd11970-Go-screencast-code/Final_Project/http/responses/error"
	"github.com/lcarr99/cd11970-Go-screencast-code/Final_Project/modules/customers"

	"github.com/gorilla/mux"
)

var connection *sql.DB
var errorResponse errorResponses.ErrorResponse

func getCustomer(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, invalidIdError := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if invalidIdError != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(errorResponses.ErrorResponse{Message: "Please ensure the id is numeric"})
		return
	}

	CustomerRepository := customers.CustomerRepository{DB: connection}

	customerStruct, err := CustomerRepository.OfId(id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			errorResponse.Message = "Customer was not found"
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse.Message = err.Error()
		}

		encoder.Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(customerStruct)
}

func getAllCustomers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	customerRepository := customers.CustomerRepository{DB: connection}

	customersSlice, err := customerRepository.All()

	if err != nil {
		errorResponse.Message = err.Error()
		encoder.Encode(errorResponse)
		return
	}

	encoder.Encode(customersSlice)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if requestData, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Message = "Invalid data passed"
		encoder.Encode(errorResponse)
	} else {
		customerStruct := customers.Customer{}
		if err := json.Unmarshal(requestData, &customerStruct); err != nil {
			errorResponse.Message = "Invalid data passed"
			encoder.Encode(errorResponse)
			return
		}

		customerRepository := customers.CustomerRepository{DB: connection}
		err := customerRepository.Create(&customerStruct)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse.Message = err.Error()
			encoder.Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		encoder.Encode(customerStruct)
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	id, error := strconv.Atoi(vars["id"])

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Message = "Please ensure the id is numeric"
		encoder.Encode(errorResponse)
		return
	}

	customerRepository := customers.CustomerRepository{DB: connection}
	customerStruct, err := customerRepository.OfId(id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			errorResponse.Message = "Customer not found"
			encoder.Encode(errorResponse)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse.Message = err.Error()
			encoder.Encode(errorResponse)
			return
		}
	}

	err = customerRepository.Delete(customerStruct)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse.Message = err.Error()
		encoder.Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Message = "Please ensure id is numeric"
		encoder.Encode(errorResponse)
		return
	}

	requestData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Message = "Invalid data given"
		encoder.Encode(errorResponse)
		return
	}

	customerRepository := customers.CustomerRepository{DB: connection}

	customerStruct, err := customerRepository.OfId(id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			errorResponse.Message = "Customer not found"
			encoder.Encode(errorResponse)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse.Message = err.Error()
			encoder.Encode(errorResponse)
			return
		}
	}

	if err := json.Unmarshal(requestData, &customerStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Message = "Invalid data given"
		encoder.Encode(errorResponse)
		return
	}

	err = customerRepository.Update(customerStruct)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse.Message = err.Error()
		encoder.Encode(errorResponse)
		return
	}

	encoder.Encode(customerStruct)
}

func main() {
	router := mux.NewRouter()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Please ensure you have a .env set")
	}

	var address = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	config := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   address,
		DBName: os.Getenv("DB_DATABASE")}

	var error error

	connection, error = sql.Open("mysql", config.FormatDSN())

	if error != nil {
		log.Fatalf("Error when connection to the database")
	}

	router.HandleFunc("/customers/{id:[0-9]+}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id:[0-9]+}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id:[0-9]+}", updateCustomer).Methods("PATCH")
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe(":3000", router)
}
