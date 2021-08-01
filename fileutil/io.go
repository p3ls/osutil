// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// CopyFile copies a file from 'src' to 'dst'. If 'src' and 'dst' files exist, and are
// the same, then return success. Otherwise, copy the file contents from 'src' to 'dst'.
//
// The file will be created if it does not already exist. If the destination file exists,
// all it's contents will be replaced by the contents of the source file.
func CopyFile(src, dst string) error {
	sInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sInfo.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("non-regular source file %s (%q)", sInfo.Name(), sInfo.Mode().String())
	}

	dInfo, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dInfo.Mode().IsRegular()) {
			return fmt.Errorf("non-regular destination file %s (%q)", dInfo.Name(), dInfo.Mode().String())
		}
		if os.SameFile(sInfo, dInfo) {
			return nil
		}
	}

	// Open original file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create new file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer func() {
		if err2 := dstFile.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	// Copy the bytes to destination from source
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	// Commit the file contents. Flushes memory to disk.
	if err = dstFile.Sync(); err != nil {
		return err
	}

	Log.Printf("File %q copied at %q", src, dst)
	return nil
}

// Create creates a new file with b bytes.
func Create(filename string, b []byte) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	err2 := file.Close()
	if err2 != nil && err == nil {
		err = err2
	}
	if err != nil {
		return err
	}

	Log.Printf("File %q created", filename)
	return nil
}

// Overwrite truncates the named file to zero and writes len(b) bytes.
func Overwrite(filename string, b []byte) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	err2 := file.Close()
	if err2 != nil && err == nil {
		err = err2
	}
	if err != nil {
		return err
	}

	Log.Printf("File %q overwritted", filename)
	return nil
}

// PrefixTemp is the prefix to add to temporary files.
const PrefixTemp = "tmp-"

// TempFile creates a temporary file from the source file into the default directory
// temporary files (see os.TempDir), whose name begins with the prefix.
// If prefix is the empty string, uses the default value PrefixTemp.
// Returns the temporary file name.
func TempFile(src, prefix string) (tmpFile string, err error) {
	if prefix == "" {
		prefix = PrefixTemp
	}

	fsrc, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer fsrc.Close()

	fdst, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	defer func() {
		if err2 := fdst.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	if _, err = io.Copy(fdst, fsrc); err != nil {
		return "", err
	}
	if err = fdst.Sync(); err != nil {
		return "", err
	}

	Log.Printf("File %q copied at %q", src, fdst.Name())
	return fdst.Name(), nil
}
