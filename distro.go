// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import (
	"os"

	"github.com/tredoe/osutil/config/shconf"
)

// Distro represents a distribution of Linux system.
type Distro int

// Most used Linux distributions.
const (
	DistroUnknown Distro = iota

	Debian
	Ubuntu

	Fedora
	CentOS

	OpenSUSE

	Arch
	Manjaro
)

var distroNames = [...]string{
	DistroUnknown: "unknown distribution",

	Debian: "Debian",
	Ubuntu: "Ubuntu",

	Fedora: "Fedora",
	CentOS: "CentOS",

	OpenSUSE: "openSUSE",

	Arch:    "Arch",
	Manjaro: "Manjaro",
}

func (s Distro) String() string { return distroNames[s] }

// * * *

var idToDistro = map[string]Distro{
	"debian": Debian,
	"ubuntu": Ubuntu,

	"centos": CentOS,
	"fedora": Fedora,

	"opensuse-leap":       OpenSUSE,
	"opensuse-tumbleweed": OpenSUSE,

	"arch":    Arch,
	"manjaro": Manjaro, // based on Arch
}

// DetectDistro returns the Linux distribution.
func DetectDistro() (Distro, error) {
	var id string
	var err error

	if _, err = os.Stat("/etc/os-release"); !os.IsNotExist(err) {
		cfg, err := shconf.ParseFile("/etc/os-release")
		if err != nil {
			return 0, err
		}
		if id, err = cfg.Get("ID"); err != nil {
			return 0, err
		}

		if v, found := idToDistro[id]; found {
			return v, nil
		}
	}

	return DistroUnknown, nil
}
