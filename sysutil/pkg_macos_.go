// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// System: macOS
// Note: running Homebrew as root is extremely dangerous and no longer supported.

package sysutil

import "github.com/tredoe/osutil/executil"

const (
	fileBrew = "brew"
	pathBrew = "/usr/local/bin/brew"
)

// ManagerBrew is the interface to handle the macOS package manager.
type ManagerBrew struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerBrew returns the Homebrew package manager.
func NewManagerBrew() ManagerBrew {
	return ManagerBrew{
		pathExec: pathBrew,
		cmd: excmd.Command("", "").
			BadExitCodes([]int{1}),
	}
}

func (m ManagerBrew) setExecPath(p string) { m.pathExec = p }

func (m ManagerBrew) ExecPath() string { return m.pathExec }

func (m ManagerBrew) PackageType() string { return Brew.String() }

func (m ManagerBrew) Install(name ...string) error {
	args := []string{"install"}

	_, err := m.cmd.Command(pathBrew, append(args, name...)...).Run()
	return err
}

func (m ManagerBrew) Remove(name ...string) error {
	args := []string{"uninstall"}

	_, err := m.cmd.Command(pathBrew, append(args, name...)...).Run()
	return err
}

func (m ManagerBrew) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerBrew) Update() error {
	_, err := m.cmd.Command(pathBrew, "update").Run()
	return err
}

func (m ManagerBrew) Upgrade() error {
	_, err := m.cmd.Command(pathBrew, "upgrade").Run()
	return err
}

//var msgWarning = []byte("Warning:")

func (m ManagerBrew) Clean() error {
	_, err := m.cmd.Command(pathBrew, "autoremove").Run()
	if err != nil {
		return err
	}

	// TODO: check exit code
	//return executil.RunToStdButErr(msgWarning, nil, pathBrew, "cleanup")
	_, err = m.cmd.Command(pathBrew, "cleanup").Run()
	return err
}

func (m ManagerBrew) AddRepo(alias string, url ...string) error {
	_, err := m.cmd.Command(pathBrew, "tap", url[0]).Run()
	return err
}

func (m ManagerBrew) RemoveRepo(r string) error {
	_, err := m.cmd.Command(pathBrew, "untap", r).Run()
	return err
}
