// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// System: FreeBSD

package osutil

import "github.com/tredoe/osutil/executil"

const pathPkg = "/usr/bin/pkg"

// ManagerPkg is the interface to handle the FreeBSD package manager,
// called 'package' or 'pkg'.
type ManagerPkg struct{}

func (m ManagerPkg) ExecPath() string {
	return pathPkg
}

func (m ManagerPkg) Install(name ...string) error {
	args := []string{pathPkg, "install", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerPkg) Remove(name ...string) error {
	args := []string{pathPkg, "delete", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerPkg) Purge(name ...string) error {
	args := []string{pathPkg, "purge", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerPkg) Update() error {
	return executil.RunToStd(nil, sudo, pathPkg, "update")
}

func (m ManagerPkg) Upgrade() error {
	return executil.RunToStd(nil, sudo, pathPkg, "upgrade")
}

func (m ManagerPkg) Clean() error {
	if err := executil.RunToStd(nil, sudo, pathPkg, "autoremove"); err != nil {
		return err
	}
	return executil.RunToStd(nil, sudo, pathPkg, "clean")
}
