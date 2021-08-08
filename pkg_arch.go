// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: Arch

package osutil

import "github.com/tredoe/osutil/executil"

const (
	filePacman = "pacman"
	pathPacman = "/usr/bin/pacman"
)

// ManagerPacman is the interface to handle the package manager of Linux systems based at Arch.
type ManagerPacman struct {
	pathExec string
}

// NewManagerPacman returns the Pacman package manager.
func NewManagerPacman() ManagerPacman {
	return ManagerPacman{pathExec: pathPacman}
}

func (m ManagerPacman) setExecPath(p string) { m.pathExec = p }

func (m ManagerPacman) ExecPath() string { return m.pathExec }

func (m ManagerPacman) PackageType() string { return Pacman.String() }

func (m ManagerPacman) Install(name ...string) error {
	args := []string{"-S", "--needed", "--noprogressbar"}

	return executil.RunToStd(nil, pathPacman, append(args, name...)...)
}

func (m ManagerPacman) Remove(name ...string) error {
	args := []string{"-Rs"}

	return executil.RunToStd(nil, pathPacman, append(args, name...)...)
}

func (m ManagerPacman) Purge(name ...string) error {
	args := []string{"-Rsn"}

	return executil.RunToStd(nil, pathPacman, append(args, name...)...)
}

func (m ManagerPacman) Update() error {
	return executil.RunToStd(nil, pathPacman, "-Syu", "--needed", "--noprogressbar")
}

func (m ManagerPacman) Upgrade() error {
	return executil.RunToStd(nil, pathPacman, "-Syu")
}

func (m ManagerPacman) Clean() error {
	return executil.RunToStd(nil, "/usr/bin/paccache", "-r")
}
