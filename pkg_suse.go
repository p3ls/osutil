// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: SUSE

package osutil

import "github.com/tredoe/osutil/executil"

const fileZypp = "zypper"

var pathZypp = "/usr/bin/zypper"

// ManagerZypp is the interface to handle the package manager of Linux systems based at SUSE.
type ManagerZypp struct {
	pathExec string
}

// NewManagerZypp returns the Zypp package manager.
func NewManagerZypp() ManagerZypp {
	return ManagerZypp{pathExec: pathZypp}
}

func (m ManagerZypp) setExecPath(p string) { m.pathExec = p }

func (m ManagerZypp) ExecPath() string { return m.pathExec }

func (m ManagerZypp) PackageType() string { return Zypp.String() }

func (m ManagerZypp) Install(name ...string) error {
	args := []string{
		pathZypp, "install", "--non-interactive", "--auto-agree-with-licenses", "-y",
	}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerZypp) Remove(name ...string) error {
	args := []string{pathZypp, "remove", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerZypp) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerZypp) Update() error {
	return executil.RunToStd(nil, sudo, pathZypp, "refresh")
}

func (m ManagerZypp) Upgrade() error {
	return executil.RunToStd(nil, sudo, pathZypp, "up", "--auto-agree-with-licenses", "-y")
}

func (m ManagerZypp) Clean() error {
	return executil.RunToStd(nil, sudo, pathZypp, "clean")
}
