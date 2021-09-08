// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

import (
	"fmt"
)

type pkgManagNotfoundError struct {
	Distro
}

func (e pkgManagNotfoundError) Error() string {
	return fmt.Sprintf(
		"package manager not found at Linux distro %s", e.Distro.String(),
	)
}
