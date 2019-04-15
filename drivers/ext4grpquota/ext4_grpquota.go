package ext4grpquota

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pouchcontainer/libquota/pkg/cmd"
	"github.com/pouchcontainer/libquota/pkg/fs"
	"github.com/pouchcontainer/libquota/types"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Ext4GrpQuota defines the ext4 group quota struct
type Ext4GrpQuota struct {
	types.BaseQuota
}

// New is used to check whether support to use ext4 group quota,
// and returns the ext4 group quota object.
func New(file string) (*Ext4GrpQuota, error) {
	// get mountpoint
	mount, err := fs.GetMountPoint(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get file(%s) mountpoint", file)
	}

	if !strings.Contains(mount.Opts, "grpquota") {
		res, err := cmd.Run(0, "mount",
			"-o", "remount,grpquota",
			mount.Device, mount.MountPoint)
		if err != nil {
			return nil, err
		}

		if res.ExitCode != 0 {
			return nil, errors.Errorf("failed to remount with prjquota, result(%v)", res)
		}
	}

	//
	vfsVersion, quotaFilename, err := getVFSVersionAndQuotaFile(mount.MountPoint)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vfs version and quota file")
	}

	filename := mount.MountPoint + "/" + quotaFilename
	if _, err := os.Stat(filename); err != nil {
		os.Remove(mount.MountPoint + "/aquota.user")

		header := []byte{0x27, 0x19, 0xc0, 0xd9, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3a, 0x09, 0x00, 0x80,
			0x3a, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x05, 0x00, 0x00, 0x00}
		if vfsVersion == "vfsv1" {
			header[4] = 0x01
		}

		if writeErr := ioutil.WriteFile(filename, header, 0644); writeErr != nil {
			return nil, errors.Wrapf(writeErr, "failed to write file, filename: (%s), vfs version: (%s)",
				filename, vfsVersion)
		}
		if res, err := cmd.Run(0, "setquota", "-g", "-t", "43200", "43200", mount.MountPoint); err != nil {
			os.Remove(filename)
			return nil, errors.Wrapf(err, "failed to setquota, result(%v)", res)
		}
		if res, err := cmd.Run(0, "setquota", "-g", "0", "0", "0", "0", "0", mount.MountPoint); err != nil {
			os.Remove(filename)
			return nil, errors.Wrapf(err, "failed to setquota, result(%v)", res)
		}
	}

	// check group quota status, on or not, pay attention, the right exit code of command 'quotaon' is '1'.
	res, err := cmd.Run(0, "quotaon", "-pg", mount.MountPoint)
	if err != nil && res.ExitCode != 1 {
		return nil, errors.Wrapf(err, "failed to quota on for mountpoint(%s), result(%v)",
			mount.MountPoint, res)
	}

	// set quota on
	if !strings.Contains(res.Stdout, " is on") {
		if res, err = cmd.Run(0, "quotaon", mount.MountPoint); err != nil {
			return nil, errors.Wrapf(err, "failed to quotaon, mountpoint(%s), result(%v)",
				mount.MountPoint, res)
		}
	}

	idFLock := flock.New(filepath.Join(mount.MountPoint, "quota-id.lock"))
	idFLock.Lock()
	defer idFLock.Unlock()

	idNext, idMap, err := loadQuotaIDs(mount.MountPoint, "-gn")
	if err != nil {
		return nil, err
	}

	return &Ext4GrpQuota{
		BaseQuota: types.BaseQuota{
			Mount:   mount,
			IDFLock: idFLock,
			IDNext:  idNext,
			IDMap:   idMap,
		},
	}, nil
}

// SetQuota is used to set the file's ext4 group quota with quota id.
func (q *Ext4GrpQuota) SetQuota(file string, id uint64, quota *types.QuotaLimit) error {
	if file == "" || id == 0 || quota == nil {
		return errors.Errorf("invalid arguments, file(%s), id(%d), quota(%v)", file, id, quota)
	}

	// set quota id on files
	args := fmt.Sprintf("-n system.subtree -v %d %s", id, file)
	res, err := cmd.Run(0, "setfattr", strings.Fields(args)...)
	if err != nil || res.ExitCode != 0 {
		return errors.Wrapf(err, "failed to set quota, mountpoint(%s), result(%v)", q.Mount.MountPoint, res)
	}

	// set quota on mountpoint
	args = fmt.Sprintf("-g %d %d %d %d %d %s",
		id, quota.BlockSoftLimit/1024, quota.BlockHardLimit/1024, quota.InodeSoftLimit, quota.InodeHardLimit, q.Mount.MountPoint)

	res, err = cmd.Run(0, "setquota", strings.Fields(args)...)
	if err != nil || res.ExitCode != 0 {
		return errors.Wrapf(err, "failed to set quota, mountpoint(%s), result(%v)", q.Mount.MountPoint, res)
	}

	return nil
}

