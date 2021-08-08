// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// System: macOS
// Note: running Homebrew as root is extremely dangerous and no longer supported.

package osutil

import "github.com/tredoe/osutil/executil"

const (
	fileBrew = "brew"
	pathBrew = "/usr/local/bin/brew"
)

// ManagerBrew is the interface to handle the macOS package manager.
type ManagerBrew struct {
	pathExec string
}

// NewManagerBrew returns the Homebrew package manager.
func NewManagerBrew() ManagerBrew {
	return ManagerBrew{pathExec: pathBrew}
}

func (m ManagerBrew) setExecPath(p string) { m.pathExec = p }

func (m ManagerBrew) ExecPath() string { return m.pathExec }

func (m ManagerBrew) PackageType() string { return Brew.String() }

func (m ManagerBrew) Install(name ...string) error {
	args := []string{"install", "-y"}

	return executil.RunToStd(nil, pathBrew, append(args, name...)...)
}

func (m ManagerBrew) Remove(name ...string) error {
	args := []string{"uninstall", "-y"}

	return executil.RunToStd(nil, pathBrew, append(args, name...)...)
}

func (m ManagerBrew) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerBrew) Update() error {
	return executil.RunToStd(nil, pathBrew, "update")
}

func (m ManagerBrew) Upgrade() error {
	return executil.RunToStd(nil, pathBrew, "upgrade")
}

func (m ManagerBrew) Clean() error {
	if err := executil.RunToStd(nil, pathBrew, "autoremove"); err != nil {
		return err
	}
	return executil.RunToStd(nil, pathBrew, "cleanup")
}
