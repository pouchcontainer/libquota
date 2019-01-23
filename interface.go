package libquota

import "github.com/pouchcontainer/libquota/types"

// Quota defines the interface of different disk quota
type Quota interface {
	SetQuota(file string, id uint64, quota *types.QuotaLimit) error

	GetQuota(file string) (*types.QuotaLimit, error)

	GetQuotaID(file string) (uint64, error)
}
