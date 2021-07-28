// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sh

import (
	"io"
	"log"
	"os"
)

const (
	path_       = "/sbin:/bin:/usr/sbin:/usr/bin"
	logFilename = "/.shutil.log" // at boot
)

var (
	logFile *os.File
	Log     = log.New(io.Discard, "", 0)
)

// To be used at functions related to 'ExecAsBash'.
var (
	env  []string
	home string // to expand symbol "~"
)

// Gets some environment variables.
func init() {
	env = os.Environ()
	home = os.Getenv("HOME")

	/*if path := os.Getenv("PATH"); path == "" {
		if err = os.Setenv("PATH", path_); err != nil {
			log.Print(err)
		}
	}*/
}

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

	env = []string{"PATH=" + path_} // from file boot
}

// CloseBootLogger closes the log file.
func CloseBootLogger() error {
	return logFile.Close()
}
