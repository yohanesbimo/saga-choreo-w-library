package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itimofeev/go-saga"
)

var order = map[string]bool{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create/{name}", create)
	r.HandleFunc("/get-order", getOrder)
	log.Println("starting server")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8081",
	}

	log.Fatal(srv.ListenAndServe())
}

func create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	s := saga.NewSaga("testing")

	s.AddStep(&saga.Step{
		Name: "order",
		Func: func(ctx context.Context) error {
			order[name] = false

			return nil
		},
		CompensateFunc: func(ctx context.Context) error {
			delete(order, name)

			return nil
		},
	})

	s.AddStep(&saga.Step{
		Name: "payment",
		Func: func(ctx context.Context) error {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/payment/%s/%d", name, 1000))
			if err != nil {
				return err
			}

			if resp.StatusCode != http.StatusOK {
				return errors.New(resp.Status)
			}

			order[name] = true

			return nil
		},
		CompensateFunc: func(ctx context.Context) error {
			order[name] = false

			return nil
		},
	})

	ctx := context.Background()
	store := saga.New()
	c := saga.NewCoordinator(ctx, ctx, s, store, "test")

	result := c.Play()

	if result.ExecutionError != nil {
		log.Println("ERROR:", result.ExecutionError.Error())
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf(`{"id": %v}`, len(order)-1)))
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	b, _ := json.Marshal(order)

	w.Write(b)
	// fmt.Println(balance)
}
