package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/models"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func TestRegisterDevice(t *testing.T) {
	os.Setenv("DEGUG", "1")
	expected := "1234-1234"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metalDevice := &models.ModelsMetalDevice{
			ID: &expected,
		}
		response, err := json.Marshal(metalDevice)
		if err != nil {
			fmt.Fprint(w, err)
		}
		fmt.Fprint(w, string(response))
	}))
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

	uuid, err := h.RegisterDevice()

	if err != nil {
		t.Error(err)
	}

	if uuid != expected {
		t.Errorf("did not get %s, got %#v ", expected, uuid)
	}
}
