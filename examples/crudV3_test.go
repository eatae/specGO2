package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowArticles(t *testing.T) {
	// запрос к нашему API
	request, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowArticles)
	handler.ServeHTTP(recorder, request)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
