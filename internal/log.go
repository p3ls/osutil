// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
	"io"
	"log"
)

// Global loggers.
// By default, they do not write logs.
// It is setup from 'osutil.SetupLogger()'.
var (
	Log      = log.New(io.Discard, "", -1)
	LogShell = log.New(io.Discard, "", -1)
)
