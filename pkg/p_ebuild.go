// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/tredoe/osutil/sh"

type ebuild struct{}

func (p ebuild) Install(name ...string) error {
	return sh.ExecToStd(nil, "/usr/bin/emerge", name...)
}

func (p ebuild) Remove(name ...string) error {
	args := []string{"--unmerge"}

	return sh.ExecToStd(nil, "/usr/bin/emerge", append(args, name...)...)
}

func (p ebuild) Purge(name ...string) error {
	return p.Remove(name...)
}

func (p ebuild) Update() error {
	return sh.ExecToStd(nil, "/usr/bin/emerge", "--sync")
}

func (p ebuild) Upgrade() error {
	return sh.ExecToStd(nil, "/usr/bin/emerge", "--update", "--deep", "--with-bdeps=y", "--newuse @world")
}

func (p ebuild) Clean() error {
	err := sh.ExecToStd(nil, "/usr/bin/emerge", "--update", "--deep", "--newuse @world")
	if err != nil {
		return err
	}

	return sh.ExecToStd(nil, "/usr/bin/emerge", "--depclean")
}
