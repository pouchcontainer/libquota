package fs

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pouchcontainer/libquota/pkg/cmd"

	"github.com/pkg/errors"
)

const (
	// ProcMount represents the /proc/mounts file
	ProcMount = "/proc/mounts"
)

// Mount defines the attribute of mount.
type Mount struct {
	Device     string
	MountPoint string
	FSType     string
	Opts       string // fs options, such as rw, ro.
	Dump       string // backup options, set 0 or 1, 1 will use dump to backup.
	Pass       string // fsck options, set 0 or 1 or 2, 0 will not check, 1 earlier than 2 to check.
}

// GetMountPoint returns the mount point information of file system
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
	if len(lines) != 3 {
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

		// check device is exist or not.
		if _, err := os.Stat(parts[0]); err != nil {
			continue
		}

		if parts[1] == mountpoint {
			mount = &Mount{
				Device:     parts[0],
				MountPoint: parts[1],
				FSType:     parts[2],
				Opts:       parts[3],
				Dump:       parts[4],
				Pass:       parts[5],
			}
			break
		}
	}

	return mount, nil
}
