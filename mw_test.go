package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mw_server *Server

func init() {
	mw_server = InitServer(true)
}

func TestValidateTokenSuccess(t *testing.T) {
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
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mw_server.Router.ServeHTTP(w, req)

	paramBody := struct {
		DeviceId string
		Token    string
		TTL      int64
	}{}
	json.Unmarshal(w.Body.Bytes(), &paramBody)

	if w.Code != http.StatusOK || paramBody.Token == "" {
		t.Errorf("Requets for token failed, got token: %v", paramBody.Token)
	}

	testTokenData := struct {
		Token string `json:"token"`
	}{
		Token: paramBody.Token,
	}

	body, _ = json.Marshal(testTokenData)
	req, _ = http.NewRequest("POST", "/close", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-jwt", paramBody.Token)
	w = httptest.NewRecorder()
	mw_server.Router.ServeHTTP(w, req)

	finalBody := struct {
		Result string `json:"result"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &finalBody)

	if finalBody.Result != "ok" {
		t.Error("Session was not closed successfuly")
	}

	res, _ := mw_server.Redis.Do("get", paramBody.Token)

	if res != nil {
		t.Error("Session token was not erased successfully")
	}
}

func TestValidateTokenInvalidToken(t *testing.T) {
	testTokenData := struct {
		Token string `json:"token"`
	}{
		Token: "wrong_token",
	}

	body, _ := json.Marshal(testTokenData)
	req, _ := http.NewRequest("POST", "/close", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-jwt", testTokenData.Token)
	w := httptest.NewRecorder()
	mw_server.Router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("Request with wrong token should be rejected")
	}
}

func TestValidateTokenMasterKey(t *testing.T) {
	testTokenData := struct {
		Token string `json:"token"`
	}{
		Token: "wrong_token",
	}

	body, _ := json.Marshal(testTokenData)
	req, _ := http.NewRequest("POST", "/close", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-master-key", LocalConfig.masterKey)
	w := httptest.NewRecorder()
	mw_server.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Request with master key in header should always be authorized")
	}
}
