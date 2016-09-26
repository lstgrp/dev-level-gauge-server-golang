package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var server *Server

func init() {
	server = InitServer()
	server.Redis.Flush()
}

func TestStoreDataSuccess(t *testing.T) {
	testData := LevelGaugeData{
		DeviceId: "test_id",
		Time:     time.Now().Unix(),
		Event:    0,
		Level:    1,
	}

	body, _ := json.Marshal(testData)

	req, _ := http.NewRequest("POST", "/store", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Correct json body (should succeed), instead received &v", w)
	}
}

func TestStoreDataIncompleteBody(t *testing.T) {
	testData := LevelGaugeData{
		DeviceId: "test_id",
		Event:    0,
		Level:    1,
	}

	body, _ := json.Marshal(testData)

	req, _ := http.NewRequest("POST", "/store", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("Incomplete JSON body should make request fail")
	}
}

func TestGenerateTokenSuccess(t *testing.T) {
	testData := TokenParameter{
		Device: struct {
			Name   string `json:"name"`
			Serial string `json:"serial"`
		}{
			Name:   "test_name",
			Serial: "test_serial",
		},
	}

	body, _ := json.Marshal(testData)

	req, _ := http.NewRequest("POST", "/device", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	paramBody := struct {
		DeviceId string
		Token    string
		TTL      int64
	}{}
	json.Unmarshal(w.Body.Bytes(), &paramBody)

	res, _ := server.Redis.Do("get", paramBody.Token)

	if res == nil {
		t.Error("Token should be stored in redis")
	}
}
