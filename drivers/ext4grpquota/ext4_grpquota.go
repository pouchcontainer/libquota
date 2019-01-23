package ext4grpquota

import "github.com/pouchcontainer/libquota/types"

type ext4GrpQuota struct {
	types.BaseQuota
}

func New(file string) (*ext4GrpQuota, error) {
	// TODO: validate fs info
	return &ext4GrpQuota{}, nil
}

func (q *ext4GrpQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	return nil
}

func (q *ext4GrpQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

func (q *ext4GrpQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
