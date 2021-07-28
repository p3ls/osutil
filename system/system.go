// Copyright Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package system defines operating systems and detects the Linux distribution.
package system

import (
	"errors"
	"runtime"
)

var errSystem = errors.New("unsopported operating system")

// listSystem is the list of allowed operating systems.
var listSystem = [...]System{SysLinux, SysFreeBSD, SysMacOS, SysWindows}

// System represents an operating system.
type System uint8

// The operating systems.
const (
	_ System = iota
	SysLinux
	SysFreeBSD
	SysMacOS
	SysWindows
)

func (s System) String() string {
	switch s {
	case SysLinux:
		return "Linux"
	case SysFreeBSD:
		return "FreeBSD"
	case SysMacOS:
		return "macOS"
	case SysWindows:
		return "Windows"
	default:
		panic("unreachable")
	}
}

// SystemFromGOOS returns the system from 'GOOS', and the distribution at Linux systems.
func SystemFromGOOS() (sys System, dist Distro, err error) {
	switch runtime.GOOS {
	case "linux":
		sys = SysLinux

		if dist, err = DetectDistro(); err != nil {
			return 0, 0, err
		}
	case "freebsd":
		sys = SysFreeBSD
	case "darwin":
		sys = SysMacOS
	case "windows":
		sys = SysWindows

	default:
		return 0, 0, errSystem
	}

	return
}
