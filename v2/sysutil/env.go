// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

import "os"

// MustDisableColor returns true to indicate that the color should be disabled into a
// command line interface (cli).
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
