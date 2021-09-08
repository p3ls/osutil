// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "errors"

var (
	ErrKeyUrl   = errors.New("the url has not a key file")
	ErrManagCmd = errors.New("unsupported command by the package manager")
)

type pkgTypeError string

func (e pkgTypeError) Error() string {
	return "invalid package type: " + string(e)
}
