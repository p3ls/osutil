// Copyright Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package sys defines operating systems and detects the Linux distribution.
package sys

import (
	"errors"
	"runtime"

	"github.com/tredoe/osutil/pkg"
)

var errSystem = errors.New("unsopported operating system")

// ListSystem is the list of allowed operating systems.
var ListSystem = [...]System{FreeBSD, Linux, MacOS, Windows}

// System represents an operating system.
type System uint8

// The operating systems.
const (
	SystemUndefined System = iota
	Linux
	FreeBSD
	MacOS
	Windows
)

func (s System) String() string {
	switch s {
	case Linux:
		return "Linux"
	case FreeBSD:
		return "FreeBSD"
	case MacOS:
		return "macOS"
	case Windows:
		return "Windows"
	default:
		panic("unreachable")
	}
}

// Manager returns the package manager.
func (s System) Manager(dist Distro) pkg.Manager {
	switch s {
	case Linux:
		return Distro.Manager(dist)
	case MacOS:
		return pkg.ManagerBrew{}
	case FreeBSD:
		return pkg.ManagerPkg{}

	default:
		panic("unimplemented")
	}
}

// * * *

// SystemFromGOOS returns the system from 'GOOS', and the distribution at Linux systems.
func SystemFromGOOS() (sys System, dist Distro, err error) {
	switch runtime.GOOS {
	case "linux":
		sys = Linux

		if dist, err = DetectDistro(); err != nil {
			return 0, 0, err
		}
	case "freebsd":
		sys = FreeBSD
	case "darwin":
		sys = MacOS
	case "windows":
		sys = Windows

	default:
		return 0, 0, errSystem
	}

	return
}
