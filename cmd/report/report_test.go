package report

import (
	"encoding/json"
	"errors"
	"github.com/metal-stack/go-hal"
	"github.com/metal-stack/go-hal/pkg/api"
	"io/ioutil"
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

	metalCoreURL := ts.Listener.Addr().String()

	transport := httptransport.New(metalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)

	r := &Report{
		Client:       client,
		InstallError: errors.New("an error occurred"),
		Hal:          &inBand{},
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

type inBand struct {
	hal.InBand
}

func (ib *inBand) PowerOff() error {
	return nil
}

func (ib *inBand) PowerCycle() error {
	return nil
}

func (ib *inBand) PowerReset() error {
	return nil
}

func (ib *inBand) IdentifyLEDState(state hal.IdentifyLEDState) error {
	return nil
}

func (ib *inBand) IdentifyLEDOn() error {
	return nil
}

func (ib *inBand) IdentifyLEDOff() error {
	return nil
}

func (ib *inBand) BootFrom(bootTarget hal.BootTarget) error {
	return nil
}

func (ib *inBand) SetFirmware(hal.FirmwareMode) error {
	return nil
}

func (ib *inBand) Describe() string {
	return "InBand mock"
}

func (ib *inBand) BMC() (*api.BMC, error) {
	return nil, nil
}

func (ib *inBand) BMCPresentSuperUser() hal.BMCUser {
	return hal.BMCUser{}
}

func (ib *inBand) BMCSuperUser() hal.BMCUser {
	return hal.BMCUser{}
}

func (ib *inBand) BMCUser() hal.BMCUser {
	return hal.BMCUser{}
}

func (ib *inBand) BMCPresent() bool {
	return false
}

func (ib *inBand) BMCCreateUserAndPassword(user hal.BMCUser, privilege api.IpmiPrivilege, constraints api.PasswordConstraints) (string, error) {
	return "", nil
}

func (ib *inBand) BMCCreateUser(user hal.BMCUser, privilege api.IpmiPrivilege, password string) error {
	return nil
}

func (ib *inBand) BMCChangePassword(user hal.BMCUser, newPassword string) error {
	return nil
}

func (ib *inBand) BMCSetUserEnabled(user hal.BMCUser, enabled bool) error {
	return nil
}

func (ib *inBand) ConfigureBIOS() (bool, error) {
	return false, nil
}

func (ib *inBand) EnsureBootOrder(bootloaderID string) error {
	return nil
}
