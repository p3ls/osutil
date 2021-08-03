// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import "github.com/tredoe/osutil/executil"

const pathPkg = "/usr/bin/pkg"

// ManagerPkg is the interface to handle the FreeBSD package manager,
// called 'package' or 'pkg'.
type ManagerPkg struct{}

func (p ManagerPkg) Install(name ...string) error {
	args := []string{pathPkg, "install", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Remove(name ...string) error {
	args := []string{pathPkg, "delete", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Purge(name ...string) error {
	args := []string{pathPkg, "purge", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Update() error {
	return executil.RunToStd(nil, sudo, pathPkg, "update")
}

func (p ManagerPkg) Upgrade() error {
	return executil.RunToStd(nil, sudo, pathPkg, "upgrade")
}

func (p ManagerPkg) Clean() error {
	err := executil.RunToStd(nil, sudo, pathPkg, "autoremove")
	if err != nil {
		return err
	}

	return executil.RunToStd(nil, sudo, pathPkg, "clean")
}
