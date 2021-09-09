// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package internal

import (
	"sync"

	"github.com/tredoe/osutil/v2/sysutil"
	"github.com/tredoe/osutil/v2/userutil"
)

var once sync.Once

// CheckSuperuser checks if it is being run by and administrator.
// Whether the system is 'SystemUndefined', it is detected.
func CheckSuperuser(sys sysutil.System) error {
	var err error
	once.Do(func() {
		err = userutil.MustBeSuperUser(sys)
	})

	return err
}
