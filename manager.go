package libquota

import "sync"

var (
	quotaMap     map[string]Quota
	quotaMapLock sync.Mutex
)

// QuotaMapAdd adds the quota of the mount point
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

// QuotaMapGet returns the quota of the mount point
func QuotaMapGet(mountpoint string) Quota {
	quotaMapLock.Lock()
	defer quotaMapLock.Unlock()

	if quotaMap == nil {
		return nil
	}
	return quotaMap[mountpoint]
}

// QuotaMapDelete deletes the quota of the mount point
func QuotaMapDelete(mountpoint string) {
	quotaMapLock.Lock()
	defer quotaMapLock.Unlock()

	if quotaMap == nil {
		return
	}
	delete(quotaMap, mountpoint)
}
