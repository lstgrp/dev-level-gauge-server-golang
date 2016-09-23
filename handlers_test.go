package main

import (
  "time"
  "testing"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "bytes"
)

func TestDataHandlerSuccess (t *testing.T) {
  testData := LevelGaugeData{
    DeviceId: "test_id",
    Time: time.Now().String(),
    Level: 1,
  }

  body, _ := json.Marshal(testData)

  req, _ := http.NewRequest("POST", "/data", bytes.NewReader(body))
  rr := httptest.NewRecorder()
  handler := http.HandlerFunc(DataHandler)

  handler.ServeHTTP(rr, req);

  if status := rr.Code; status != http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }
}

func TestDataHandlerBadData (t *testing.T) {
  req, _ := http.NewRequest("POST", "/data", bytes.NewReader([]byte("")))
  rr := httptest.NewRecorder()
  handler := http.HandlerFunc(DataHandler)

  handler.ServeHTTP(rr, req);

  if status := rr.Code; status == http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusBadRequest)
  }
}
