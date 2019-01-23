package xfsprjquota

import (
	"fmt"
	"strings"

	"github.com/pouchcontainer/libquota/pkg/cmd"
	"github.com/pouchcontainer/libquota/pkg/fs"
	"github.com/pouchcontainer/libquota/types"

	"github.com/pkg/errors"
)

const (
	XfsQuota = "xfs_quota"
)

type xfsPrjQuota struct {
	types.BaseQuota
}

func New(file string) (*xfsPrjQuota, error) {
	// check xfs_quota tool and its version
	res, err := cmd.Run(0, "xfs_quota", "-V")
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, errors.Errorf("failed to get tool xfs_quota, result(%v)", res)
	}

	// remount with prjquota
	mount, err := fs.GetMountPoint(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get file(%s) mountpoint", file)
	}

	if !strings.Contains(mount.Opts, "prjquota") {
		res, err := cmd.Run(0, "mount",
			"-o", "remount,prjquota",
			mount.Device, mount.MountPoint)
		if err != nil {
			return nil, err
		}

		if res.ExitCode != 0 {
			return nil, errors.Errorf("failed to remount with prjquota, result(%v)", res)
		}
	}

	return &xfsPrjQuota{
		BaseQuota: types.BaseQuota{
			IDMap: make(map[uint64]*types.QuotaLimit),
		},
	}, nil
}

func (q *xfsPrjQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	// get mountpoint
	mount, err := fs.GetMountPoint(file)
	if err != nil {
		return errors.Wrapf(err, "failed to get file(%s) mountpoint", file)
	}

	// set quota id
	setID := fmt.Sprintf("project -s -p %s %d", file, id)

	ret, err := cmd.Run(0, XfsQuota, "-xc", setID)
	if err != nil {
		return errors.Wrapf(err, "failed to run (%s -xc %s)", XfsQuota, setID)
	}
	if ret.ExitCode != 0 {
		return errors.Errorf("failed to run (%s -xc %s), result(%v)",
			XfsQuota, setID, ret)
	}

	// set quota limit
	setLimit := fmt.Sprintf("limit -p bhard=%dm bsoft=%dm ihard=%d isoft=%d %d %s",
		quota.BlockHardLimit, quota.BlockSoftLimit,
		quota.InodeHardLimit, quota.InodeSoftLimit,
		id, mount.MountPoint)
	ret, err = cmd.Run(0, XfsQuota, "-xc", setLimit)
	if err != nil {
		return errors.Wrapf(err, "failed to run (%s -xc %s)", XfsQuota, setLimit)
	}
	if ret.ExitCode != 0 {
		return errors.Errorf("failed to run (%s -xc %s), result(%v)",
			XfsQuota, setLimit, ret)
	}

	return nil
}

func (q *xfsPrjQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	// TODO: Not implemented
	return nil, nil
}

func (q *xfsPrjQuota) GetQuotaID(file string) (uint64, error) {
	// TODO: Not implemented
	// xfsctl(path, fd, XFS_IOC_FSGETXATTR, &fsx)
	return 0, nil
}
