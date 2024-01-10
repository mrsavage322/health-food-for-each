package main

import (
	_ "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		request                string
		body                   string
		expectedStatusCode     int
		expectedLocationHeader string
		expectedResponseBody   string
	}{
		{
			name:               "POST request with a valid link",
			method:             http.MethodPost,
			request:            "/",
			body:               "https://example.com",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "POST request with an empty link",
			method:             http.MethodPost,
			request:            "/",
			body:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:                 "JSON POST request with a valid link",
			method:               http.MethodPost,
			request:              "/api/shorten",
			body:                 `{"url": "https://example.com"}`,
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"result": "YourShortURLLogicHere"}`,
		},
		{
			name:                 "JSON POST request with an empty link",
			method:               http.MethodPost,
			request:              "/api/shorten",
			body:                 ``,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "",
		},
	}
