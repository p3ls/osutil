// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package executil

import (
	"io"
	"log"
	"os"
)

const logFilename = "/.shutil.log" // at boot

var (
	logFile *os.File
	Log     = log.New(io.Discard, "", 0)
)

// StartBootLogger initializes the log file to be used during the boot.
func StartBootLogger() {
	var err error

	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	logFile, err = os.OpenFile(logFilename, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		log.Print(err)
	} else {
		Log = log.New(logFile, "", log.Lshortfile)
	}
}

// CloseBootLogger closes the log file.
func CloseBootLogger() error {
	return logFile.Close()
}
