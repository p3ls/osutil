// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import "testing"

func TestDetect(t *testing.T) {
	if distro, err := DetectDistro(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(distro)
	}
}
