package register

import (
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

// Register the Machine
type Register struct {
	MachineUUID string
	Client      *machine.Client
	Network     *network.Network
}

// RegisterMachine register a machine at the metal-api via metal-core
func (r *Register) RegisterMachine(hw *models.DomainMetalHammerRegisterMachineRequest) error {
	params := machine.NewRegisterParams()
	params.SetBody(hw)
	params.ID = hw.UUID
	resp, err := r.Client.Register(params)

	if err != nil {
		return errors.Wrapf(err, "unable to register machine:%#v", hw)
	}
	if resp == nil {
		return errors.Errorf("unable to register machine:%#v response payload is nil", hw)
	}

	log.Info("register machine returned", "response", resp.Payload)
	return nil
}
