// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

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
}

// NewManagerDnf returns the DNF package manager.
func NewManagerDnf() ManagerDnf {
	return ManagerDnf{pathExec: pathDnf}
}

func (m ManagerDnf) setExecPath(p string) { m.pathExec = p }

func (m ManagerDnf) ExecPath() string { return m.pathExec }

func (m ManagerDnf) PackageType() string { return Dnf.String() }

func (m ManagerDnf) Install(name ...string) error {
	args := []string{"install"}

	return executil.RunToStd(nil, pathDnf, append(args, name...)...)
}

func (m ManagerDnf) Remove(name ...string) error {
	args := []string{"remove"}

	return executil.RunToStd(nil, pathDnf, append(args, name...)...)
}

func (m ManagerDnf) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerDnf) Update() error {
	return executil.RunToStd(nil, pathDnf, "check-update")
}

func (m ManagerDnf) Upgrade() error {
	return executil.RunToStd(nil, pathDnf, "update")
}

func (m ManagerDnf) Clean() error {
	if err := executil.RunToStd(nil, pathDnf, "autoremove"); err != nil {
		return err
	}
	return executil.RunToStd(nil, pathDnf, "clean", "all")
}

// * * *

// ManagerYum is the interface to handle the package manager YUM of Linux systems
// based at Red Hat.
type ManagerYum struct {
	pathExec string
}

// NewManagerYum returns the YUM package manager.
func NewManagerYum() ManagerYum {
	return ManagerYum{pathExec: pathYum}
}

func (m ManagerYum) setExecPath(p string) { m.pathExec = p }

func (m ManagerYum) ExecPath() string { return m.pathExec }

func (m ManagerYum) PackageType() string { return Yum.String() }

func (m ManagerYum) Install(name ...string) error {
	args := []string{"install", "-y"}

	return executil.RunToStd(nil, pathYum, append(args, name...)...)
}

func (m ManagerYum) Remove(name ...string) error {
	args := []string{"remove", "-y"}

	return executil.RunToStd(nil, pathYum, append(args, name...)...)
}

func (m ManagerYum) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerYum) Update() error {
	return executil.RunToStd(nil, pathYum, "check-update")
}

func (m ManagerYum) Upgrade() error {
	return executil.RunToStd(nil, pathYum, "update")
}

func (m ManagerYum) Clean() error {
	return executil.RunToStd(nil, pathYum, "clean", "packages")
}

// * * *

// ManagerRpm is the interface to handle the package manager RPM of Linux systems
// based at Red Hat.
type ManagerRpm struct {
	pathExec string
}

// NewManagerRpm returns the RPM package manager.
func NewManagerRpm() ManagerRpm {
	return ManagerRpm{pathExec: pathRpm}
}

func (m ManagerRpm) setExecPath(p string) { m.pathExec = p }

func (m ManagerRpm) ExecPath() string { return m.pathExec }

func (m ManagerRpm) PackageType() string { return Rpm.String() }

func (m ManagerRpm) Install(name ...string) error {
	args := []string{"-i"}

	return executil.RunToStd(nil, pathRpm, append(args, name...)...)
}

func (m ManagerRpm) Remove(name ...string) error {
	args := []string{"-e"}

	return executil.RunToStd(nil, pathRpm, append(args, name...)...)
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
