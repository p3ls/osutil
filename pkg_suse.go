// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: SUSE

package osutil

import "github.com/tredoe/osutil/executil"

const pathZypp = "/usr/bin/zypper"

// ManagerZypp is the interface to handle the package manager of Linux systems based at SUSE.
type ManagerZypp struct{}

func (m ManagerZypp) ExecPath() string {
	return pathZypp
}

func (m ManagerZypp) Install(name ...string) error {
	args := []string{"install", "--auto-agree-with-licenses"}

	return executil.RunToStd(nil, pathZypp, append(args, name...)...)
}

func (m ManagerZypp) Remove(name ...string) error {
	args := []string{"remove"}

	return executil.RunToStd(nil, pathZypp, append(args, name...)...)
}

func (m ManagerZypp) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerZypp) Update() error {
	return executil.RunToStd(nil, pathZypp, "refresh")
}

func (m ManagerZypp) Upgrade() error {
	return executil.RunToStd(nil, pathZypp, "up", "--auto-agree-with-licenses")
}

func (m ManagerZypp) Clean() error {
	return executil.RunToStd(nil, pathZypp, "clean")
}
