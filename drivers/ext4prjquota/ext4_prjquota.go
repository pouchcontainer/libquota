package ext4prjquota

import "github.com/pouchcontainer/libquota/types"

// Ext4PrjQuota defines the ext4 project quota struct
type Ext4PrjQuota struct {
	types.BaseQuota
}

// New is used to check whether support to use ext4 project quota,
// and returns the ext4 project quota object.
func New(file string) (*Ext4PrjQuota, error) {
	// TODO: validate fs info
	return &Ext4PrjQuota{}, nil
}

// SetQuota is used to set the file's ext4 project quota with quota id.
func (q *Ext4PrjQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	return nil
}

// GetQuota returns the file's ext4 project quota information
func (q *Ext4PrjQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

// GetQuotaID returns the file's ext4 project quota id.
func (q *Ext4PrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
