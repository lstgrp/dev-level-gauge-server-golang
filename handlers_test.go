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
	server.Redis.Do("flushall")

	for i := 1; i < 11; i++ {
		data := LevelGaugeData{
			DeviceId: "test_device",
			Time:     now + int64(i*10),
			Event:    uint8(i % 2),
			Level:    uint8(i),
		}

		dataJSONStr, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", "/store", bytes.NewReader(dataJSONStr))
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
	}
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

func TestRetrieveDataSuccess(t *testing.T) {
	query := LevelGaugeDataQuery{
		DeviceId: "test_device",
		Date:     []int64{0, -1},
		Event:    -1,
	}

	body, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", "/retrieve", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)
	var res struct {
		Result string `json:"result"`
	}
	json.Unmarshal(w.Body.Bytes(), &res)
	var finalData []LevelGaugeData
	json.Unmarshal([]byte(res.Result), &finalData)

	if length := len(finalData); length != 10 {
		t.Errorf("Should have received all data, instead got length: %v", length)
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
