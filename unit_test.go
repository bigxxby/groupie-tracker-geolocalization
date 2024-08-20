package main

import (
	handlers "bigxxby/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.MainHandler)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusOK {
		t.Errorf("Main handler returned wrong status code. Got %v, want %v", status, http.StatusOK)
	}
}

func TestMainHandler_NotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.MainHandler)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusNotFound {
		t.Errorf("Home handler returned wrong status code. Got %v, want %v", status, http.StatusNotFound)
	}
}

func TestMainHandler_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.MainHandler)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Home handler returned wrong status code. Got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestArtistPageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ArtistIdHandler)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusOK {
		t.Errorf("ArtistPage handler returned wrong status code. Got %v, want %v", status, http.StatusOK)
	}
}

func TestArtistPageHandler_InvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists/53", nil)
	if err != nil {
		t.Fatal(err)
	}

	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ArtistIdHandler)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusNotFound {
		t.Errorf("ArtistPage handler returned wrong status code. Got %v, want %v", status, http.StatusNotFound)
	}
}
