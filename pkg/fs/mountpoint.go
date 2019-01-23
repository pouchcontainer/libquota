package fs

import (
	"io/ioutil"
	"strings"

	"github.com/pouchcontainer/libquota/pkg/cmd"

	"github.com/pkg/errors"
)

const (
	ProcMount = "/proc/mounts"
)

type Mount struct {
	Device     string
	MountPoint string
	FSType     string
	Opts       string // fs options, such as rw, ro.
	Dump       string // backup options, set 0 or 1, 1 will use dump to backup.
	Pass       string // fsck options, set 0 or 1 or 2, 0 will not check, 1 earlier than 2 to check.
}

// return mountpoint, fs type, error
func GetMountPoint(file string) (*Mount, error) {
	var (
		mountpoint string
		mount      *Mount
	)

	res, err := cmd.Run(0, "df", file)
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, errors.Errorf("failed to run (df %s), result(%v)",
			file, res)
	}

	// example output is:
	// Filesystem     1K-blocks  Used Available Use% Mounted on
	// /dev/sdb        41922560 32940  41889620   1% /mnt/data
	lines := strings.Split(res.Stdout, "\n")
	if len(lines) != 2 {
		return nil, errors.Errorf("failed to use df to get mountpoint, "+
			"invalid output(%s)", res.Stdout)
	}

	parts := strings.Fields(lines[1])
	if len(parts) != 6 {
		return nil, errors.Errorf("failed to use df to get mountpoint, "+
			"invalid output(%s)", lines[1])
	}
	mountpoint = parts[5]

	mountInfo, err := ioutil.ReadFile(ProcMount)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file(%s)", ProcMount)
	}

	for _, line := range strings.Split(string(mountInfo), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) != 6 {
			continue
		}
		if parts[1] == mountpoint {
			mount.Device = parts[0]
			mount.MountPoint = parts[1]
			mount.FSType = parts[2]
			mount.Opts = parts[3]
			mount.Dump = parts[4]
			mount.Pass = parts[5]
			break
		}
	}

	return mount, nil
}
