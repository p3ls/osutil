// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package executil

import (
	"errors"
	"fmt"
)

var (
	errEnvVar      = errors.New("the format of the variable has to be VAR=value")
	errNoCmdInPipe = errors.New("no command around of pipe")
)

// ErrProcKilled reports an error by a process killed.
var ErrProcKilled = errors.New("the process hasn't exited or was terminated by a signal")

// errorFromStderr returns the standard error like a Go error.
func errorFromStderr(e []byte) error {
	return fmt.Errorf("\n[stderr]\n%s", e)
}

// * * *

// extraCmdError reports an error due to the lack of an extra command.
type extraCmdError string

func (e extraCmdError) Error() string {
	return "command not added to " + string(e)
}
