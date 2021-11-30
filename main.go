package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var file string

func main() {

	file = os.Getenv("LARGE_FILE")
	// Create a new ServerMux
	mux := http.NewServeMux()
	// Add home HandlerFunc to handle the root path
	mux.HandleFunc("/", home)
	mux.HandleFunc("/stress", stress)

	// Create the Server srtuct abd initialized with
	// localhost on port 5000 and with the mux created
	srv := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	// idleClose channel signals when shutdown is completed and
	// connections are closed
	idleClose := make(chan interface{})

	// OS signal notification and server shutdown is running
	// in a separate goroutine
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		// sigint halts until OS signal received
		<-sigint

		log.Println("HTTP Server is shutting down....")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown error: %v", err)
		}
		// when shutdown is completed and connections closed,
		// close the idleClose channel and main goroutine will continue
		close(idleClose)
	}()

	// start server and waiting for incoming connections
	log.Println("HTTP Server listening...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP Server error: %v", err)
	}
	// main goroutine halts until idleClose channel is closed
	<-idleClose
	log.Println("HTTP Server is down")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome memory stress!")
}

func stress(w http.ResponseWriter, r *http.Request) {
	//log.Println(file)
	dat, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error reading file: %s", err)
		fmt.Fprint(w, "NOK")
		return
	}
	l := len(dat)
	log.Printf("the file is %d bytes", l)
	fmt.Fprint(w, "OK")
}
