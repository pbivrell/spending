package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/pbivrell/gatekeeper/log"
	"github.com/pbivrell/gatekeeper/server"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var wait time.Duration
var influxAddr string
var sqlLitePath string
var proxyPath string

func init() {
	flag.StringVar(&proxyPath, "proxy-config", "", "json file describing server addrs and roles to proxy to")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&server.JwtKey, "jwt-secret", "", "jwt secret for generating jwts")
	flag.StringVar(&influxAddr, "influx-addr", "", "influx address to log to")
	flag.StringVar(&sqlLitePath, "sqllite", "./test.db", "sql lite file path to store credentials")
	flag.Parse()

}

func main() {

	creds, err := sqlCreds()
	if err != nil {
		panic(err)
	}

	err = initCreds(creds)
	if err != nil {
		fmt.Printf("Failed to init creds: %v\n", err)
	}

	var logger log.Logger
	logger = log.NewLogrusLogger()
	if influxAddr != "" {
		logger = log.NewInfluxLogger(influxdb2.NewClient(influxAddr, ""))
	}

	lock := &sync.Mutex{}
	proxys, err := server.LoadProxies(proxyPath)
	if err != nil {
		fmt.Printf("Failed to load proxy(s): %v\n", err)
	}

	s := server.New(creds)

	r := mux.NewRouter()

	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	r.HandleFunc("/api/v1/login", server.LogMiddlewear(logger, s.Login)).
		Methods("POST")
	r.HandleFunc("/api/v1/refresh", server.LogMiddlewear(logger, server.AuthMiddlewear(server.ReadRole, server.Refresh))).
		Methods("GET")
	r.HandleFunc("/api/v1/users", server.LogMiddlewear(logger, server.AuthMiddlewear(server.AdminRole, s.Users))).
		Methods("GET")
	r.HandleFunc("/api/v1/user/{user}", server.LogMiddlewear(logger, server.AuthMiddlewear(server.ReadRole, s.User))).
		Methods("GET")
	r.HandleFunc("/api/v1/user/{user}", server.LogMiddlewear(logger, server.AuthMiddlewear(server.ReadRole, s.DeleteUser))).
		Methods("DELETE")
	r.HandleFunc("/api/v1/user", server.LogMiddlewear(logger, server.AuthMiddlewear(server.AdminRole, s.CreateUser))).
		Methods("POST")
	r.HandleFunc("/api/v1/user", server.LogMiddlewear(logger, server.AuthMiddlewear(server.ReadRole, s.GetUser))).
		Methods("GET")

	r.HandleFunc("/api/v1/user", server.LogMiddlewear(logger, server.AuthMiddlewear(server.ReadRole, s.GetUser))).
		Methods("GET")

	r.HandleFunc(`/api/v1/proxy/{proxy:[A-z]+}/{path:[a-zA-Z0-9=\-\/_]*}`, server.LogMiddlewear(logger, s.Proxy(lock, proxys)))

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Failed serving", err)
		}
	}()

	// Handle reload to update proxies
	p := make(chan os.Signal, 1)
	signal.Notify(p, syscall.SIGHUP)
	go func(lock *sync.Mutex, proxys *map[string]server.Proxy) {
		<-p
		ps, err := server.LoadProxies(proxyPath)
		if err != nil {
			fmt.Printf("Failed to load proxy(s): %v\n", err)
		}
		lock.Lock()
		proxys = &ps
		lock.Unlock()
	}(lock, &proxys)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	fmt.Println("shutting down")
	os.Exit(0)
}

func sqlCreds() (server.UserdataStorage, error) {

	_, err := os.Stat(sqlLitePath)
	if os.IsNotExist(err) {
		_, err := os.Create(sqlLitePath)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", sqlLitePath)
	if err != nil {
		return nil, err
	}

	return server.NewSQLUserdataStorage(db)
}

func initCreds(storage server.UserdataStorage) error {

	defaultPassword := os.Getenv("GATEKEEPER_ADMIN_PASSWORD")
	if defaultPassword == "" {
		return fmt.Errorf("no default password provided")
	}

	password, err := server.HashPassword(defaultPassword)
	if err != nil {
		return err
	}

	return storage.Insert(server.Userdata{
		Username: "owner",
		Password: password,
		Email:    "",
		Role:     server.OwnerRole,
	})
}
