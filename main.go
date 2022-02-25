package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hrasiter/go-mongodb-docker/services"
	"github.com/joho/godotenv"
)

const (
	addr = "8080"
	get  = "GET"
	post = "POST"
)

func main() {

	enverr := godotenv.Load("environment.env")
	if enverr != nil {
		log.Fatal("invalid environment file", enverr.Error())
		return
	}

	// create anew http router
	rtr := mux.NewRouter()

	// health endpoint
	rtr.HandleFunc("/health", services.Health).Methods(get)
	rtr.HandleFunc("/fill", services.Fill).Methods(get)

	// use go routinue to serve endpoint
	ctx := context.Background()

	GracefullyListenAndServe(ctx, addr, rtr)
}

func GracefullyListenAndServe(ctx context.Context, servePort string, rtr *mux.Router) {
	http.Handle("/", rtr)

	h := &http.Server{
		Addr:    fmt.Sprintf(":%v", servePort),
		Handler: handlers.CORS()(rtr),
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		log.Printf("serving on port: %v", servePort)
		if err := h.ListenAndServe(); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	// wait for signal to end
	<-sig

	log.Println("Shutting down server...")

	_ = h.Shutdown(ctx)

	log.Println("server gracefully shutdown")
}
