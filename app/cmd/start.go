package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/backing-image-manager/pkg/manager"
	filesync "github.com/longhorn/backing-image-manager/pkg/sync"
	"github.com/longhorn/backing-image-manager/pkg/types"
	"github.com/longhorn/backing-image-manager/pkg/util"
)

func StartCmd() cli.Command {
	return cli.Command{
		Name: "daemon",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen",
				Value: "localhost:" + strconv.Itoa(types.DefaultManagerPort),
			},
			cli.StringFlag{
				Name:  "sync-listen",
				Value: "localhost:" + strconv.Itoa(types.DefaultSyncServerPort),
			},
			cli.StringFlag{
				Name:  "disk-uuid",
				Usage: "The corresponding disk uuid stored in the metafile of the disk path",
			},
			cli.StringFlag{
				Name:  "port-range",
				Value: "30001-31000",
			},
		},
		Action: func(c *cli.Context) {
			if err := start(c); err != nil {
				logrus.Fatalf("Error running start command: %v.", err)
			}
		},
	}
}

func start(c *cli.Context) error {
	listen := c.String("listen")
	syncListen := c.String("sync-listen")
	diskUUID := c.String("disk-uuid")
	portRange := c.String("port-range")

	diskUUIDInFile, err := util.GetDiskConfig(types.DiskPathInContainer)
	if err != nil {
		return err
	}
	if diskUUID == "" {
		diskUUID = diskUUIDInFile
	} else if diskUUID != diskUUIDInFile {
		return fmt.Errorf("invalid input disk UUID %v, which doesn't match disk UUID %v the disk config file", diskUUID, diskUUIDInFile)
	}

	return manager.NewServer(context.Background(), listen, syncListen, diskUUID, types.DiskPathInContainer, portRange, &filesync.HTTPHandler{})
}
