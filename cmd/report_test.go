package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportInstallation(t *testing.T) {
	expected := "an error occured"
	resp := &Report{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(body, resp)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	err := ReportInstallation(ts.URL, expected, errors.New("an error occured"))
	if err != nil {
		t.Error(err)
	}

	if resp.Message != expected {
		t.Errorf("response message:%s expected:%s", resp.Message, expected)
	}
	if resp.Success {
		t.Errorf("response success:%t expected:False", resp.Success)
	}
}
