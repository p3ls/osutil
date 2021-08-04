// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

const sudo = "sudo"

// Manager is the common interface to handle different package systems.
type Manager interface {
	// ExecPath returns the executable path.
	ExecPath() string

	// Install installs packages.
	Install(name ...string) error

	// Remove removes packages.
	Remove(name ...string) error

	// Purge removes packages and its configuration files.
	Purge(name ...string) error

	// Update resynchronizes the package index files from their sources.
	Update() error

	// Upgrade upgrades all the packages on the system.
	Upgrade() error

	// Clean erases both packages downloaded and orphaned dependencies.
	Clean() error
}

// PackageType represents a package management system.
type PackageType int8

const (
	// Linux
	Deb PackageType = iota + 1
	Dnf
	Ebuild
	Pacman
	Rpm
	Yum
	Zypp

	// BSD
	Brew
	Pkg
)

func (pkg PackageType) String() string {
	switch pkg {
	// Linux
	case Deb:
		return "Deb"
	case Dnf:
		return "DNF"
	case Ebuild:
		return "Ebuild"
	case Pacman:
		return "Pacman"
	case Rpm:
		return "RPM"
	case Yum:
		return "YUM"
	case Zypp:
		return "ZYpp"

	// BSD
	case Brew:
		return "brew"
	case Pkg:
		return "pkg"
	}
	panic("unreachable")
}

// ExecPath returns the executable path.
func (pkg PackageType) ExecPath() string {
	switch pkg {
	// Linux
	case Deb:
		return pathDeb
	case Dnf:
		return pathDnf
	case Ebuild:
		return pathEbuild
	case Pacman:
		return pathPacman
	case Rpm:
		return pathRpm
	case Yum:
		return pathYum
	case Zypp:
		return pathZypp

	// BSD
	case Brew:
		return pathBrew
	case Pkg:
		return pathPkg
	}
	panic("unreachable")
}

// NewPkgFromType returns the interface to handle the package manager.
func NewPkgFromType(pkg PackageType) Manager {
	switch pkg {
	// Linux
	case Deb:
		return new(ManagerDeb)
	case Dnf:
		return new(ManagerDnf)
	case Ebuild:
		return new(ManagerEbuild)
	case Pacman:
		return new(ManagerPacman)
	case Rpm:
		return new(ManagerRpm)
	case Yum:
		return new(ManagerYum)
	case Zypp:
		return new(ManagerZypp)

	// BSD
	case Brew:
		return new(ManagerBrew)
	case Pkg:
		return new(ManagerPkg)
	}
	panic("unreachable")
}

// * * *

type pkgMngNotfoundError struct {
	Distro
}

func (e pkgMngNotfoundError) Error() string {
	return fmt.Sprintf(
		"package manager not found at Linux distro %s", e.Distro.String(),
	)
}

// NewPkgFromSystem returns the package manager used by a system.
func NewPkgFromSystem(sys System, dis Distro) (Manager, error) {
	switch sys {
	case Linux:
		pkg := newPkgFromDistro(dis)
		if len(pkg) == 1 {
			return pkg[0], nil
		}
		for _, v := range pkg {
			_, err := exec.LookPath(v.ExecPath())
			if err == nil {
				return v, nil
			}
		}
		return ManagerVoid{}, pkgMngNotfoundError{dis}

	case MacOS:
		return ManagerBrew{}, nil
	case FreeBSD:
		return ManagerPkg{}, nil

	default:
		panic("unimplemented")
	}
}

// newPkgFromDistro returns the package manager used by a Linux distro.
func newPkgFromDistro(dis Distro) []Manager {
	switch dis {
	case Debian:
		return []Manager{ManagerDeb{}}
	case Ubuntu:
		return []Manager{ManagerDeb{}}

	case CentOS:
		return []Manager{ManagerDnf{}, ManagerYum{}} // ManagerRpm{}
	case Fedora:
		return []Manager{ManagerDnf{}, ManagerYum{}} // ManagerRpm{}

	case OpenSUSE:
		return []Manager{ManagerZypp{}}

	case Arch:
		return []Manager{ManagerPacman{}}
	case Manjaro:
		return []Manager{ManagerPacman{}}

	default:
		panic("unimplemented")
	}
}

// * * *

// execPackage is a list of package managers executables.
var execPackage = [...]string{
	// Linux
	filepath.Base(Deb.ExecPath()),
	filepath.Base(Dnf.ExecPath()),
	filepath.Base(Yum.ExecPath()),
	filepath.Base(Zypp.ExecPath()),
	filepath.Base(Pacman.ExecPath()),
	filepath.Base(Ebuild.ExecPath()),
	//filepath.Base(Rpm.ExecPath()),

	// BSD
	filepath.Base(Brew.ExecPath()),
	filepath.Base(Pkg.ExecPath()),
}

// Detect tries to get the package system used in the system, looking for
// executables in directories "/usr/bin" and "/usr/local/bin".
func Detect() (PackageType, error) {
	for _, p := range []string{"/usr/bin/", "/usr/local/bin/"} {
		for k, v := range execPackage {
			_, err := exec.LookPath(p + v)
			if err == nil {
				return PackageType(k), nil
			}
		}
	}
	return -1, fmt.Errorf(
		"package manager not found in directories '/usr/bin' neither '/usr/local/bin'",
	)
}
