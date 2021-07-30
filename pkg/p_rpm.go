// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/sh"

const pathYum = "/usr/bin/yum"

// ManagerRpm is the interface to handle the package manager of Linux systems based at Red Hat.
type ManagerRpm struct{}

func (p ManagerRpm) Install(name ...string) error {
	args := []string{"install"}

	return sh.ExecToStd(nil, pathYum, append(args, name...)...)
}

func (p ManagerRpm) Remove(name ...string) error {
	args := []string{"remove"}

	return sh.ExecToStd(nil, pathYum, append(args, name...)...)
}

func (p ManagerRpm) Purge(name ...string) error {
	return p.Remove(name...)
}

func (p ManagerRpm) Update() error {
	return sh.ExecToStd(nil, pathYum, "update")
}

func (p ManagerRpm) Upgrade() error {
	return sh.ExecToStd(nil, pathYum, "update")
}

func (p ManagerRpm) Clean() error {
	return sh.ExecToStd(nil, pathYum, "clean", "packages")
}
