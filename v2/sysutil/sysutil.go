// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

import (
	"io"
	"os"

	"github.com/tredoe/osutil/v2/executil"
)

var (
	cmd = executil.NewCommand("", "").
		Env(append([]string{"LANG=C"}, os.Environ()...))

	cmdWin = executil.NewCommand("", "").
		Env(os.Environ())
)

// CommandStdout sets the standard out at the commands of the package manager.
func CommandStdout(out io.Writer) {
	cmd.Stdout(out)
	cmdWin.Stdout(out)
}

// MustDisableColor returns true to indicate that the color should be disabled.
// It is useful to know when disable the color into a command line interface (cli).
//
// It checks the next environment variables:
//
//  + The NO_COLOR environment variable is set.
//  + The TERM environment variable has the value 'dumb'.
func MustDisableColor() bool {
	_, found := os.LookupEnv("NO_COLOR")
	if found {
		return true
	}

	if v := os.Getenv("TERM"); v == "dumb" {
		return true
	}

	return false
}
