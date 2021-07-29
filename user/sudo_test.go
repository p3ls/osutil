// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"testing"

	"github.com/tredoe/osutil/system"
)

func TestSudo(t *testing.T) {
	if err := CheckSudo(system.SysUndefined); err != nil {
		t.Error(err)
	}
}
