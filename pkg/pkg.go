// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package pkg handles basic operations in the management of packages in
// FreeBSD, Linux and macOs operating systems.
package pkg

import (
	"fmt"
	"os/exec"
)

const sudo = "sudo"

// Manager is the common interface to handle different package systems.
type Manager interface {
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
	RPM
	Pacman
	Ebuild
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
	case RPM:
		return "RPM"
	case Pacman:
		return "Pacman"
	case Ebuild:
		return "Ebuild"
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

// New returns the interface to handle the package manager.
func New(pkg PackageType) Manager {
	switch pkg {
	// Linux
	case Deb:
		return new(ManagerDeb)
	case RPM:
		return new(ManagerRpm)
	case Pacman:
		return new(ManagerPacman)
	case Ebuild:
		return new(ManagerEbuild)
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

// execPackage is a list of executables of package managers.
var execPackage = [...]string{
	// Linux
	Deb:    "apt-get",
	RPM:    "yum",
	Pacman: "pacman",
	Ebuild: "emerge",
	Zypp:   "zypper",

	// BSD
	Brew: "brew",
	Pkg:  "pkg",
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
