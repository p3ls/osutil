// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"fmt"
	"os/user"

	"github.com/tredoe/osutil/sh"
	"github.com/tredoe/osutil/sys"
)

// AddGroupFromCmd adds the given group to the original user.
// Returns an information message, if the command is run.
func AddGroupFromCmd(syst sys.System, group string) error {
	switch syst {
	case sys.Linux:
	default:
		panic("unimplemented: " + syst.String())
	}

	username, err := RealUser(syst)
	if err != nil {
		return err
	}

	grp, err := user.LookupGroup(group)
	if err != nil {
		return err
	}
	gid := grp.Gid

	usr, err := user.Lookup(username)
	if err != nil {
		return err
	}
	groups, err := usr.GroupIds()
	if err != nil {
		return err
	}

	found := false
	for _, v := range groups {
		if v == gid {
			found = true
			break
		}
	}
	if !found {
		_, stderr, err := sh.Exec("usermod", "-aG", group, usr.Username)
		if stderr != nil {
			return fmt.Errorf("%s", stderr)
		}
		if err != nil {
			return err
		}

		Log.Printf(
			"the user %q has been added to the group %q.\nYou MUST reboot the system.\n",
			username, group,
		)
		return nil
	}

	return nil
}
