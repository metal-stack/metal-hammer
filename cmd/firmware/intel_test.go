package firmware

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseNVMVersion(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    string
		wantErr bool
	}{
		{
			name: "e810 port",
			output: `driver: ice
version: 6.8.0-generic
firmware-version: 4.60 0x8001e8e4 1.3653.0
expansion-rom-version:
bus-info: 0000:51:00.1
supports-statistics: yes
`,
			want: "4.60",
		},
		{
			name: "firmware version without additional fields",
			output: `driver: ice
firmware-version: 4.80
bus-info: 0000:51:00.0
`,
			want: "4.80",
		},
		{
			name: "no firmware version reported",
			output: `driver: virtio_net
version: 1.0.0
firmware-version:
bus-info: 0000:00:03.0
`,
			wantErr: true,
		},
		{
			name:    "empty output",
			output:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNVMVersion("eth0", tt.output)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestUpdateRequired(t *testing.T) {
	tests := []struct {
		name    string
		current string
		desired string
		want    bool
		wantErr bool
	}{
		{name: "current behind desired", current: "4.60", desired: "4.80", want: true},
		{name: "current equals desired", current: "4.80", desired: "4.80", want: false},
		{name: "current ahead of desired", current: "5.00", desired: "4.80", want: false},
		{name: "versions are not compared as strings", current: "4.90", desired: "10.00", want: true},
		{name: "unparsable current version", current: "n/a", desired: "4.80", wantErr: true},
		{name: "unparsable desired version", current: "4.60", desired: "unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, current, desired, err := updateRequired(t.Context(), &fakeUpdater{currentVersion: tt.current, desiredVersion: tt.desired})
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
			if current != tt.current || desired != tt.desired {
				t.Errorf("got versions %q/%q want %q/%q", current, desired, tt.current, tt.desired)
			}
		})
	}
}

// fakeUpdater reports fixed versions and the restart its update returns, ActivationNone
// meaning no component took the update.
type fakeUpdater struct {
	name           string
	currentVersion string
	desiredVersion string
	activation     Activation
}

func (f *fakeUpdater) String() string                                 { return f.name }
func (f *fakeUpdater) applicable(ctx context.Context) bool            { return true }
func (f *fakeUpdater) current(ctx context.Context) (string, error)    { return f.currentVersion, nil }
func (f *fakeUpdater) desired(ctx context.Context) (string, error)    { return f.desiredVersion, nil }
func (f *fakeUpdater) update(ctx context.Context) (Activation, error) { return f.activation, nil }

// TestUpdate covers what a run reports to its caller: the strongest activation of all
// components which were actually flashed, and separately the components which stay behind
// because they are retried on every boot without ever converging.
func TestUpdate(t *testing.T) {
	tests := []struct {
		name           string
		updaters       []updater
		wantActivation Activation
		wantNotFlashed []string
	}{
		{
			name:           "component took the update and needs a power cycle",
			updaters:       []updater{&fakeUpdater{name: "fake", currentVersion: "4.60", desiredVersion: "4.80", activation: ActivationPowerCycle}},
			wantActivation: ActivationPowerCycle,
		},
		{
			name:           "component took the update and needs a warm reboot",
			updaters:       []updater{&fakeUpdater{name: "fake", currentVersion: "4.60", desiredVersion: "4.80", activation: ActivationWarmReboot}},
			wantActivation: ActivationWarmReboot,
		},
		{
			name: "strongest activation of several flashed components wins",
			updaters: []updater{
				&fakeUpdater{name: "warm", currentVersion: "4.60", desiredVersion: "4.80", activation: ActivationWarmReboot},
				&fakeUpdater{name: "cold", currentVersion: "4.60", desiredVersion: "4.80", activation: ActivationPowerCycle},
			},
			wantActivation: ActivationPowerCycle,
		},
		{
			name:           "component did not take the update",
			updaters:       []updater{&fakeUpdater{name: "fake", currentVersion: "4.60", desiredVersion: "4.80", activation: ActivationNone}},
			wantActivation: ActivationNone,
			wantNotFlashed: []string{"fake"},
		},
		{
			name:           "component is already up to date",
			updaters:       []updater{&fakeUpdater{name: "fake", currentVersion: "4.80", desiredVersion: "4.80", activation: ActivationPowerCycle}},
			wantActivation: ActivationNone,
		},
		{
			name:           "component is ahead of the packaged version",
			updaters:       []updater{&fakeUpdater{name: "fake", currentVersion: "5.01", desiredVersion: "4.80", activation: ActivationPowerCycle}},
			wantActivation: ActivationNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Firmware{updaters: tt.updaters, log: slog.Default()}

			got := f.Update(t.Context())
			if got.Activation != tt.wantActivation {
				t.Errorf("activation got %v want %v", got.Activation, tt.wantActivation)
			}
			if !slices.Equal(got.NotFlashed, tt.wantNotFlashed) {
				t.Errorf("not flashed got %v want %v", got.NotFlashed, tt.wantNotFlashed)
			}
		})
	}
}

