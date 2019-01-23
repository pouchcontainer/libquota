package libquota

import (
	"github.com/pouchcontainer/libquota/drivers/ext4grpquota"
	"github.com/pouchcontainer/libquota/drivers/ext4prjquota"
	"github.com/pouchcontainer/libquota/drivers/xfsprjquota"
	"github.com/pouchcontainer/libquota/pkg/fs"
	"github.com/pouchcontainer/libquota/pkg/kernel"
	"github.com/pouchcontainer/libquota/types"

	"github.com/pkg/errors"
)

func NewQuota(file string) (Quota, error) {
	var quota Quota

	mount, err := fs.GetMountPoint(file)
	if err != nil {
		return nil, errors.Errorf("failed to get (%s) file system information", file)
	}

	quota = QuotaMapGet(mount.MountPoint)
	if quota != nil {
		return quota, nil
	}

	kernelVersion, err := kernel.GetKernelVersion()

	switch types.FSType(mount.FSType) {
	case types.Ext4:
		if (kernelVersion.Kernel == 4 && kernelVersion.Major >= 5) ||
			(kernelVersion.Kernel > 4) {
			quota, err = ext4prjquota.New(file)
			if err != nil {
				return nil, err
			}
		} else {
			quota, err = ext4grpquota.New(file)
			if err != nil {
				return nil, err
			}
		}
	case types.Xfs:
		if (kernelVersion.Kernel == 3 && kernelVersion.Major >= 10) ||
			(kernelVersion.Kernel > 3) {
			quota, err = xfsprjquota.New(file)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("Unsupport kernel version(%s) for xfs file system",
				kernelVersion.String())
		}
	default:
		return nil, errors.Errorf("Unsupport file system type(%s)", mount.FSType)
	}

	if quota == nil {
		return nil, errors.Errorf("failed to new quota for file(%s)", file)
	}

	QuotaMapAdd(mount.MountPoint, quota)

	return quota, nil
}
