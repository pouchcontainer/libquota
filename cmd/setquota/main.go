package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pouchcontainer/libquota"
	"github.com/pouchcontainer/libquota/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// version will be populated by the Makefile, read from
// VERSION file of the source code.
var version = ""

// gitCommit will be the hash that the binary was built from
// and will be populated by the Makefile
var gitCommit = ""

const (
	usage = `Set disk quota tool, support ext4 and xfs file system`
)

func main() {
	app := cli.NewApp()
	app.Name = "setquota"
	app.Usage = usage

	var v []string
	if version != "" {
		v = append(v, version)
	}
	if gitCommit != "" {
		v = append(v, fmt.Sprintf("commit: %s", gitCommit))
	}
	app.Version = strings.Join(v, "\n")

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "/var/log/setquota.log",
			Usage: "set the log file path where internal debug information is written",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "set the file or directory's disk quota",
		},
		cli.Uint64Flag{
			Name:  "id",
			Value: 0,
			Usage: "set the quota's id",
		},
		cli.Uint64Flag{
			Name:  "bhard",
			Usage: "set block hard limit, unit(MB)",
		},
		cli.Uint64Flag{
			Name:  "bsoft",
			Usage: "set block soft limit, unit(MB)",
		},
		cli.Uint64Flag{
			Name:  "ihard",
			Usage: "set inode hard limit",
		},
		cli.Uint64Flag{
			Name:  "isoft",
			Usage: "set inode soft limit",
		},
	}

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if path := context.GlobalString("log"); path != "" {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
			if err != nil {
				return err
			}
			logrus.SetOutput(f)
		}
		return nil
	}

	app.Action = func(context *cli.Context) error {
		return setquota(context)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setquota(context *cli.Context) error {
	file := context.String("file")
	if file == "" {
		return errors.Errorf("")
	}

	id := context.Uint64("id")

	bhard := context.Uint64("bhard")
	bsoft := context.Uint64("bsoft")
	ihard := context.Uint64("ihard")
	isoft := context.Uint64("isoft")
	if bhard == 0 && bsoft == 0 && ihard == 0 && isoft == 0 {
		return errors.Errorf("haven't set quota limit, " +
			"set one of 'bhard' or 'bsoft' or 'ihard' or 'isoft'")
	}
	limit := &types.QuotaLimit{
		BlockHardLimit: bhard * 1024 * 1024,
		BlockSoftLimit: bsoft * 1024 * 1024,
		InodeHardLimit: ihard,
		InodeSoftLimit: isoft,
	}

	quota, err := libquota.NewQuota(file)
	if err != nil {
		return err
	}

	if id == 0 {
		id, err = quota.GetQuotaID(file)
		if err != nil {
			return err
		}
	}

	return quota.SetQuota(file, id, limit)
}
