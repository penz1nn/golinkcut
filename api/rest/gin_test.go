package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cfg = config.Config{
	"memory":   true,
	"debug":    false,
	"validate": true,
	"httpHost": "localhost",
	"httpPort": "8080",
	"redirect": false,
}

func setupRouter() *gin.Engine {
	s := link.NewDbStorage(cfg)
	uc := link.NewUseCase(s, cfg)
	return SetupRouter(uc, cfg)
}

func TestRestApi_CreateLink(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	values := map[string]string{"Url": "ozon.ru/JosdnnHUSIDFmfklsnJSKLf9234fdy2"}
	json_data, err := json.Marshal(values)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	req, err := http.NewRequest("POST", "/new", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Wrong response status: %v", w.Code)
	}
}

func TestRestApi_GetLink(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	values := map[string]string{"Url": "wildberries.ru/sale/?p=jhus087hsdkfJhdsf6gdshfj"}
	json_data, err := json.Marshal(values)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	req, err := http.NewRequest("POST", "/new", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Wrong response status code: %v", w.Code)
	}
	var resp = map[string]string{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	shortLink := resp["shortLink"]
	alias := shortLink[len(shortLink)-10 : len(shortLink)]
	if len(alias) != 10 {
		t.Errorf("Error getting alias: %s", alias)
	}
	req, err = http.NewRequest("GET", "/"+alias, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Wrong response status code: %v", w.Code)
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	got := resp["url"]
	want := "wildberries.ru/sale/?p=jhus087hsdkfJhdsf6gdshfj"
	if got != want {
		t.Errorf("Got wrong link.\ngot:\t%v\nwant\t%v\nfull response: %v\nalias: %v", got, want, w.Body.String(), alias)
	}
}

func TestRestApi_CreateLink_Exists(t *testing.T) {
	url := "https://pkg.go.dev/net/url#URL"
	r := setupRouter()
	w := httptest.NewRecorder()
	values := map[string]string{"Url": url}
	json_data, err := json.Marshal(values)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	req, err := http.NewRequest("POST", "/new", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Wrong response status: %v", w.Code)
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/new", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 409 {
		t.Errorf("Wrong response status: %v", w.Code)
	}
}

func TestRestApi_CreateLink_Bad(t *testing.T) {
	url := "not-really-a-url"
	r := setupRouter()
	w := httptest.NewRecorder()
	values := map[string]string{"Url": url}
	json_data, err := json.Marshal(values)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	req, err := http.NewRequest("POST", "/new", bytes.NewBuffer(json_data))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("Wrong response status: %v", w.Code)
	}
}

func TestRestApi_GetLink_NotExist(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/0123456789", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Wrong response status code: %v", w.Code)
	}
}
