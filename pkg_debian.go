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
	cmd      *executil.Command
}

// NewManagerDeb returns the Deb package manager.
func NewManagerDeb() ManagerDeb {
	return ManagerDeb{
		pathExec: pathDeb,
		cmd: excmd.Command("", "").
			AddEnv([]string{"DEBIAN_FRONTEND=noninteractive"}).
			BadExitCodes([]int{100}),
	}
}

func (m ManagerDeb) setExecPath(p string) { m.pathExec = p }

func (m ManagerDeb) ExecPath() string { return m.pathExec }

func (m ManagerDeb) PackageType() string { return Deb.String() }

func (m ManagerDeb) Install(name ...string) error {
	args := []string{pathDeb, "install", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerDeb) Remove(name ...string) error {
	args := []string{pathDeb, "remove", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerDeb) Purge(name ...string) error {
	args := []string{pathDeb, "purge", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerDeb) Update() error {
	_, err := m.cmd.Command(sudo, pathDeb, "update", "-qq").Run()
	return err
}

func (m ManagerDeb) Upgrade() error {
	_, err := m.cmd.Command(sudo, pathDeb, "upgrade", "-y").Run()
	return err
}

func (m ManagerDeb) Clean() error {
	_, err := m.cmd.Command(sudo, pathDeb, "autoremove", "-y").Run()
	if err != nil {
		return err
	}

	_, err = m.cmd.Command(sudo, pathDeb, "clean").Run()
	return err
}
