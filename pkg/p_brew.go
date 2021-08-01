// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/executil"

const pathBrew = "/usr/local/bin/brew"

// ManagerBrew is the interface to handle the macOS package manager.
type ManagerBrew struct{}

func (p ManagerBrew) Install(name ...string) error {
	args := []string{pathBrew, "install", "-y"}

	return executil.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerBrew) Remove(name ...string) error {
	args := []string{pathBrew, "uninstall", "-y"}

	return executil.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p ManagerBrew) Purge(name ...string) error {
	return p.Remove(name...)
}

func (p ManagerBrew) Update() error {
	return executil.ExecToStd(nil, sudo, pathBrew, "update")
}

func (p ManagerBrew) Upgrade() error {
	return executil.ExecToStd(nil, sudo, pathBrew, "upgrade")
}

func (p ManagerBrew) Clean() error {
	err := executil.ExecToStd(nil, sudo, pathBrew, "autoremove")
	if err != nil {
		return err
	}

	return executil.ExecToStd(nil, sudo, pathBrew, "cleanup")
}
