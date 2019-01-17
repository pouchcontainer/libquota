package libquota

type Quota interface {
	SetQuota(file string, id uint64, quota *QuotaLimit) error

	GetQuota(file string) (*QuotaLimit, error)

	GetQuotaID(file string) (uint64, error)
}
