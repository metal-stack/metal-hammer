package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func TestReportInstallation(t *testing.T) {
	expected := "an error occured"
	resp := &models.DomainReport{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(body, resp)
		if err != nil {
			t.Error(err)
		}
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	spec := &Specification{
		MetalCoreURL: ts.Listener.Addr().String(),
	}

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := device.New(transport, strfmt.Default)

	h := &Hammer{
		Client: client,
		Spec:   spec,
	}

	err := h.ReportInstallation(expected, errors.New("an error occured"))
	if err != nil {
		t.Error(err)
	}

	if *resp.Message != expected {
		t.Errorf("response message:%s expected:%s", *resp.Message, expected)
	}
	if resp.Success {
		t.Errorf("response success:%t expected:False", resp.Success)
	}
}
