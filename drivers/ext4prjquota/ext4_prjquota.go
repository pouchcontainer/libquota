package ext4prjquota

import "github.com/pouchcontainer/libquota/types"

type ext4PrjQuota struct {
	types.BaseQuota
}

func New(file string) (*ext4PrjQuota, error) {
	// TODO: validate fs info
	return &ext4PrjQuota{}, nil
}

func (q *ext4PrjQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	return nil
}

func (q *ext4PrjQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

func (q *ext4PrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
