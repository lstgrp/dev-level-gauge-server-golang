package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func DataHandler(w http.ResponseWriter, r *http.Request) {
	var data LevelGaugeData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Error while parsing JSON data: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := data.Validate(); err != nil {
		log.Printf("Error while validating JSON data: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		ExecuteAllHandlers(data)
	}()

	log.Printf("Received Data: %v\n", data)
}
