// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: Debian

package osutil

import "github.com/tredoe/osutil/executil"

// 'apt' is for the terminal and gives beautiful output.
// 'apt-get' and 'apt-cache' are for scripts and give stable, parsable output.

const (
	fileDeb = "apt-get"
	pathDeb = "/usr/bin/apt-get"
)

// ManagerDeb is the interface to handle the package manager of Linux systems based at Debian.
type ManagerDeb struct {
	pathExec string
}

// NewManagerDeb returns the Deb package manager.
func NewManagerDeb() ManagerDeb {
	return ManagerDeb{pathExec: pathDeb}
}

func (m ManagerDeb) setExecPath(p string) { m.pathExec = p }

func (m ManagerDeb) ExecPath() string { return m.pathExec }

func (m ManagerDeb) PackageType() string { return Deb.String() }

func (m ManagerDeb) Install(name ...string) error {
	args := []string{pathDeb, "install", "-y"}

	return executil.RunToStd(
		[]string{"DEBIAN_FRONTEND=noninteractive"},
		sudo, append(args, name...)...,
	)
}

func (m ManagerDeb) Remove(name ...string) error {
	args := []string{pathDeb, "remove", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerDeb) Purge(name ...string) error {
	args := []string{pathDeb, "purge", "-y"}

	return executil.RunToStd(nil, sudo, append(args, name...)...)
}

func (m ManagerDeb) Update() error {
	return executil.RunToStd(nil, sudo, pathDeb, "update", "-qq")
}

func (m ManagerDeb) Upgrade() error {
	return executil.RunToStd(nil, sudo, pathDeb, "upgrade", "-y")
}

func (m ManagerDeb) Clean() error {
	if err := executil.RunToStd(nil, sudo, pathDeb, "autoremove", "-y"); err != nil {
		return err
	}
	return executil.RunToStd(nil, sudo, pathDeb, "clean")
}
