// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package userutil

import (
	"bytes"
	"errors"
	"os/exec"
	"os/user"

	"github.com/tredoe/osutil"
)

var errNoSuperUser = errors.New("you MUST have superuser privileges")

var (
	newLineB = []byte{'\n'}
	emptyB   = []byte{}
)

// CheckSudo executes command 'sudo' to check that the user has permission.
func CheckSudo(syst osutil.System) (err error) {
	if syst == osutil.SystemUndefined {
		if syst, _, err = osutil.SystemFromGOOS(); err != nil {
			return err
		}
	}

	switch syst {
	case osutil.Linux, osutil.FreeBSD, osutil.MacOS:
		cmd := exec.Command("sudo", "true")

		return cmd.Run()
	case osutil.Windows:
		return MustBeSuperUser(osutil.Windows)
	}
	panic("unimplemented: " + syst.String())
}

// MustBeSuperUser checks if the current user is in the superusers group.
// Panics if it is not being run as superuser.
func MustBeSuperUser(syst osutil.System) (err error) {
	if syst == osutil.SystemUndefined {
		if syst, _, err = osutil.SystemFromGOOS(); err != nil {
			return err
		}
	}

	switch syst {
	case osutil.Linux, osutil.FreeBSD, osutil.MacOS, osutil.Windows:
		usr, err := user.Current()
		if err != nil {
			return err
		}
		groups, err := usr.GroupIds()
		if err != nil {
			return err
		}

		findGroup := ""
		switch syst {
		case osutil.Linux, osutil.FreeBSD:
			findGroup = "root"
		case osutil.MacOS:
			findGroup = "admin"
		case osutil.Windows:
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
		panic("unimplemented: " + syst.String())
	}
}

// RealUser returns the original user at Unix systems.
func RealUser(syst osutil.System) (string, error) {
	switch syst {
	default:
		panic("unimplemented: " + syst.String())

	case osutil.Linux:
		username, err := exec.Command("logname").Output()
		if err != nil {
			return "", err
		}
		username = bytes.Replace(username, newLineB, emptyB, 1) // Remove the new line.

		return string(username), nil
	}
}
