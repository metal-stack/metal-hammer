package storage

import (
	"reflect"
	"testing"
)

func Test_assembleSGDiskCommand(t *testing.T) {
	type disk struct {
		disk Disk
	}
	tests := []struct {
		name string
		disk disk
		want []string
	}{
		{
			name: "working",
			disk: disk{
				Disk{
					Device: "/dev/sda",
					Partitions: []*Partition{
						{
							Label:      "efi",
							Number:     1,
							MountPoint: "/boot/efi",
							Filesystem: VFAT,
							GPTType:    GPTBoot,
							GPTGuid:    EFISystemPartition,
							Size:       300,
						},
						{
							Label:      "root",
							Number:     2,
							MountPoint: "/",
							Filesystem: EXT4,
							GPTType:    GPTLinux,
							Size:       -1,
						},
					},
				},
			},
			want: []string{
				"-n=1:0:300M",
				"-c=1:efi",
				"-t=1:ef00",
				"-u=1:C12A7328-F81F-11D2-BA4B-00A0C93EC93B",
				"-n=2:0:0",
				"-c=2:root",
				"-t=2:8300",
				"/dev/sda"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := assembleSGDiskCommand(tt.disk.disk); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("assembleSGDiskCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
