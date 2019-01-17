package libquota

import (
	"github.com/pouchcontainer/libquota/pkg/kernel"

	"github.com/pkg/errors"
)

func NewQuota(file string) (Quota, error) {
	var quota Quota

	mountpoint, fsType, err := getMountPoint(file)
	if err != nil {
		return nil, errors.Errorf("failed to get (%s) file system information", file)
	}

	quota = QuotaMapGet(mountpoint)
	if quota != nil {
		return quota, nil
	}

	kernelVersion, err := kernel.GetKernelVersion()

	switch FSType(fsType) {
	case Ext4:
		if (kernelVersion.Kernel == 4 && kernelVersion.Major >= 5) ||
			(kernelVersion.Kernel > 4) {
			quota, err = NewExt4PrjQuota(file)
			if err != nil {
				return nil, err
			}
		} else {
			quota, err = NewExt4GrpQuota(file)
			if err != nil {
				return nil, err
			}
		}
	case Xfs:
		if (kernelVersion.Kernel == 3 && kernelVersion.Major >= 10) ||
			(kernelVersion.Kernel > 3) {
			quota, err = NewXfsPrjQuota(file)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("Unsupport kernel version(%s) for xfs file system",
				kernelVersion.String())
		}
	default:
		return nil, errors.Errorf("Unsupport file system type(%s)", fsType)
	}

	if quota == nil {
		return nil, errors.Errorf("failed to new quota for file(%s)", file)
	}

	QuotaMapAdd(mountpoint, quota)

	return quota, nil
}
