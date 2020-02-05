package main

import (
	"log"
	"net/http"
	"os"
)

type task struct {
	message
}

type Face uint8

const (
	One   Face = iota
	Two   Face = iota
	Three Face = iota
	Four  Face = iota
	Five  Face = iota
	Six   Face = iota
)

func main() {
	s := &server{
		token: apiToken,
	}

	http.HandleFunc("/", s.apiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
