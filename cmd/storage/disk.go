package storage

// FIXME legacy structs remove once old images are gone

type (
	// Disk is a physical Disk
	Disk struct {
		// Device the name of the disk device visible from kernel side, e.g. sda
		Device string
		// Partitions to create on this disk, order is preserved
		Partitions []Partition
	}
	Partition struct {
		Label      string
		Filesystem string
		Properties map[string]string
	}
)
