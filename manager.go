package libquota

import "sync"

var (
	quotaMap     map[string]Quota
	quotaMapLock sync.Mutex
)

func QuotaMapAdd(mountpoint string, quota Quota) {
	quotaMapLock.Lock()
	defer quotaMapLock.Unlock()

	if mountpoint == "" || quota == nil {
		return
	}

	if quotaMap == nil {
		quotaMap = make(map[string]Quota)
	}

	quotaMap[mountpoint] = quota
}

func QuotaMapGet(mountpoint string) Quota {
	quotaMapLock.Lock()
	defer quotaMapLock.Unlock()

	if quotaMap == nil {
		return nil
	}
	return quotaMap[mountpoint]
}

func QuotaMapDelete(mountpoint string) {
	quotaMapLock.Lock()
	defer quotaMapLock.Unlock()

	if quotaMap == nil {
		return
	}
	delete(quotaMap, mountpoint)
}
