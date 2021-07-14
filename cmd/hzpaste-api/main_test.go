package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCNotEnoughDataToCollect(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/gc", nil)
	router.ServeHTTP(w, req)

	ErrNotEnoughDataForGC := `{"error":"Not enough data for GC"}`

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, ErrNotEnoughDataForGC, w.Body.String())
}

func TestNoPastesFound(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pastes", nil)
	router.ServeHTTP(w, req)

	ErrNotFound := `{"error":"Paste not found"}`

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, ErrNotFound, w.Body.String())
}

func TestInvalidPasteUUID(t *testing.T) {

	router := setupRouter()

	wrongUUID := "I am wrongUUID"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pastes/"+wrongUUID, nil)
	router.ServeHTTP(w, req)

	ErrInvalidID := `{"error":"Paste UUID is invalid"}`

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, ErrInvalidID, w.Body.String())
}

func TestPasteUUIDNotFound(t *testing.T) {

	router := setupRouter()

	randomUUID := "ee6edb3a-db64-11eb-8d19-0242ac130003" // random UUID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pastes/"+randomUUID, nil)
	router.ServeHTTP(w, req)

	ErrInvalidID := `{"error":"Paste not found"}`

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, ErrInvalidID, w.Body.String())
}

func TestPasteAdd(t *testing.T) {

	router := setupRouter()

	newPaste := `{"name":"Very important paste",
	              "body":"I can see the world!"}`
	post := strings.NewReader(newPaste)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pastes", post)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Decode newPaste string to struct
	type Paste struct {
		Name, Body string
	}
	dec := json.NewDecoder(strings.NewReader(newPaste))
	var p Paste
	err := dec.Decode(&p)
	if err != nil {
		t.Fatal(err)
	}

	// Unpacking the map[string]interface{} to map[string]string
	var created map[string]interface{}
	createdString := make(map[string]string)
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decoding: %s", err)
	}
	for key, value := range created["new_paste"].(map[string]interface{}) {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		createdString[strKey] = strValue
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, p.Name, createdString["name"])
	assert.Equal(t, p.Body, createdString["body"])
}

func TestKillPastes(t *testing.T) {

	router := setupRouter()

	newPaste := `{"name":"Very important paste",
	              "body":"I can see the world!"}`
	for i := 0; i < 2; i++ {
		post := strings.NewReader(newPaste)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/pastes", post)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
	req, _ := http.NewRequest("DELETE", "/killall", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	OkClearPastes := `{"pastes":"null"}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, OkClearPastes, w.Body.String())
}

func TestPing(t *testing.T) {

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	Ready := `{"ready":true}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, Ready, w.Body.String())
}
