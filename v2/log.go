// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import (
	"log"

	"github.com/tredoe/osutil/v2/internal"
)

// SetupLogger setups the loggers used by some packages.
// 'log' is used by the packages 'edi', 'fileutil' and 'userutil'.
// 'logShell' is used by the package 'executil' and 'sysutil'.
func SetupLogger(log, logShell *log.Logger) {
	internal.Log = log
	internal.LogShell = logShell
}
