package types

import "sync"

// FSType defines the type of file system struct
type FSType string

const (
	// Ext4 defines the ext4 file system
	Ext4 FSType = "ext4"

	// Xfs defines the xfs file system
	Xfs FSType = "xfs"
)

// BaseQuota defines the basic attribute of quota
type BaseQuota struct {
	IDMap     map[uint64]*QuotaLimit
	IDMapLock sync.Mutex
}

// QuotaLimit defines the attribute of quota limit
type QuotaLimit struct {
	BlockSoftLimit uint64
	BlockHardLimit uint64

	InodeSoftLimit uint64
	InodeHardLimit uint64
}
