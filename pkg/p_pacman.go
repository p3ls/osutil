// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/executil"

const pathPacman = "/usr/bin/pacman"

// ManagerPacman is the interface to handle the package manager of Linux systems based at Arch.
type ManagerPacman struct{}

func (p ManagerPacman) Install(name ...string) error {
	args := []string{"-S", "--needed", "--noprogressbar"}

	return executil.ExecToStd(nil, pathPacman, append(args, name...)...)
}

func (p ManagerPacman) Remove(name ...string) error {
	args := []string{"-Rs"}

	return executil.ExecToStd(nil, pathPacman, append(args, name...)...)
}

func (p ManagerPacman) Purge(name ...string) error {
	args := []string{"-Rsn"}

	return executil.ExecToStd(nil, pathPacman, append(args, name...)...)
}

func (p ManagerPacman) Update() error {
	return executil.ExecToStd(nil, pathPacman, "-Syu", "--needed", "--noprogressbar")
}

func (p ManagerPacman) Upgrade() error {
	return executil.ExecToStd(nil, pathPacman, "-Syu")
}

func (p ManagerPacman) Clean() error {
	return executil.ExecToStd(nil, "/usr/bin/paccache", "-r")
}
