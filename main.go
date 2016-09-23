package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/data", DataHandler)

	log.Print("Server will listen on port 5656")
	log.Fatal(http.ListenAndServe(":5656", nil))
}
