package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/pbivrell/spending/server"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	// Add your routes as needed

	s := &server.Server{}

	r.HandleFunc("/api/1.0/user", s.GetUser).Methods("GET")
	r.HandleFunc("/api/1.0/user", s.CreateUser).Methods("POST")
	r.HandleFunc("/api/1.0/user", s.UpdateUser).Methods("PATCH")
	r.HandleFunc("/api/1.0/user", s.DeleteUser).Methods("DELETE")

	r.HandleFunc("/api/1.0/page", s.GetPage).Methods("GET")
	r.HandleFunc("/api/1.0/page", s.CreatePage).Methods("POST")
	r.HandleFunc("/api/1.0/page", s.UpdatePage).Methods("PATCH")
	r.HandleFunc("/api/1.0/page", s.DeletePage).Methods("DELETE")

	r.HandleFunc("/api/1.0/estimate", s.GetEstimate).Methods("GET")
	r.HandleFunc("/api/1.0/estimate", s.CreateEstimate).Methods("POST")
	r.HandleFunc("/api/1.0/estimate", s.UpdateEstimate).Methods("PATCH")
	r.HandleFunc("/api/1.0/estimate", s.DeleteEstimate).Methods("DELETE")

	r.HandleFunc("/api/1.0/transaction", s.GetTransaction).Methods("GET")
	r.HandleFunc("/api/1.0/transaction", s.CreateTransaction).Methods("POST")
	r.HandleFunc("/api/1.0/transaction", s.UpdateTransaction).Methods("PATCH")
	r.HandleFunc("/api/1.0/transaction", s.DeleteTransaction).Methods("DELETE")

	r.HandleFunc("/api/1.0/role", s.GetRole).Methods("GET")
	r.HandleFunc("/api/1.0/role", s.CreateRole).Methods("POST")
	r.HandleFunc("/api/1.0/role", s.UpdateRole).Methods("PATCH")
	r.HandleFunc("/api/1.0/role", s.DeleteRole).Methods("DELETE")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
