// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

import (
	"os"

	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/internal"
)

var excmd = executil.NewCommand("", "").
	Stdout(internal.LogShell.Writer()).
	Env(append([]string{"LANG=C"}, os.Environ()...))