func TestCurrent(t *testing.T) {
	tests := []struct {
		name    string
		cards   []card
		want    string
		wantErr bool
	}{
		{
			name: "single card",
			cards: []card{
				{interfaceName: "eth0", address: "0000:51:00", nvmVersion: "4.60"},
			},
			want: "4.60",
		},
		{
			name: "two cards on different versions",
			cards: []card{
				{interfaceName: "eth0", address: "0000:51:00", nvmVersion: "4.80"},
				{interfaceName: "eth2", address: "0000:98:00", nvmVersion: "4.60"},
			},
			want: "4.60",
		},
		{
			name: "lowest is not the lexicographically smallest",
			cards: []card{
				{interfaceName: "eth0", address: "0000:51:00", nvmVersion: "10.00"},
				{interfaceName: "eth2", address: "0000:98:00", nvmVersion: "4.80"},
			},
			want: "4.80",
		},
		{
			name:    "no cards at all",
			cards:   nil,
			wantErr: true,
		},
		{
			name: "unparsable version",
			cards: []card{
				{interfaceName: "eth0", address: "0000:51:00", nvmVersion: "garbage"},
			},
			wantErr: true,
		},
		{
			name: "one unparsable version does not hide the remaining cards",
			cards: []card{
				{interfaceName: "eth0", address: "0000:51:00", nvmVersion: "garbage"},
				{interfaceName: "eth2", address: "0000:98:00", nvmVersion: "4.60"},
			},
			want: "4.60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newIntel(slog.Default())
			// the cards are injected, so detect() must not look at the real machine
			i.detected = true
			i.cards = tt.cards

			got, err := i.current(t.Context())
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestPciAddress(t *testing.T) {
	root := t.TempDir()

	tests := []struct {
		name    string
		target  string
		want    string
		wantErr bool
	}{
		{name: "second port of a card", target: "../../../0000:51:00.1", want: "0000:51:00"},
		{name: "no pci function", target: "../../../virtual", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devicePath := filepath.Join(root, tt.name, "device")
			if err := os.MkdirAll(filepath.Dir(devicePath), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(tt.target, devicePath); err != nil {
				t.Fatal(err)
			}

			got, err := pciAddress(devicePath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}

// TestApplicable covers the packaging part of applicable(), the card detection needs an
// ethtool reporting a firmware version and is therefore only exercised on real hardware.
func TestApplicable(t *testing.T) {
	tests := []struct {
		name        string
		version     *string
		noConfig    bool
		wantVersion string
	}{
		{name: "nvm package present", version: ptr("4.80\n"), wantVersion: "4.80"},
		{name: "unpadded version in the package", version: ptr("4.8\n"), wantVersion: "4.80"},
		{name: "no nvm package in this initrd", version: nil},
		{name: "empty version file", version: ptr("  \n")},
		{name: "package without an inventory config", version: ptr("4.80\n"), noConfig: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newIntel(slog.Default())
			i.sysClassNet = fakeSysClassNet(t)
			i.nvmDir = t.TempDir()

			if tt.version != nil {
				if err := os.WriteFile(i.nvmPath("nvmupdate64e"), nil, 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(i.nvmPath("version"), []byte(*tt.version), 0644); err != nil {
					t.Fatal(err)
				}
				if !tt.noConfig {
					if err := os.WriteFile(i.nvmPath(nvmUpdateCfg), nil, 0644); err != nil {
						t.Fatal(err)
					}
				}
			}

			if got := i.applicable(t.Context()); tt.wantVersion == "" && got {
				t.Errorf("applicable got %v want false", got)
			}
			if i.desiredVersion != tt.wantVersion {
				t.Errorf("desired version got %q want %q", i.desiredVersion, tt.wantVersion)
			}
		})
	}
}

func TestNormalizeNVMVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{name: "version as reported by ethtool", version: "4.60", want: "4.60"},
		{name: "unpadded minor of an update package", version: "4.8", want: "4.80"},
		{name: "major with two digits", version: "10.00", want: "10.00"},
		{name: "no minor at all", version: "4", wantErr: true},
		{name: "not numeric", version: "n/a", wantErr: true},
		{name: "more than a major and a minor", version: "4.60.1", wantErr: true},
		{name: "empty", version: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeNVMVersion(tt.version)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}

// TestParseActivation covers the reading of the XML report the update utility writes with
// -o. A card is only activated for what the utility asks for via its top level counters,
// so a run which flashed nothing must not power cycle or reboot this machine.
func TestParseActivation(t *testing.T) {
	tests := []struct {
		name    string
		report  string
		want    Activation
		wantErr bool
	}{
		{
			name: "a card was flashed and needs a power cycle",
			report: `<?xml version="1.0" encoding="UTF-8"?>
<DeviceUpdate lang="en">
	<Instance vendor="8086" device="1593" bus="81" dev="0" func="0">
		<Status result="Success" id="0">Update successful.</Status>
	</Instance>
	<NextUpdateAvailable> 0 </NextUpdateAvailable>
	<RebootRequired> 0 </RebootRequired>
	<PowerCycleRequired> 1 </PowerCycleRequired>
</DeviceUpdate>`,
			want: ActivationPowerCycle,
		},
		{
			name: "a card was flashed and needs a warm reboot",
			report: `<?xml version="1.0" encoding="UTF-8"?>
<DeviceUpdate lang="en">
	<Instance vendor="8086" device="1593" bus="81" dev="0" func="0">
		<Status result="Success" id="0">Update successful.</Status>
	</Instance>
	<NextUpdateAvailable> 0 </NextUpdateAvailable>
	<RebootRequired> 1 </RebootRequired>
	<PowerCycleRequired> 0 </PowerCycleRequired>
</DeviceUpdate>`,
			want: ActivationWarmReboot,
		},
		{
			name: "a power cycle wins over a reboot",
			report: `<DeviceUpdate lang="en">
	<RebootRequired> 1 </RebootRequired>
	<PowerCycleRequired> 1 </PowerCycleRequired>
</DeviceUpdate>`,
			want: ActivationPowerCycle,
		},
		{
			name: "no card took the update",
			report: `<?xml version="1.0" encoding="UTF-8"?>
<DeviceUpdate lang="en">
	<Instance vendor="8086" device="1593" bus="81" dev="0" func="0">
		<Status result="Success" id="0">Update not available.</Status>
	</Instance>
	<NextUpdateAvailable> 0 </NextUpdateAvailable>
	<RebootRequired> 0 </RebootRequired>
	<PowerCycleRequired> 0 </PowerCycleRequired>
</DeviceUpdate>`,
			want: ActivationNone,
		},
		{
			name: "no devices at all",
			report: `<?xml version="1.0" encoding="UTF-8"?>
<DeviceInventory lang="en">
	<Status result="Fail" id="8">No devices on the list.</Status>
</DeviceInventory>`,
			want: ActivationNone,
		},
		{
			name:    "the report is not valid xml",
			report:  "<DeviceUpdate><PowerCycleRequired></DeviceUpdate>",
			wantErr: true,
		},
		{
			name: "a counter is not a number",
			report: `<DeviceUpdate lang="en">
	<PowerCycleRequired> yes </PowerCycleRequired>
</DeviceUpdate>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseActivation([]byte(tt.report))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got activation %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

// fakeSysClassNet creates a sysfs tree with eth0 and eth1 being the two ports of an
// e810, eth2 a card with another driver and lo a virtual interface without a device.
func fakeSysClassNet(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	for interfaceName, pci := range map[string]string{
		"eth0": "0000:51:00.0",
		"eth1": "0000:51:00.1",
		"eth2": "0000:98:00.0",
	} {
		driver := iceDriver
		if interfaceName == "eth2" {
			driver = "i40e"
		}

		device := filepath.Join(root, interfaceName, "device")
		driverPath := filepath.Join(root, "drivers", driver)
		pciPath := filepath.Join(root, "devices", pci)
		for _, dir := range []string{filepath.Join(root, interfaceName), driverPath, pciPath} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.Symlink(pciPath, device); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(driverPath, filepath.Join(device, "driver")); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "lo"), 0755); err != nil {
		t.Fatal(err)
	}

	return root
}

func ptr(s string) *string {
	return &s
}
