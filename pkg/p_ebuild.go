// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/sh"

const pathEbuild = "/usr/bin/emerge"

// ManagerEbuild is the interface to handle the package manager of Linux systems based at Gentoo.
type ManagerEbuild struct{}

func (p ManagerEbuild) Install(name ...string) error {
	return sh.ExecToStd(nil, pathEbuild, name...)
}

func (p ManagerEbuild) Remove(name ...string) error {
	args := []string{"--unmerge"}

	return sh.ExecToStd(nil, pathEbuild, append(args, name...)...)
}

func (p ManagerEbuild) Purge(name ...string) error {
	return p.Remove(name...)
}

func (p ManagerEbuild) Update() error {
	return sh.ExecToStd(nil, pathEbuild, "--sync")
}

func (p ManagerEbuild) Upgrade() error {
	return sh.ExecToStd(nil, pathEbuild, "--update", "--deep", "--with-bdeps=y", "--newuse @world")
}

func (p ManagerEbuild) Clean() error {
	err := sh.ExecToStd(nil, pathEbuild, "--update", "--deep", "--newuse @world")
	if err != nil {
		return err
	}

	return sh.ExecToStd(nil, pathEbuild, "--depclean")
}
