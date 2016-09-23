package main

import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/data", DataHandler)
	log.Print("Server will listen on port 5656")
	log.Fatal(http.ListenAndServe(":5656", nil))
}
