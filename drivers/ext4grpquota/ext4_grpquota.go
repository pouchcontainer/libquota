package ext4grpquota

import "github.com/pouchcontainer/libquota/types"

// Ext4GrpQuota defines the ext4 group quota struct
type Ext4GrpQuota struct {
	types.BaseQuota
}

// New is used to check whether support to use ext4 group quota,
// and returns the ext4 group quota object.
func New(file string) (*Ext4GrpQuota, error) {
	// TODO: validate fs info
	return &Ext4GrpQuota{}, nil
}

// SetQuota is used to set the file's ext4 group quota with quota id.
func (q *Ext4GrpQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	return nil
}

// GetQuota returns the file's ext4 group quota information
func (q *Ext4GrpQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

// GetQuotaID returns the file's ext4 group quota id.
func (q *Ext4GrpQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
