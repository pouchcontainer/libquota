package libquota

type xfsPrjQuota struct {
	BaseQuota
}

func NewXfsPrjQuota(file string) (*xfsPrjQuota, error) {
	// TODO: validate fs info
	return &xfsPrjQuota{}, nil
}

func (q *xfsPrjQuota) SetQuota(file string, id uint64, quota *QuotaLimit) error {
	return nil
}

func (q *xfsPrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
