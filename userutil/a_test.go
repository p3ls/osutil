// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tredoe/fileutil"
	"github.com/tredoe/osutil/sys"
)

const (
	USER     = "u_foo"
	USER2    = "u_foo2"
	SYS_USER = "usys_bar"

	GROUP     = "g_foo"
	SYS_GROUP = "gsys_bar"
)

const (
	prefixTemp = "test_osutil-"

	FILE_USER    = prefixTemp + "passwd"
	FILE_GROUP   = prefixTemp + "group"
	FILE_SHADOW  = prefixTemp + "shadow"
	FILE_GSHADOW = prefixTemp + "gshadow"
)

var MEMBERS = []string{USER, SYS_USER}

// Stores the ids at creating the groups.
var GID, SYS_GID int

// == Copy the system files before of be edited.

func init() {
	err := MustBeSuperUser(sys.SystemUndefined)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if fileUser, err = fileutil.CopytoTemp(fileUser, FILE_USER); err != nil {
		goto _error
	}
	if fileGroup, err = fileutil.CopytoTemp(fileGroup, FILE_GROUP); err != nil {
		goto _error
	}
	if fileShadow, err = fileutil.CopytoTemp(fileShadow, FILE_SHADOW); err != nil {
		goto _error
	}
	if fileGShadow, err = fileutil.CopytoTemp(fileGShadow, FILE_GSHADOW); err != nil {
		goto _error
	}

	return

_error:
	removeTempFiles()
	log.Fatalf("%s", err)
}

func removeTempFiles() {
	files, _ := filepath.Glob(filepath.Join(os.TempDir(), prefixTemp+"*"))

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Printf("%s", err)
		}
	}
}
