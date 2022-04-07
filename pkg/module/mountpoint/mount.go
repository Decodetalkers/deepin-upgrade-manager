package mountpoint

import (
	"deepin-upgrade-manager/pkg/module/util"
	"fmt"
	"strings"
)

const (
	_CMD_MOUNT  = "mount"
	_CMD_UMOUNT = "umount"
)

type MountPoint struct {
	Src     string
	Dest    string
	FSType  string
	Options string

	Bind bool
}
type MountPointList []*MountPoint

func (list MountPointList) Mount() (MountPointList, error) {
	var mounted MountPointList
	var err error
	for _, mp := range list {
		err = mp.Mount()
		if err != nil {
			break
		}
		mounted = append(mounted, mp)
	}

	return mounted, err
}

func (list MountPointList) Umount() error {
	var items []string
	for _, mp := range list {
		if isInListByPrefix(mp.Dest, items) {
			continue
		}
		items = append(items, mp.Dest)
	}
	args := []string{"-R"}
	args = append(args, items...)
	fmt.Println("Will umount:", args)
	return util.ExecCommand(_CMD_UMOUNT, args)
}

func (mp *MountPoint) Mount() error {
	err := util.Mkdir(mp.Src, mp.Dest)
	if err != nil {
		return err
	}

	var args []string
	if mp.Bind {
		args = append(args, []string{"-o", "bind"}...)
	} else {
		args = append(args, []string{"-t", mp.FSType}...)
		args = append(args, []string{"-o", mp.Options}...)
	}
	args = append(args, mp.Src)
	args = append(args, mp.Dest)
	fmt.Println("Will mount:", args)
	return util.ExecCommand(_CMD_MOUNT, args)
}

func (mp *MountPoint) Umount() error {
	return util.ExecCommand(_CMD_UMOUNT, []string{"-R", mp.Dest})
}

func isInListByPrefix(item string, list []string) bool {
	for _, v := range list {
		if strings.HasPrefix(item, v) {
			return true
		}
	}
	return false
}
