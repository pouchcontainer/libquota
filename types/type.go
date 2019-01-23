package types

import "sync"

type FSType string

const (
	Ext4 FSType = "ext4"

	Xfs FSType = "xfs"
)

type BaseQuota struct {
	IDMap     map[uint64]*QuotaLimit
	IDMapLock sync.Mutex
}

type QuotaLimit struct {
	BlockSoftLimit uint64
	BlockHardLimit uint64

	InodeSoftLimit uint64
	InodeHardLimit uint64
}
