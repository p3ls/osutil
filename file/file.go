// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package file

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// CheckFile checks if the path exists and if it is a file.
func CheckFile(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("expect a regular file at \"%s\"", p)
	}
	return nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherwise, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	return copyFileContents(src, dst)
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := out.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	return out.Sync()
}

// CreateFile creates a file.
func CreateFile(filePath string, content []byte) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := file.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	if _, err = file.Write(content); err != nil {
		return err
	}

	return nil
}

// Dload downloads a file.
func Dload(urlFile string, dst io.Writer) error {
	_, err := url.Parse(urlFile)
	if err != nil {
		return err
	}

	/*fileURL, err := url.Parse(urlFile)
	if err != nil {
		return err
	}*/
	/*path := fileURL.Path
	segments := strings.Split(path, "/")
	filename = segments[len(segments)-1]*/
	/*filename := filepath.Base(fileURL.Path)

	// Create blank file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := file.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()*/

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(urlFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	return err
}

// FilePathRelative returns the relative path of a file.
func FilePathRelative(dir, file string) string {
	relFile, _ := filepath.Rel(dir, file)
	if relFile == "" {
		relFile = file
	}

	return relFile
}

// Untar uncompresses a 'tar.gz' or 'tar' file.
func Untar(filename, dirDst string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := file.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	var tr *tar.Reader

	// TODO: maybe used gzip at '.tar' files
	if strings.HasSuffix(filename, ".tar.gz") {
		uncompressedStream, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		tr = tar.NewReader(uncompressedStream)
	} else if strings.HasSuffix(filename, ".tar") {
		tr = tar.NewReader(file)
	} else {
		return errNotTar
	}

	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// If the header is nil, just skip it (not sure how this happens)
		if header == nil {
			continue
		}

		// The target location where the dir/file should be created.
		target := filepath.Join(dirDst, header.Name)

		// Check the file type
		switch header.Typeflag {
		default: // TODO: remove?
			return fmt.Errorf(
				"Untar: unknown type: %s in %s",
				string(header.Typeflag),
				header.Name,
			)

		// If it's a dir and it doesn't exist, create it.
		case tar.TypeDir:
			if _, err = os.Stat(target); err != nil {
				if err = os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// If it's a file, create it.
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// Copy over contents
			if _, err = io.Copy(f, tr); err != nil {
				return err
			}

			if err2 := f.Close(); err2 != nil && err == nil {
				err = err2
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteTempFile writes bytes to a temporary file and returns its name.
func WriteTempFile(b []byte, name, ext string) (filename string, err error) {
	tmpfile, err := os.CreateTemp("", fmt.Sprintf("%s_*%s", name, ext))
	if err != nil {
		return "", err
	}
	filename = tmpfile.Name()

	defer func() {
		if err2 := tmpfile.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()

	if _, err = tmpfile.Write(b); err != nil {
		return "", err
	}
	if err = tmpfile.Sync(); err != nil {
		return "", err
	}

	return
}

// == Errors
//

var errNotTar = errors.New("not a tar file")
