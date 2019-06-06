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

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/machine"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
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
	}
	r := &Register{
		Client:      client,
		Network:     n,
		MachineUUID: expected,
	}

	eth0Mac = "00:00:00:00:00:01"
	_, uuid, err := r.RegisterMachine()

	if err != nil {
		t.Error(err)
	}

	if uuid != expected {
		t.Errorf("did not get %s, got %#v ", expected, uuid)
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
		t.Run(tt.name, func(t *testing.T) {
			r := &Register{
				Client:      tt.fields.Client,
				MachineUUID: "00000000-0000-0000-0000-000000000000",
			}
			eth0Mac = "00:00:00:00:00:01"
			got, err := r.readHardwareDetails()
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
	type fields struct {
		Client *machine.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *models.ModelsV1MachineIPMI
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readIPMIDetails("00:00:00:00:00:01")
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
