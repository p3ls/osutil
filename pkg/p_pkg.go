// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/sh"

const pathPkg = "/usr/bin/pkg"

// ManagerPkg is the interface to handle the FreeBSD package manager,
// called 'package' or 'pkg'.
type ManagerPkg struct{}

func (p ManagerPkg) Install(name ...string) error {
	args := []string{pathPkg, "install", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Remove(name ...string) error {
	args := []string{pathPkg, "delete", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Purge(name ...string) error {
	args := []string{pathPkg, "purge", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerPkg) Update() error {
	return sh.ExecToStd(nil, sudo, pathPkg, "update")
}

func (p ManagerPkg) Upgrade() error {
	return sh.ExecToStd(nil, sudo, pathPkg, "upgrade")
}

func (p ManagerPkg) Clean() error {
	err := sh.ExecToStd(nil, sudo, pathPkg, "autoremove")
	if err != nil {
		return err
	}

	return sh.ExecToStd(nil, sudo, pathPkg, "clean")
}
