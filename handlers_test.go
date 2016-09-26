package main

import (
  "testing"
  "net/http"
  "time"
  "encoding/json"
  "bytes"
  "net/http/httptest"
)

var server *Server

func init() {
  server = InitServer()
}

func TestStoreDataSuccess(t *testing.T) {
  testData := LevelGaugeData{
    DeviceId: "test_id",
    Time:     time.Now().Unix(),
    Event: 0,
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
    Event: 0,
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
