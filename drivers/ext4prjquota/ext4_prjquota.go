package ext4prjquota

import (
	"strings"

	"github.com/pouchcontainer/libquota/pkg/cmd"
	"github.com/pouchcontainer/libquota/pkg/fs"
	"github.com/pouchcontainer/libquota/types"

	"github.com/pkg/errors"
)

// Ext4PrjQuota defines the ext4 project quota struct
type Ext4PrjQuota struct {
	types.BaseQuota
}

// New is used to check whether support to use ext4 project quota,
// and returns the ext4 project quota object.
func New(file string) (*Ext4PrjQuota, error) {
	// get mountpoint
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

	// use tool quotaon to set disk quota for mountpoint
	res, err := cmd.Run(0, "quotaon", "-P", mount.MountPoint)
	if err != nil {
		if !strings.Contains(res.Stderr, " File exists") {
			return nil, err
		}
	}

	return &Ext4PrjQuota{
		BaseQuota: types.BaseQuota{
			Mount:  mount,
			IDNext: types.BaseQuotaID,
			IDMap:  make(map[uint64]*types.QuotaLimit),
		},
	}, nil
}

// SetQuota is used to set the file's ext4 project quota with quota id.
func (q *Ext4PrjQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	return nil
}

// GetQuota returns the file's ext4 project quota information
func (q *Ext4PrjQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

// GetQuotaID returns the file's ext4 project quota id.
func (q *Ext4PrjQuota) GetQuotaID(file string) (uint64, error) {
	return 0, nil
}
