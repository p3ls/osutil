// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import (
	"github.com/tredoe/osutil/sh"
)

type deb struct{}

func (p deb) Install(name ...string) error {
	args := []string{"/usr/bin/apt-get", "install", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p deb) Remove(name ...string) error {
	args := []string{"/usr/bin/apt-get", "remove", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p deb) Purge(name ...string) error {
	args := []string{"/usr/bin/apt-get", "purge", "-y"}

	return sh.ExecToStd(nil, sudo, append(args, name...)...)
}

func (p deb) Update() error {
	return sh.ExecToStd(nil, sudo, "/usr/bin/apt-get", "update", "-qq")
}

func (p deb) Upgrade() error {
	return sh.ExecToStd(nil, sudo, "/usr/bin/apt-get", "upgrade", "-y")
}

func (p deb) Clean() error {
	err := sh.ExecToStd(nil, sudo, "/usr/bin/apt-get", "autoremove", "-y")
	if err != nil {
		return err
	}

	return sh.ExecToStd(nil, sudo, "/usr/bin/apt-get", "clean")
}
