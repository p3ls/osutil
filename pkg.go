// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const sudo = "sudo"

// PkgManager is the common interface to handle different package systems.
type PkgManager interface {
	// setExecPath sets the executable path.
	setExecPath(p string)

	// ExecPath returns the executable path.
	ExecPath() string

	// PackageType returns the package type.
	PackageType() string

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

// NewPkgTypeFromstr returns a package management system from the string.
func NewPkgTypeFromstr(s string) (PackageType, error) {
	switch strings.ToLower(s) {
	case fileDeb:
		return Deb, nil
	case fileDnf:
		return Dnf, nil
	case fileEbuild:
		return Ebuild, nil
	case filePacman:
		return Pacman, nil
	case fileRpm:
		return Rpm, nil
	case fileYum:
		return Yum, nil
	case fileZypp:
		return Zypp, nil

	case fileBrew:
		return Brew, nil
	case filePkg:
		return Pkg, nil

	default:
		return -1, pkgTypeError(s)
	}
}

// * * *

// NewPkgManagerFromType returns the interface to handle the package manager.
func NewPkgManagerFromType(pkg PackageType) PkgManager {
	switch pkg {
	// Linux
	case Deb:
		return NewManagerDeb()
	case Dnf:
		return NewManagerDnf()
	case Ebuild:
		return NewManagerEbuild()
	case Pacman:
		return NewManagerPacman()
	case Rpm:
		return NewManagerRpm()
	case Yum:
		return NewManagerYum()
	case Zypp:
		return NewManagerZypp()

	// BSD
	case Brew:
		return NewManagerBrew()
	case Pkg:
		return NewManagerPkg()
	}
	panic("unreachable")
}

// * * *

// NewPkgManagerFromSystem returns the package manager used by a system.
func NewPkgManagerFromSystem(sys System, dis Distro) (PkgManager, error) {
	switch sys {
	case Linux:
		return newPkgManagerFromDistro(dis)
		/*pkg, err := newPkgManagerFromDistro(dis)
		if err != nil {
			return ManagerVoid{}, err
		}

		if len(pkg) == 1 {
			return pkg[0], nil
		}
		for _, v := range pkg {
			_, err := exec.LookPath(v.ExecPath())
			if err == nil {
				return v, nil
			}
		}
		return ManagerVoid{}, pkgMngNotfoundError{dis}*/

	case MacOS:
		return NewManagerBrew(), nil
	case FreeBSD:
		return NewManagerPkg(), nil

	default:
		panic("unimplemented")
	}
}

// newPkgManagerFromDistro returns the package manager used by a Linux distro.
func newPkgManagerFromDistro(dis Distro) (PkgManager, error) {
	switch dis {
	case Debian, Ubuntu:
		return NewManagerDeb(), nil

	case OpenSUSE:
		return NewManagerZypp(), nil

	case Arch, Manjaro:
		return NewManagerPacman(), nil

	// DNF is the default package manager of Fedora 22, CentOS8, and RHEL8.
	case CentOS, Fedora:
		verStr, err := DetectSystemVer(Linux)
		if err != nil {
			return ManagerVoid{}, err
		}
		ver, err := strconv.Atoi(verStr)
		if err != nil {
			return ManagerVoid{}, err
		}

		useDnf := true
		switch dis {
		case CentOS:
			if ver < 8 {
				useDnf = false
			}
		case Fedora:
			if ver < 22 {
				useDnf = false
			}
		}

		if useDnf {
			return NewManagerDnf(), nil
		}
		return NewManagerYum(), nil

	default:
		panic("unimplemented")
	}
}

// * * *

// execPackage is a list of package managers executables.
var execPackage = [...]string{
	// Linux
	fileDeb,
	fileDnf,
	fileYum,
	fileZypp,
	filePacman,
	fileEbuild,
	fileRpm,

	// BSD
	fileBrew,
	filePkg,
}

// DetectPkgManager tries to get the package manager used in the system, looking for
// executables at directories in $PATH.
func DetectPkgManager() (PkgManager, error) {
	for _, p := range execPackage {
		pathExec, err := exec.LookPath(p)
		if err == nil {
			pkg, err := NewPkgTypeFromstr(p)
			if err != nil {
				return ManagerVoid{}, err
			}
			mng := NewPkgManagerFromType(pkg)

			if mng.ExecPath() != pathExec {
				mng.setExecPath(pathExec)
			}

			return mng, nil
		}
	}

	return ManagerVoid{}, fmt.Errorf("package manager not found in $PATH")
}

// == Errors

type pkgMngNotfoundError struct {
	Distro
}

func (e pkgMngNotfoundError) Error() string {
	return fmt.Sprintf(
		"package manager not found at Linux distro %s", e.Distro.String(),
	)
}

type pkgTypeError string

func (e pkgTypeError) Error() string {
	return "invalid package type: " + string(e)
}