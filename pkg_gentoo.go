// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: Gentoo

package osutil

import "github.com/tredoe/osutil/executil"

const pathEbuild = "/usr/bin/emerge"

// ManagerEbuild is the interface to handle the package manager of Linux systems based at Gentoo.
type ManagerEbuild struct{}

func (m ManagerEbuild) ExecPath() string {
	return pathEbuild
}

func (m ManagerEbuild) Install(name ...string) error {
	return executil.RunToStd(nil, pathEbuild, name...)
}

func (m ManagerEbuild) Remove(name ...string) error {
	args := []string{"--unmerge"}

	return executil.RunToStd(nil, pathEbuild, append(args, name...)...)
}

func (m ManagerEbuild) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerEbuild) Update() error {
	return executil.RunToStd(nil, pathEbuild, "--sync")
}

func (m ManagerEbuild) Upgrade() error {
	return executil.RunToStd(nil, pathEbuild, "--update", "--deep", "--with-bdeps=y", "--newuse @world")
}

func (m ManagerEbuild) Clean() error {
	err := executil.RunToStd(nil, pathEbuild, "--update", "--deep", "--newuse @world")
	if err != nil {
		return err
	}
	return executil.RunToStd(nil, pathEbuild, "--depclean")
}
