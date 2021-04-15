package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
)

func TestRegisterMachine(t *testing.T) {
	// FIXME test is disabled
	t.Skip()
	os.Setenv("DEGUG", "1")
	expected := "1234-1234"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metalMachine := &models.ModelsV1MachineResponse{
			ID: &expected,
		}
		response, err := json.Marshal(metalMachine)
		if err != nil {
			fmt.Fprint(w, err)
		}
		fmt.Fprint(w, string(response))
	}))
	defer ts.Close()
	metalCoreURL := ts.Listener.Addr().String()
	transport := httptransport.New(metalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)

	interfaces := make([]string, 0)
	lldpc := network.NewLLDPClient(interfaces, 0, 0, 2*time.Second)
	go lldpc.Start()
	n := &network.Network{
		LLDPClient: lldpc,
		Eth0Mac:    "00:00:00:00:00:01",
	}
	r := &Register{
		Client:      client,
		Network:     n,
		MachineUUID: expected,
	}

	hw, err := r.ReadHardwareDetails()
	if err != nil {
		t.Error(err)
	}

	err = r.RegisterMachine(hw)

	if err != nil {
		t.Error(err)
	}

	if hw.UUID != expected {
		t.Errorf("did not get %s, got %#v ", expected, hw.UUID)
	}
}

func Test_readHardwareDetails(t *testing.T) {
	// FIXME test is disabled
	t.Skip()
	type fields struct {
		Client *machine.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *models.DomainMetalHammerRegisterMachineRequest
		wantErr bool
	}{
		{
			name: "simple",
			want: &models.DomainMetalHammerRegisterMachineRequest{
				UUID: "00000000-0000-0000-0000-000000000000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := &Register{
				Client:      tt.fields.Client,
				MachineUUID: "00000000-0000-0000-0000-000000000000",
			}
			got, err := r.ReadHardwareDetails()
			if (err != nil) != tt.wantErr {
				t.Errorf("readHardwareDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.UUID) != len(tt.want.UUID) {
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

func TestHammer_readIPMIDetails(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.ModelsV1MachineIPMI
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := readIPMIDetails("00:00:00:00:00:01", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hammer.readIPMIDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hammer.readIPMIDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDCreation(t *testing.T) {
	uuidAsString, err := uuid.FromBytes([]byte("S167357X6205283" + " "))
	if err != nil {
		t.Error(err)
	}
	t.Logf("got: %s", uuidAsString)

	uuidAsString2, err := uuid.FromBytes([]byte("S167357X6205283" + " "))
	if err != nil {
		t.Error(err)
	}
	if uuidAsString != uuidAsString2 {
		t.Errorf("expected same uuid, got different: %s vs: %s", uuidAsString, uuidAsString2)
	}
}