// GetQuota returns the file's ext4 group quota information
func (q *Ext4GrpQuota) GetQuota(file string) (*types.QuotaLimit, error) {
	return nil, nil
}

// GetQuotaID returns the file's ext4 group quota id.
func (q *Ext4GrpQuota) GetQuotaID(file string) (uint64, error) {
	var (
		getNextID bool
		id        uint64
	)

	if file == "" {
		getNextID = true
	}

	if _, err := os.Stat(file); err != nil {
		return 0, err
	}

	res, err := cmd.Run(0, "getfattr", "-n", "system.subtree", "--only-values", "--absolute-names", file)
	if err == nil && res.ExitCode == 0 {
		num, err := strconv.Atoi(res.Stdout)
		if err != nil {
			return 0, err
		}
		id = uint64(num)

		if id == 0 {
			getNextID = true
		}
	} else {
		getNextID = true
	}

	if getNextID {
		id = q.GetNextQuotaID()
	}

	return id, nil
}

// GetNextQuotaID returns the quota of this mountpoint that can be used.
func (q *Ext4GrpQuota) GetNextQuotaID() uint64 {
	q.IDFLock.Lock()
	defer func() {
		q.IDNext++
		q.IDFLock.Unlock()
	}()

	return q.IDNext
}

func getVFSVersionAndQuotaFile(file string) (string, string, error) {
	devID, err := fs.GetDevID(file)
	if err != nil {
		return "", "", err
	}

	output, err := ioutil.ReadFile(fs.ProcMount)
	if err != nil {
		logrus.Warnf("failed to read file: (%s), err: (%v)", fs.ProcMount, err)
		return "", "", errors.Wrap(err, "failed to read /proc/mounts")
	}

	vfsVersion := "vfsv0"
	quotaFilename := "aquota.group"
	for _, line := range strings.Split(string(output), "\n") {
		// TODO: add an example here to make following code readable.
		// /dev/sdb1 /home/pouch ext4 rw,relatime,prjquota,data=ordered 0 0 ?
		parts := strings.Split(line, " ")
		if len(parts) != 6 {
			continue
		}

		devID2, _ := fs.GetDevID(parts[1])
		if devID != devID2 {
			continue
		}

		for _, opt := range strings.Split(parts[3], ",") {
			items := strings.SplitN(opt, "=", 2)
			if len(items) != 2 {
				continue
			}
			switch items[0] {
			case "jqfmt":
				vfsVersion = items[1]
			case "grpjquota":
				quotaFilename = items[1]
			}
		}
		return vfsVersion, quotaFilename, nil
	}

	return vfsVersion, quotaFilename, nil
}

func loadQuotaIDs(mountpoint string, opts string) (uint64, map[uint64]*types.QuotaLimit, error) {
	idMap := make(map[uint64]*types.QuotaLimit)

	res, err := cmd.Run(0, "repquota", opts, mountpoint)
	if err != nil || res.ExitCode != 0 {
		return 0, nil, errors.Wrapf(err, "failed to execute [repquota %s %s], result(%v)",
			opts, mountpoint, res)
	}

	// example output:
	//                          Block limits                File limits
	// Group           used    soft    hard  grace    used  soft  hard  grace
	// ----------------------------------------------------------------------
	// #0        -- 12693804       0       0         163764     0     0
	// #16777217 --      44       0 10485760             18     0     0

	currentID := types.BaseQuotaID
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		// find all lines with prefix '#'
		if len(line) == 0 || line[0] != '#' {
			continue
		}

		parts := strings.Fields(line)
		// right line is:
		// #16777217 --      44       0 10485760             18     0     0
		if len(parts) != 8 {
			continue
		}

		id, err := strconv.Atoi(parts[0][1:])
		if err != nil {
			continue
		}

		qid := uint64(id)
		if qid < types.BaseQuotaID {
			continue
		}

		if qid > currentID {
			currentID = qid
		}

		bsLimit, _ := strconv.Atoi(parts[3])
		bhLimit, _ := strconv.Atoi(parts[4])
		isLimit, _ := strconv.Atoi(parts[6])
		ihLimit, _ := strconv.Atoi(parts[7])
		idMap[qid] = &types.QuotaLimit{
			BlockSoftLimit: uint64(bsLimit) * 1024,
			BlockHardLimit: uint64(bhLimit) * 1024,
			InodeSoftLimit: uint64(isLimit),
			InodeHardLimit: uint64(ihLimit),
		}
	}

	return currentID + 1, idMap, nil
}
