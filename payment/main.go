package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var balance = 0

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/payment/{orderId}/{amount}", payment)
	r.HandleFunc("/get-balance", getBalance)
	log.Println("starting server")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}

	log.Fatal(srv.ListenAndServe())
}

func payment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	amount := vars["amount"]

	v, _ := strconv.Atoi(amount)
	balance = balance + v

	w.WriteHeader(http.StatusOK)
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"balance":%v}`, balance)))
	// fmt.Println(balance)
}
