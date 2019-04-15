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
	// XfsQuota represents the xfs_quota tool
	XfsQuota = "xfs_quota"
)

// XfsPrjQuota defines the xfs project quota struct
type XfsPrjQuota struct {
	types.BaseQuota
}

// New is used to check whether support to use xfs project quota,
// and returns the xfs project quota object.
func New(file string) (*XfsPrjQuota, error) {
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

	return &XfsPrjQuota{
		BaseQuota: types.BaseQuota{
			Mount:  mount,
			IDNext: types.BaseQuotaID,
			IDMap:  make(map[uint64]*types.QuotaLimit),
		},
	}, nil
}

// SetQuota is used to set the file's xfs project quota with quota id.
func (q *XfsPrjQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
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
		id, q.Mount.MountPoint)
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

// GetQuota returns the file's xfs project quota information
func (q *XfsPrjQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	// TODO: Not implemented
	return nil, nil
}

// GetQuotaID returns the file's xfs project quota id.
func (q *XfsPrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
