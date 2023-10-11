package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestServer(db *sql.DB) *httptest.Server {
	server := App(db)
	return httptest.NewServer(server.Handler)
}

func TestApp_Integration(t *testing.T) {
	dbUrl := os.Getenv("TEST_DATABASE_URL")

	if dbUrl == "" {
		t.Fatalf("TEST_DATABASE_URL is required to run integration tests")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ts := newTestServer(db)
	defer ts.Close()

	t.Run("cell not found", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/sheet_not_found", ts.URL))
		if err != nil {
			t.Fatalf("expected no error, got (%v)", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want (%v) got (%v)", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("sheet not found", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/sheet_not_found/cell_not_found", ts.URL))
		if err != nil {
			t.Fatalf("expected no error, got (%v)", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want (%v) got (%v)", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("create cell", func(t *testing.T) {
		cellID := "cell_create_cell"
		sheetID := "sheet_create_cell"

		value := "2+2"
		result := "4"

		body := bytes.NewBufferString(fmt.Sprintf("{\"value\": \"%s\"}", value))

		resp, err := http.Post(fmt.Sprintf("%s/api/v1/%s/%s", ts.URL, sheetID, cellID), "application/json", body)
		if err != nil {
			t.Fatalf("expected no error, got (%v)", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("want (%v) got (%v)", http.StatusNotFound, resp.StatusCode)
		}

		respBody := struct {
			Result string `json:"result"`
			Value  string `json:"value"`
		}{}

		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			t.Fatalf("could not decode a response body: %v", err)
		}

		if respBody.Value != value {
			t.Fatalf("want (%s) got (%v)", value, respBody.Result)
		}

		if respBody.Result != result {
			t.Fatalf("want (%s) got (%v)", result, respBody.Result)
		}

		t.Run("get cell", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/api/v1/%s/%s", ts.URL, sheetID, cellID))
			if err != nil {
				t.Fatalf("expected no error, got (%v)", err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("want (%v) got (%v)", http.StatusOK, resp.StatusCode)
			}

			respBody := struct {
				Result string `json:"result"`
				Value  string `json:"value"`
			}{}

			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatalf("could not decode a response body: %v", err)
			}

			if respBody.Value != value {
				t.Fatalf("want (%s) got (%v)", value, respBody.Result)
			}

			if respBody.Result != result {
				t.Fatalf("want (%s) got (%v)", result, respBody.Result)
			}
		})

		t.Run("using exising cell in formula", func(t *testing.T) {
			currentCellID := "cell_using_existing_cell_in_formula"
			currentValue := fmt.Sprintf("%s+(%s-1)", cellID, cellID)
			currentResult := "7"

			body := bytes.NewBufferString(fmt.Sprintf("{\"value\": \"%s\"}", currentValue))

			resp, err := http.Post(fmt.Sprintf("%s/api/v1/%s/%s", ts.URL, sheetID, currentCellID), "application/json", body)
			if err != nil {
				t.Fatalf("expected no error, got (%v)", err)
			}

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("want (%v) got (%v)", http.StatusNotFound, resp.StatusCode)
			}

			respBody := struct {
				Result string `json:"result"`
				Value  string `json:"value"`
			}{}

			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatalf("could not decode a response body: %v", err)
			}

			if respBody.Value != currentValue {
				t.Fatalf("want (%s) got (%v)", value, respBody.Result)
			}

			if respBody.Result != currentResult {
				t.Fatalf("want (%s) got (%v)", result, respBody.Result)
			}

		})
	})

	t.Run("incorrect formula", func(t *testing.T) {
		t.Run("invalid parentheses", func(t *testing.T) {

			cellID := "cell_invalid_parentheses"
			sheetID := "sheet_invalid_parentheses"

			value := "2+((-4"
			result := "ERROR"
			message := "invalid parentheses"

			body := bytes.NewBufferString(fmt.Sprintf("{\"value\": \"%s\"}", value))

			resp, err := http.Post(fmt.Sprintf("%s/api/v1/%s/%s", ts.URL, sheetID, cellID), "application/json", body)
			if err != nil {
				t.Fatalf("expected no error, got (%v)", err)
			}

			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("want (%v) got (%v)", http.StatusUnprocessableEntity, resp.StatusCode)
			}

			respBody := struct {
				Result  string `json:"result"`
				Value   string `json:"value"`
				Message string `json:"message"`
			}{}

			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatalf("could not decode a response body: %v", err)
			}

			if respBody.Value != value {
				t.Fatalf("want (%s) got (%v)", value, respBody.Result)
			}

			if respBody.Result != result {
				t.Fatalf("want (%s) got (%v)", result, respBody.Result)
			}

			if respBody.Message != message {
				t.Fatalf("want (%s) got (%v)", message, respBody.Message)
			}
		})

		t.Run("invalid operation", func(t *testing.T) {

			cellID := "cell_invalin_operation"
			sheetID := "sheet_invalid_operation"

			value := "*(3+3)"
			result := "ERROR"
			message := "invalid operation"

			body := bytes.NewBufferString(fmt.Sprintf("{\"value\": \"%s\"}", value))

			resp, err := http.Post(fmt.Sprintf("%s/api/v1/%s/%s", ts.URL, sheetID, cellID), "application/json", body)
			if err != nil {
				t.Fatalf("expected no error, got (%v)", err)
			}

			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("want (%v) got (%v)", http.StatusUnprocessableEntity, resp.StatusCode)
			}

			respBody := struct {
				Result  string `json:"result"`
				Value   string `json:"value"`
				Message string `json:"message"`
			}{}

			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatalf("could not decode a response body: %v", err)
			}

			if respBody.Value != value {
				t.Fatalf("want (%s) got (%v)", value, respBody.Result)
			}

			if respBody.Result != result {
				t.Fatalf("want (%s) got (%v)", result, respBody.Result)
			}

			if respBody.Message != message {
				t.Fatalf("want (%s) got (%v)", message, respBody.Message)
			}
		})
	})
}
