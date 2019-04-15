package types

import (
	"github.com/gofrs/flock"
	"github.com/pouchcontainer/libquota/pkg/fs"
)

// FSType defines the type of file system struct
type FSType string

const (
	// Ext4 defines the ext4 file system
	Ext4 FSType = "ext4"

	// Xfs defines the xfs file system
	Xfs FSType = "xfs"

	// BaseQuotaID represents the minimize quota id.
	// The value is unit32(2^24).
	BaseQuotaID = uint64(16777216)
)

// BaseQuota defines the basic attribute of quota
type BaseQuota struct {
	Mount   *fs.Mount
	IDFLock *flock.Flock

	IDNext uint64
	IDMap  map[uint64]*QuotaLimit
}

// QuotaLimit defines the attribute of quota limit
type QuotaLimit struct {
	BlockSoftLimit uint64 // unit is bytes
	BlockHardLimit uint64 // unit is bytes

	InodeSoftLimit uint64
	InodeHardLimit uint64
}
