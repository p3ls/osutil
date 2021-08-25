// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: Gentoo

package sysutil

import "github.com/tredoe/osutil/v2/executil"

const (
	fileEbuild = "emerge"
	pathEbuild = "/usr/bin/emerge"
)

// ManagerEbuild is the interface to handle the package manager of Linux systems based at Gentoo.
type ManagerEbuild struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerEbuild returns the Ebuild package manager.
func NewManagerEbuild() ManagerEbuild {
	return ManagerEbuild{
		pathExec: pathEbuild,
		cmd:      excmd.Command("", ""),
		//BadExitCodes([]int{1}),
	}
}

func (m ManagerEbuild) setExecPath(p string) { m.pathExec = p }

func (m ManagerEbuild) ExecPath() string { return m.pathExec }

func (m ManagerEbuild) PackageType() string { return Ebuild.String() }

func (m ManagerEbuild) Install(name ...string) error {
	_, err := m.cmd.Command(pathEbuild, name...).Run()
	return err
}

func (m ManagerEbuild) Remove(name ...string) error {
	args := []string{"--unmerge"}

	_, err := m.cmd.Command(pathEbuild, append(args, name...)...).Run()
	return err
}

func (m ManagerEbuild) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerEbuild) Update() error {
	_, err := m.cmd.Command(pathEbuild, "--sync").Run()
	return err
}

func (m ManagerEbuild) Upgrade() error {
	_, err := m.cmd.Command(pathEbuild, "--update", "--deep", "--with-bdeps=y", "--newuse @world").Run()
	return err
}

func (m ManagerEbuild) Clean() error {
	_, err := m.cmd.Command(pathEbuild, "--update", "--deep", "--newuse @world").Run()
	if err != nil {
		return err
	}
	_, err = m.cmd.Command(pathEbuild, "--depclean").Run()
	return err
}

func (m ManagerEbuild) ImportKey(alias, keyUrl string) error {
	return ErrRepo
}

func (m ManagerEbuild) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrRepo
}

func (m ManagerEbuild) RemoveKey(alias string) error {
	return ErrRepo
}

func (m ManagerEbuild) AddRepo(alias string, url ...string) error {
	return ErrRepo
}

func (m ManagerEbuild) RemoveRepo(r string) error {
	return ErrRepo
}
