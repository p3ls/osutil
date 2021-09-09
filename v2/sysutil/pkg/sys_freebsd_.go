// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// System: FreeBSD

package pkg

import (
	"io"

	"github.com/tredoe/osutil/v2"
	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/internal"
	"github.com/tredoe/osutil/v2/sysutil"
)

const (
	filePkg = "pkg"
	pathPkg = "/usr/sbin/pkg"
)

// ManagerPkg is the interface to handle the FreeBSD package manager,
// called 'package' or 'pkg'.
type ManagerPkg struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerPkg returns the Pkg package manager.
// It checks if it is being run by an administrator.
func NewManagerPkg() (ManagerPkg, error) {
	if err := internal.CheckSuperuser(sysutil.FreeBSD); err != nil {
		return ManagerPkg{}, err
	}

	return ManagerPkg{
		pathExec: pathPkg,
		cmd:      cmd.Command("", "").BadExitCodes([]int{1}),
	}, nil
}

func (m ManagerPkg) setPathExec(p string) { m.pathExec = p }

func (m ManagerPkg) Cmd() *executil.Command { return m.cmd }

func (m ManagerPkg) PackageType() string { return Pkg.String() }

func (m ManagerPkg) PathExec() string { return m.pathExec }

func (m ManagerPkg) PreUsage() error { return nil }

func (m ManagerPkg) SetStdout(out io.Writer) { m.cmd.Stdout(out) }

// * * *

func (m ManagerPkg) Install(name ...string) error {
	osutil.Log.Print(taskInstall)
	args := append([]string{"install", "-y"}, name...)

	_, err := m.cmd.Command(pathPkg, args...).Run()
	return err
}

func (m ManagerPkg) Remove(name ...string) error {
	osutil.Log.Print(taskRemove)
	args := append([]string{"delete", "-y"}, name...)

	_, err := m.cmd.Command(pathPkg, args...).Run()
	return err
}

func (m ManagerPkg) Purge(name ...string) error {
	osutil.Log.Print(taskPurge)
	return m.Remove(name...)
}

func (m ManagerPkg) UpdateIndex() error {
	osutil.Log.Print(taskUpdate)
	stderr, err := m.cmd.Command(pathPkg, "update").OutputStderr()

	return executil.CheckStderr(stderr, err)
}

func (m ManagerPkg) Update() error {
	osutil.Log.Print(taskUpgrade)
	_, err := m.cmd.Command(pathPkg, "upgrade", "-y").Run()
	return err
}

func (m ManagerPkg) Clean() error {
	osutil.Log.Print(taskClean)
	_, err := m.cmd.Command(pathPkg, "autoremove", "-y").Run()
	if err != nil {
		return err
	}
	_, err = m.cmd.Command(pathPkg, "clean", "-y").Run()
	return err
}

// * * *

func (m ManagerPkg) ImportKey(alias, keyUrl string) error {
	return ErrManagCmd
}

func (m ManagerPkg) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrManagCmd
}

func (m ManagerPkg) RemoveKey(alias string) error {
	return ErrManagCmd
}

func (m ManagerPkg) AddRepo(alias string, url ...string) error {
	return ErrManagCmd
}

func (m ManagerPkg) RemoveRepo(r string) error {
	return ErrManagCmd
}
