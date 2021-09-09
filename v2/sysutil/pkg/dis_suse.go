// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Distro: SUSE

package pkg

import (
	"io"

	"github.com/tredoe/osutil/v2"
	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/internal"
	"github.com/tredoe/osutil/v2/sysutil"
)

const fileZypp = "zypper"

var pathZypp = "/usr/bin/zypper"

// ManagerZypp is the interface to handle the package manager of Linux systems based at SUSE.
type ManagerZypp struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerZypp returns the Zypp package manager.
// It checks if it is being run by an administrator.
func NewManagerZypp() (ManagerZypp, error) {
	if err := internal.CheckSuperuser(sysutil.Linux); err != nil {
		return ManagerZypp{}, err
	}

	return ManagerZypp{
		pathExec: pathZypp,
		cmd: cmd.Command("", "").
			// https://www.unix.com/man-page/suse/8/zypper/
			BadExitCodes([]int{1, 2, 3, 4, 5, 104}),
	}, nil
}

func (m ManagerZypp) setPathExec(p string) { m.pathExec = p }

func (m ManagerZypp) Cmd() *executil.Command { return m.cmd }

func (m ManagerZypp) PackageType() string { return Zypp.String() }

func (m ManagerZypp) PathExec() string { return m.pathExec }

func (m ManagerZypp) PreUsage() error { return nil }

func (m ManagerZypp) SetStdout(out io.Writer) { m.cmd.Stdout(out) }

// * * *

func (m ManagerZypp) Install(name ...string) error {
	osutil.Log.Print(taskInstall)
	args := append(
		[]string{
			"--non-interactive",
			"install", "--auto-agree-with-licenses", "-y",
		}, name...)

	_, err := m.cmd.Command(pathZypp, args...).Run()
	return err
}

func (m ManagerZypp) Remove(name ...string) error {
	osutil.Log.Print(taskRemove)
	args := append([]string{"remove", "-y"}, name...)

	_, err := m.cmd.Command(pathZypp, args...).Run()
	return err
}

func (m ManagerZypp) Purge(name ...string) error {
	osutil.Log.Print(taskPurge)
	return m.Remove(name...)
}

func (m ManagerZypp) UpdateIndex() error {
	osutil.Log.Print(taskUpdate)
	stderr, err := m.cmd.Command(pathZypp, "refresh").OutputStderr()

	return executil.CheckStderr(stderr, err)
}

func (m ManagerZypp) Update() error {
	osutil.Log.Print(taskUpgrade)
	_, err := m.cmd.Command(
		pathZypp, "up", "--auto-agree-with-licenses", "-y",
	).Run()
	return err
}

func (m ManagerZypp) Clean() error {
	osutil.Log.Print(taskClean)
	_, err := m.cmd.Command(pathZypp, "clean").Run()
	return err
}

// * * *

// https://opensuse-guide.org/repositories.php

func (m ManagerZypp) ImportKey(alias, keyUrl string) error {
	return ErrManagCmd
}

func (m ManagerZypp) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrManagCmd
}

func (m ManagerZypp) RemoveKey(alias string) error {
	return ErrManagCmd
}

func (m ManagerZypp) AddRepo(alias string, url ...string) error {
	osutil.Log.Print(taskAddRepo)
	_, err := m.cmd.Command(pathZypp, "addrepo", "-f", url[0], alias).Run()
	if err != nil {
		return err
	}

	return m.UpdateIndex()
}

func (m ManagerZypp) RemoveRepo(r string) error {
	osutil.Log.Print(taskRemoveRepo)
	if _, err := m.cmd.Command(pathZypp, "removerepo", r).Run(); err != nil {
		return err
	}

	return m.UpdateIndex()
}
