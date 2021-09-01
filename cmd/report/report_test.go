package report

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func TestReportInstallation(t *testing.T) {
	expected := "an error occurred"
	resp := &models.DomainReport{}

	handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
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

	metalCoreURL := ts.Listener.Addr().String()

	transport := httptransport.New(metalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)

	r := &Report{
		Client:       client,
		InstallError: errors.New("an error occurred"),
	}

	err := r.ReportInstallation()
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
