package libquota

type ext4PrjQuota struct {
	BaseQuota
}

func NewExt4PrjQuota(file string) (*ext4PrjQuota, error) {
	// TODO: validate fs info
	return &ext4PrjQuota{}, nil
}

func (q *ext4PrjQuota) SetQuota(file string, id uint64, quota *QuotaLimit) error {
	return nil
}

func (q *ext4PrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
