// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package file

import (
	"errors"
	"fmt"
	"os"
)

// CheckDir checks if the path exists and if it is a directory.
func CheckDir(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		return err
	}

	if !info.Mode().IsDir() {
		return fmt.Errorf("expect a directory at \"%s\"", p)
	}
	return nil
}

// CreateDir creates a directory if it does not exist.
func CreateDir(dir string) error {
	existDir := false

	stat, err := os.Stat(dir)
	if err == nil {
		existDir = true

		if !stat.IsDir() {
			return fmt.Errorf("file with name %q exists", dir)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if !existDir {
		if err = os.MkdirAll(dir, 0774); err != nil {
			return err
		}
	}
	return nil
}
