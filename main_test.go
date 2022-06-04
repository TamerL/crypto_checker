package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var tests = []struct {
	name            string
	endpoint        string
	expected_header string
	expected_body   string
	isErr           bool
}{
	{"valid-data", "/usd", "application/json", "", false},
	{"valid-data", "/EUR", "application/json", "", false},
	{"valid-data", "/health", "application/json", `{"alive": true}`, false},
	{"invalid-data", "/AUD", "application/json", "", true},
	{"invalid-data", "/", "application/json", "", true},
	{"invalid-data", "/asdfasdf", "application/json", "", true},
}

func TestHealthHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	for _, tt := range tests {
		req, err := http.NewRequest("GET", tt.endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		test_handler := http.HandlerFunc(Handler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		test_handler.ServeHTTP(rr, req)

		// Testing valid-data scenarios
		if !tt.isErr {
			// Check the status code is 200 OK.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			// Check the response content-type header is json.
			if ctype := rr.Header().Get("Content-Type"); ctype != tt.expected_header {
				t.Errorf("content type header does not match: got %v want %v",
					ctype, tt.expected_header)
			}
			if tt.endpoint == "/health" {
				if rr.Body.String() != tt.expected_body {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tt.expected_body)
				}
			}
		} else {
			// Testing invalid-data scenarios
			if rr.Code == http.StatusOK {
				t.Error("expected error but didn't get any")
			} else if tt.endpoint == "/AUD" && rr.Code != http.StatusMethodNotAllowed {
				t.Error("expected but not received error: method not allowed")
			} else if tt.endpoint == "/" && rr.Code != http.StatusNotFound {
				t.Error("expected but not received error: not found")
			}
		}
	}
}
