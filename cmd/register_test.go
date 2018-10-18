package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterDevice(t *testing.T) {
	expected := "1234-1234"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("{\"id\": \"%s\"}", expected))
	}))
	defer ts.Close()

	spec := &Specification{
		RegisterURL: ts.URL,
	}
	uuid, err := RegisterDevice(spec)

	if err != nil {
		t.Error(err)
	}

	if uuid != expected {
		t.Errorf("did not get %s, got %s ", uuid, expected)
	}
}
