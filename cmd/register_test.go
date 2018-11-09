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

func Test_readHardwareDetails(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.DomainMetalHammerRegisterDeviceRequest
		wantErr bool
	}{
		{
			name: "simple",
			want: &models.DomainMetalHammerRegisterDeviceRequest{
				UUID: "00000000-0000-0000-0000-000000000000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readHardwareDetails()
			if (err != nil) != tt.wantErr {
				t.Errorf("readHardwareDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UUID != tt.want.UUID {
				t.Errorf("readHardwareDetails() expected uuid: %s got %s", tt.want.UUID, got.UUID)
			}
			if *got.CPUCores == 0 {
				t.Errorf("readHardwareDetails() expected cpucores: %d", got.CPUCores)
			}
			if *got.Memory == 0 {
				t.Errorf("readHardwareDetails() expected memory: %d", got.Memory)
			}
		})
	}
}
