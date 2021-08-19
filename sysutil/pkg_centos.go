// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

import "github.com/tredoe/osutil/executil"

const (
	fileDnf = "dnf" // Preferable to YUM
	pathDnf = "/usr/bin/dnf"

	fileYum = "yum"
	pathYum = "/usr/bin/yum"

	// RPM is used to install/uninstall local packages.
	fileRpm = "rpm"
	pathRpm = "/usr/bin/rpm"
)

// ManagerDnf is the interface to handle the package manager DNG of Linux systems
// based at Red Hat.
type ManagerDnf struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerDnf returns the DNF package manager.
func NewManagerDnf() ManagerDnf {
	return ManagerDnf{
		pathExec: pathDnf,
		cmd: excmd.Command("", "").
			// https://dnf.readthedocs.io/en/latest/command_ref.html
			BadExitCodes([]int{1, 3, 200}),
	}
}

func (m ManagerDnf) setExecPath(p string) { m.pathExec = p }

func (m ManagerDnf) ExecPath() string { return m.pathExec }

func (m ManagerDnf) PackageType() string { return Dnf.String() }

func (m ManagerDnf) Install(name ...string) error {
	args := []string{pathDnf, "install", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerDnf) Remove(name ...string) error {
	args := []string{pathDnf, "remove", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerDnf) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerDnf) Update() error {
	// check-update does not update else it checks the updating.
	return nil

	// check-update returns exit value of 100 if there are packages available for an update.
	// Also returns a list of the packages to be updated in list format.
	// Returns 0 if no packages are available for update.
	// Returns 1 if an error occurred.
	/*err := m.cmd.Command(sudo, pathDnf, "check-update")
	if err != nil {
		// Check the exit code
	}
	return err*/
}

func (m ManagerDnf) Upgrade() error {
	_, err := m.cmd.Command(sudo, pathDnf, "update", "-y").Run()
	return err
}

func (m ManagerDnf) Clean() error {
	_, err := m.cmd.Command(sudo, pathDnf, "autoremove", "-y").Run()
	if err != nil {
		return err
	}
	_, err = m.cmd.Command(sudo, pathDnf, "clean", "all").Run()
	return err
}

// * * *

// ManagerYum is the interface to handle the package manager YUM of Linux systems
// based at Red Hat.
type ManagerYum struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerYum returns the YUM package manager.
func NewManagerYum() ManagerYum {
	return ManagerYum{
		pathExec: pathYum,
		cmd: excmd.Command("", "").
			BadExitCodes([]int{1, 2, 3, 16}),
	}
}

func (m ManagerYum) setExecPath(p string) { m.pathExec = p }

func (m ManagerYum) ExecPath() string { return m.pathExec }

func (m ManagerYum) PackageType() string { return Yum.String() }

func (m ManagerYum) Install(name ...string) error {
	args := []string{pathYum, "install", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerYum) Remove(name ...string) error {
	args := []string{pathYum, "remove", "-y"}

	_, err := m.cmd.Command(sudo, append(args, name...)...).Run()
	return err
}

func (m ManagerYum) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerYum) Update() error {
	// check-update does not update else it checks the updating.
	return nil
}

func (m ManagerYum) Upgrade() error {
	_, err := m.cmd.Command(sudo, pathYum, "update", "-y").Run()
	return err
}

func (m ManagerYum) Clean() error {
	_, err := m.cmd.Command(sudo, pathYum, "clean", "packages").Run()
	return err
}

// * * *

// ManagerRpm is the interface to handle the package manager RPM of Linux systems
// based at Red Hat.
type ManagerRpm struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerRpm returns the RPM package manager.
func NewManagerRpm() ManagerRpm {
	return ManagerRpm{
		pathExec: pathRpm,
		cmd:      excmd.Command("", ""),
		//BadExitCodes([]int{}),
	}
}

func (m ManagerRpm) setExecPath(p string) { m.pathExec = p }

func (m ManagerRpm) ExecPath() string { return m.pathExec }

func (m ManagerRpm) PackageType() string { return Rpm.String() }

func (m ManagerRpm) Install(name ...string) error {
	args := []string{"-i"}

	_, err := m.cmd.Command(pathRpm, append(args, name...)...).Run()
	return err
}

func (m ManagerRpm) Remove(name ...string) error {
	args := []string{"-e"}

	_, err := m.cmd.Command(pathRpm, append(args, name...)...).Run()
	return err
}

func (m ManagerRpm) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerRpm) Update() error {
	return nil
}

func (m ManagerRpm) Upgrade() error {
	return nil
}

func (m ManagerRpm) Clean() error {
	return nil
}
