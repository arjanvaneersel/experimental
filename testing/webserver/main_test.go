package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDoubleHandler(t *testing.T) {
	tt := []struct {
		name   string
		value  string
		double int
		err    string
	}{
		{name: "double of two", value: "2", double: 4},
		{name: "missing value", value: "", err: "missing value"},
		{name: "invald value", value: "a", err: "not a number: a"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "localhost:8000/double?v="+tc.value, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			w := httptest.NewRecorder()
			doubleHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.err != "" {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status %v, but got status %v", http.StatusBadRequest, res.StatusCode)
				}

				if msg := string(bytes.TrimSpace(b)); msg != tc.err {
					t.Errorf("expected message %q, but got %q", tc.err, msg)
				}
				return
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK, but got status %v", res.StatusCode)
			}

			d, err := strconv.Atoi(string(b))
			if err != nil {
				t.Fatalf("could not convert response %q to int: %v", b, err)
			}

			if d != tc.double {
				t.Fatalf("expected d to be %d, but got %d", tc.double, d)
			}
		})
	}
}

func TestRouting(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()
	res, err := http.Get(fmt.Sprintf("%s/double?v=2", srv.URL))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, but got status %v", res.StatusCode)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	d, err := strconv.Atoi(string(b))
	if err != nil {
		t.Fatalf("could not convert response %q to int: %v", b, err)
	}

	if d != 4 {
		t.Fatalf("expected d to be 4, but got %d", d)
	}

}
