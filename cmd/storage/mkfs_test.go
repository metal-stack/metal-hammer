package storage

import (
	"reflect"
	"testing"
)

func Test_assembleMKFSCommand(t *testing.T) {
	type partition struct {
		p *Partition
	}
	tests := []struct {
		name      string
		partition partition
		mkfs      string
		mkfsargs  []string
		wantErr   bool
	}{
		{
			name:      "Working ext4",
			partition: partition{&Partition{Device: "/dev/sda", Filesystem: EXT4}},
			mkfs:      ext4MkFsCommand,
			mkfsargs:  []string{"-v", "-F", "/dev/sda"},
			wantErr:   false,
		},
		{
			name:      "Working ext3",
			partition: partition{&Partition{Device: "/dev/sda", Filesystem: EXT3, Label: "root"}},
			mkfs:      ext3MkFsCommand,
			mkfsargs:  []string{"-v", "-F", "-L", "root", "/dev/sda"},
			wantErr:   false,
		},
		{
			name:      "Working vfat",
			partition: partition{&Partition{Device: "/dev/sda", Filesystem: VFAT, Label: "efi"}},
			mkfs:      fat32MkFsCommand,
			mkfsargs:  []string{"-v", "-F", "32", "-n", "EFI", "/dev/sda"},
			wantErr:   false,
		},
		{
			name:      "Working swap",
			partition: partition{&Partition{Device: "/dev/sda", Filesystem: SWAP, Label: "swap"}},
			mkfs:      mkswapCommand,
			mkfsargs:  []string{"-f", "-L", "swap", "/dev/sda"},
			wantErr:   false,
		},
		{
			name:      "Not working",
			partition: partition{&Partition{Device: "/dev/sda", Filesystem: "unknown"}},
			mkfs:      "",
			mkfsargs:  []string{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mkfs, mkfsargs, err := assembleMKFSCommand(tt.partition.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("assembleMKFSCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if mkfs != tt.mkfs {
				t.Errorf("assembleMKFSCommand() mkfs = %v, want %v", mkfs, tt.mkfs)
			}
			if !reflect.DeepEqual(mkfsargs, tt.mkfsargs) {
				t.Errorf("assembleMKFSCommand() mkfsargs = %v, want %v", mkfsargs, tt.mkfsargs)
			}
		})
	}
}
