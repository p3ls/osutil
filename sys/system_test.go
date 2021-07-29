// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sys

import (
	"runtime"
	"testing"

	"github.com/tredoe/osutil/sh"
)

func TestExecWinshell(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}
	for _, v := range sh.ListWinShell {
		_, err := sh.ExecWinshell(v, false, `dir C:\`)
		if err != nil {
			t.Fatal(err)
		}
	}
}
