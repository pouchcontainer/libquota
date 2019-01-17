package libquota

type ext4GrpQuota struct {
	BaseQuota
}

func NewExt4GrpQuota(file string) (*ext4GrpQuota, error) {
	// TODO: validate fs info
	return &ext4GrpQuota{}, nil
}

func (q *ext4GrpQuota) SetQuota(file string, id uint64, quota *QuotaLimit) error {
	return nil
}

func (q *ext4GrpQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
