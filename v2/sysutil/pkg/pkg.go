// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package pkg handles the basic operations in the management of packages at
FreeBSD, Linux and macOS operating systems.

By default, the output of the commands run by the package managers are not printed.
To set the standard output, to use the method 'SetStdout()'.

NOTE: the next package management systems are untested:

 + Packman (Arch)
 + ebuild  (Gentoo)
*/
package pkg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/sysutil"
)

// sudo is the path by default at Linux systems.
const sudo = "/usr/bin/sudo"

const (
	taskInstall             = "Installing ..."
	taskRemove              = "Removing ..."
	taskPurge               = "Purging ..."
	taskUpdate              = "Updating repositories ..."
	taskUpgrade             = "Upgrading packages ..."
	taskClean               = "Cleaning ..."
	taskImportKey           = "Importing key ..."
	taskImportKeyFromServer = "Importing key from server ..."
	taskRemoveKey           = "Removing key ..."
	taskAddRepo             = "Adding repository ..."
	taskRemoveRepo          = "Removing repository ..."
)

var (
	cmd = executil.NewCommand("", "").
		Env(append([]string{"LANG=C"}, os.Environ()...))

	cmdWin = executil.NewCommand("", "").
		Env(os.Environ())
)

// Manager is the common interface to handle different package systems.
type Manager interface {
	// setPathExec sets the path of the executable.
	setPathExec(p string)

	// Cmd returns the command configured for the package manager.
	Cmd() *executil.Command

	// PackageType returns the package type.
	PackageType() string

	// PathExec returns the path of the executable.
	PathExec() string

	// PreUsage creates the files or installs the programs which are necessary for that
	// all methods can be used.
	PreUsage() error

	// SetStdout sets the standard out for the commands of the package manager.
	SetStdout(out io.Writer)

	// * * *

	// Install installs packages.
	Install(name ...string) error

	// Remove removes packages.
	Remove(name ...string) error

	// Purge removes packages and its configuration files (if the package system does it).
	Purge(name ...string) error

	// UpdateIndex resynchronizes the package index files from their sources.
	UpdateIndex() error

	// Update installs the newest versions of all packages currently installed.
	Update() error

	// Clean erases both packages downloaded and orphaned dependencies.
	Clean() error

	// * * *

	// ImportKey downloads the OpenPGP key and add it to the system.
	ImportKey(alias, keyUrl string) error

	// ImportKeyFromServer imports OpenPGP key directly from a keyserver.
	// Whether 'keyServer' is empty, then it uses a value by default.
	ImportKeyFromServer(alias, keyServer, key string) error

	// RemoveKey removes the OpenPGP key.
	RemoveKey(alias string) error

	// AddRepo adds a repository.
	AddRepo(alias string, url ...string) error

	// RemoveRepo removes a repository.
	RemoveRepo(string) error
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

	// Windows
	Choco
	Winget
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

	case Choco:
		return "Chocolatey"
	case Winget:
		return "winget"
	}
	panic("unreachable")
}

// NewTypeFromStr returns a package management system from the string.
func NewTypeFromStr(s string) (PackageType, error) {
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

	case fileChoco:
		return Choco, nil
	case fileWinget:
		return Winget, nil

	default:
		return -1, pkgTypeError(s)
	}
}

// * * *

// NewManagerFromType returns the interface to handle the package manager.
func NewManagerFromType(pkg PackageType) Manager {
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

	// Windows
	case Choco:
		return NewManagerChoco()
	case Winget:
		return NewManagerWinget()
	}
	panic("unreachable")
}

// * * *

// NewManagerFromSystem returns the package manager used by a system.
func NewManagerFromSystem(sys sysutil.System, dis sysutil.Distro) (Manager, error) {
	switch sys {
	case sysutil.Linux:
		return NewManagerFromDistro(dis)
		/*pkg, err := NewManagerFromDistro(dis)
		if err != nil {
			return ManagerVoid{}, err
		}

		if len(pkg) == 1 {
			return pkg[0], nil
		}
		for _, v := range pkg {
			_, err := exec.LookPath(v.PathExec())
			if err == nil {
				return v, nil
			}
		}
		return ManagerVoid{}, pkgManagNotfoundError{dis}*/

	case sysutil.MacOS:
		return NewManagerBrew(), nil
	case sysutil.FreeBSD:
		return NewManagerPkg(), nil
	case sysutil.Windows:
		// TODO: in the future, to use winget
		return NewManagerChoco(), nil

	default:
		panic("unimplemented")
	}
}

// NewManagerFromDistro returns the package manager used by a Linux distro.
func NewManagerFromDistro(dis sysutil.Distro) (Manager, error) {
	switch dis {
	case sysutil.Debian, sysutil.Ubuntu:
		return NewManagerDeb(), nil

	case sysutil.OpenSUSE:
		return NewManagerZypp(), nil

	case sysutil.Arch, sysutil.Manjaro:
		return NewManagerPacman(), nil

	// DNF is the default package manager of Fedora 22, CentOS8, and RHEL8.
	case sysutil.CentOS, sysutil.Fedora:
		verStr, err := sysutil.DetectSystemVer(sysutil.Linux)
		if err != nil {
			return ManagerVoid{}, err
		}
		ver, err := strconv.Atoi(verStr)
		if err != nil {
			return ManagerVoid{}, err
		}

		useDnf := true
		switch dis {
		case sysutil.CentOS:
			if ver < 8 {
				useDnf = false
			}
		case sysutil.Fedora:
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

// execPkgLinux is a list of package managers executables for Linux.
var execPkgLinux = []string{
	fileDeb,
	fileDnf,
	fileYum,
	fileZypp,
	filePacman,
	fileEbuild,
	fileRpm,
}

// execPkgFreebsd is a list of package managers executables for FreeBSD.
var execPkgFreebsd = []string{
	fileBrew,
}

// execPkgMacos is a list of package managers executables for MacOS.
var execPkgMacos = []string{
	filePkg,
}

// execPkgWindows is a list of package managers executables for Windows.
var execPkgWindows = []string{
	fileChoco,
	fileWinget,
}

// DetectManager tries to get the package manager used in the system, looking for
// executables at directories in $PATH.
func DetectManager(sys sysutil.System) (Manager, error) {
	var execPkg []string
	switch sys {
	case sysutil.Linux:
		execPkg = execPkgLinux
	case sysutil.FreeBSD:
		execPkg = execPkgFreebsd
	case sysutil.MacOS:
		execPkg = execPkgMacos
	case sysutil.Windows:
		execPkg = execPkgWindows

	default:
		panic("unimplemented: " + sys.String())
	}

	for _, p := range execPkg {
		pathExec, err := exec.LookPath(p)
		if err == nil {
			pkg, err := NewTypeFromStr(p)
			if err != nil {
				return ManagerVoid{}, err
			}
			mng := NewManagerFromType(pkg)

			if mng.PathExec() != pathExec {
				mng.setPathExec(pathExec)
			}

			return mng, nil
		}
	}

	return ManagerVoid{}, fmt.Errorf("package manager not found in $PATH")
}
