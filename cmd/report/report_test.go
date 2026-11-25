package report

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
	"google.golang.org/protobuf/testing/protocmp"
)

type bootClient struct {
	infrav2connect.BootServiceClient

	requests []*infrav2.BootServiceInstallationSucceededRequest
	err      error
}

func (b *bootClient) InstallationSucceeded(_ context.Context, req *infrav2.BootServiceInstallationSucceededRequest) (*infrav2.BootServiceInstallationSucceededResponse, error) {
	b.requests = append(b.requests, req)
	return &infrav2.BootServiceInstallationSucceededResponse{}, b.err
}

func TestReport_ReportInstallation(t *testing.T) {
	tests := []struct {
		name         string
		installError error
		clientError  error
		want         []*infrav2.BootServiceInstallationSucceededRequest
		wantErr      bool
	}{
		{
			name: "successful installation is reported",
			want: []*infrav2.BootServiceInstallationSucceededRequest{
				{Uuid: "machine-uuid", ConsolePassword: "password"},
			},
		},
		{
			name:         "failed installation is not reported as succeeded",
			installError: errors.New("unable to umount"),
			want:         nil,
		},
		{
			name:        "error from the api is returned",
			clientError: errors.New("unavailable"),
			want: []*infrav2.BootServiceInstallationSucceededRequest{
				{Uuid: "machine-uuid", ConsolePassword: "password"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &bootClient{err: tt.clientError}
			r := &Report{
				MachineUUID:     "machine-uuid",
				ConsolePassword: "password",
				InstallError:    tt.installError,
				Log:             slog.Default(),
			}

			err := r.ReportInstallation(client)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReportInstallation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, client.requests, protocmp.Transform()); diff != "" {
				t.Errorf("InstallationSucceeded() requests diff = %s", diff)
			}
		})
	}
}
