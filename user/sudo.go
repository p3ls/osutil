// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"bytes"
	"errors"
	"os/exec"
	"os/user"

	"github.com/tredoe/osutil/sys"
)

var errNoSuperUser = errors.New("you MUST have superuser privileges")

var (
	newLineB = []byte{'\n'}
	emptyB   = []byte{}
)

// CheckSudo executes command 'sudo' to check that the user has permission.
func CheckSudo(sys_ sys.System) (err error) {
	if sys_ == sys.SysUndefined {
		if sys_, _, err = sys.SystemFromGOOS(); err != nil {
			return err
		}
	}

	switch sys_ {
	case sys.SysLinux, sys.SysFreeBSD, sys.SysMacOS:
		cmd := exec.Command("sudo", "true")

		return cmd.Run()
	case sys.SysWindows:
		return MustBeSuperUser(sys.SysWindows)
	}
	panic("unimplemented: " + sys_.String())
}

// MustBeSuperUser checks if the current user is in the superusers group.
// Panics if it is not being run as superuser.
func MustBeSuperUser(sys_ sys.System) error {
	switch sys_ {
	case sys.SysLinux, sys.SysFreeBSD, sys.SysMacOS, sys.SysWindows:
		usr, err := user.Current()
		if err != nil {
			return err
		}
		groups, err := usr.GroupIds()
		if err != nil {
			return err
		}

		findGroup := ""
		switch sys_ {
		case sys.SysLinux, sys.SysFreeBSD:
			findGroup = "root"
		case sys.SysMacOS:
			findGroup = "admin"
		case sys.SysWindows:
			findGroup = "Administrators"
		}

		for _, v := range groups {
			grp, err := user.LookupGroupId(v)
			if err != nil {
				return err
			}
			if grp.Name == findGroup {
				return nil
			}
		}
		//return errNoSuperUser
		panic(errNoSuperUser)

	default:
		panic("unimplemented: " + sys_.String())
	}
}

// RealUser returns the original user at Unix systems.
func RealUser(sys_ sys.System) (string, error) {
	switch sys_ {
	default:
		panic("unimplemented: " + sys_.String())

	case sys.SysLinux:
		username, err := exec.Command("logname").Output()
		if err != nil {
			return "", err
		}
		username = bytes.Replace(username, newLineB, emptyB, 1) // Remove the new line.

		return string(username), nil
	}
}
