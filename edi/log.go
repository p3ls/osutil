// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package edi

import (
	"io"
	"log"
)

// Logger is the global logger.
// By default, it does not write logs.
var Log = log.New(io.Discard, "", -1)
