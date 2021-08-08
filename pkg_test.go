// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package osutil

import (
	"testing"
)

func TestPackager(t *testing.T) {
	pkg, err := DetectPkgManager()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Package type: %s", pkg.PackageType())

	if !testing.Verbose() {
		return
	}
	//cmd := "mtr-tiny"

	/*if err = pkg.Update(); err != nil {
		t.Fatal(err)
	}
	if err = pkg.Upgrade(); err != nil {
		t.Fatal(err)
	}*/

	/*if err = pkg.Install(cmd); err != nil {
		t.Errorf("\n%s", err)
	}
	if err = pkg.Remove(cmd); err != nil {
		t.Errorf("\n%s", err)
	}
	if err = pkg.Purge(cmd); err != nil {
		t.Errorf("\n%s", err)
	}

	if err = pkg.Clean(); err != nil {
		t.Errorf("\n%s", err)
	}*/
}
