package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStoreData(t *testing.T) {
	server := InitServer(false)
	defer server.Teardown()
	server.Redis.Do("flushall")

	t.Run("Success", func(t *testing.T) {
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
	})

	t.Run("Fail, incomplete JSON body", func(t *testing.T) {
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
	})
}

func TestRetrieveData(t *testing.T) {
	server := InitServer(false)
	defer server.Teardown()
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

	t.Run("Success", func(t *testing.T) {
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
	})
}

func TestGenerateToken(t *testing.T) {
	server := InitServer(false)
	defer server.Teardown()
	server.Redis.Do("flushall")

	t.Run("Success", func(t *testing.T) {
		testData := TokenParameter{
			Device: struct {
				Name   string `json:"name" binding:"required"`
				Serial string `json:"serial" binding:"required"`
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
	})
}

func TestCloseSession(t *testing.T) {
	server := InitServer(false)
	defer server.Teardown()
	server.Redis.Do("flushall")

	t.Run("Success", func(t *testing.T) {
		testData := TokenParameter{
			Device: struct {
				Name   string `json:"name" binding:"required"`
				Serial string `json:"serial" binding:"required"`
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

		testTokenData := struct {
			Token string `json:"token"`
		}{
			Token: paramBody.Token,
		}

		body, _ = json.Marshal(testTokenData)
		req, _ = http.NewRequest("POST", "/close", bytes.NewReader(body))
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)
		finalBody := struct {
			Result string `json:"result"`
		}{}
		json.Unmarshal(w.Body.Bytes(), &finalBody)

		if finalBody.Result != "ok" {
			t.Error("Session was not closed successfuly")
		}

		res, _ := server.Redis.Do("get", paramBody.Token)

		if res != nil {
			t.Error("Session token was not erased successfully")
		}
	})
}
