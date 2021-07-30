// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sh

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Exec executes a command. Returns both standard output and error.
func Exec(cmd string, args ...string) (stdout, stderr []byte, err error) {
	return ExecWithTime(0, cmd, args...)
}

// ExecWithTime executes a command waiting to finish the command before of kill it ('timeKillCmd'),
// or waits without kill it when the duration is lesser or equal to zero).
// Logs the command and returns both standard output and error.
func ExecWithTime(timeKillCmd time.Duration, cmd string, args ...string,
) (stdout, stderr []byte, err error) {
	Log.Printf("%s %s", cmd, strings.Join(args, " "))

	var outPipe, errPipe io.ReadCloser
	var ctx context.Context
	var cancel context.CancelFunc

	if timeKillCmd > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeKillCmd)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	c := exec.CommandContext(ctx, cmd, args...)

	// Using 'bytes.Buffer' to get stdout and stderr gives error:
	// https://github.com/golang/go/issues/23019
	//
	// var bufOut, bufStderr bytes.Buffer
	//c.Stdout = &bufOut
	//c.Stderr = &bufStderr

	if outPipe, err = c.StdoutPipe(); err != nil {
		return
	}
	if errPipe, err = c.StderrPipe(); err != nil {
		return
	}
	if err = c.Start(); err != nil {
		goto _checkErr
	}

	// Std out
	go func() {
		var bufOut bytes.Buffer
		buf := bufio.NewReader(outPipe)
		for {
			line, err2 := buf.ReadBytes('\n')
			if len(line) > 0 {
				bufOut.Write(line)
			}
			if err2 != nil {
				stdout = bufOut.Bytes()
				if err2 != io.EOF && !errors.Is(err2, os.ErrClosed) && err == nil {
					err = err2
				}

				return
			}
		}
	}()
	// Std error
	go func() {
		var bufStderr bytes.Buffer
		buf := bufio.NewReader(errPipe)
		for {
			line, err2 := buf.ReadBytes('\n')
			if len(line) > 0 {
				bufStderr.Write(line)
			}
			if err2 != nil {
				stderr = bufStderr.Bytes()
				if err2 != io.EOF && !errors.Is(err2, os.ErrClosed) && err == nil {
					err = err2
				}

				return
			}
		}
	}()

	if err = c.Wait(); err == nil {
		return
	}

_checkErr:
	switch errType := err.(type) {
	case *exec.ExitError:
		exitCode := errType.ExitCode()

		if exitCode == -1 {
			err = ErrProcKilled
		}
	}

	return
}

// ExecToStd executes a command setting both standard output and error.
// Logs the command.
func ExecToStd(extraEnv []string, cmd string, args ...string) error {
	Log.Printf("%s %s", cmd, strings.Join(args, " "))

	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	c.Env = append(c.Env, "LANG=C")

	if extraEnv != nil {
		//fmt.Printf("%s", strings.Join(extraEnv, " "))
		c.Env = append(c.Env, extraEnv...)
	}

	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// ExecNoStdErr executes a command setting only standard output.
// Logs the command.
func ExecNoStdErr(extraEnv []string, cmd string, args ...string) error {
	Log.Printf("%s %s", cmd, strings.Join(args, " "))

	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	//c.Stderr = os.Stderr
	c.Env = os.Environ()
	c.Env = append(c.Env, "LANG=C")

	if extraEnv != nil {
		//fmt.Printf("%s", strings.Join(extraEnv, " "))
		c.Env = append(c.Env, extraEnv...)
	}

	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// ExecToStdButErr executes a command setting the standard output.
// Logs the command.
//
// checkStderr (if any) checks if it is found in the standard error to know whether the standard
// error is not really an error.
func ExecToStdButErr(checkStderr []byte, extraEnv []string, cmd string, args ...string) error {
	Log.Printf("%s %s", cmd, strings.Join(args, " "))

	var bufStderr bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = &bufStderr
	c.Env = os.Environ()
	c.Env = append(c.Env, "LANG=C")

	if extraEnv != nil {
		//fmt.Printf("%s", strings.Join(extraEnv, " "))
		c.Env = append(c.Env, extraEnv...)
	}

	err := c.Start()
	if err != nil {
		return err
	}
	if err = c.Wait(); err != nil {
		return err
	}

	stderr := bufStderr.Bytes()

	if len(stderr) == 0 {
		return nil
	}
	if checkStderr == nil {
		return errorFromStderr(stderr)
	}

	s := bufio.NewScanner(bytes.NewReader(stderr))
	found := false
	for s.Scan() {
		line := s.Bytes()
		if bytes.Contains(line, checkStderr) {
			found = true
			break
		}
	}
	if !found {
		return errorFromStderr(stderr)
	}

	fmt.Printf("%s", stderr)
	return nil
}

// SaveCmdOut saves both standard out and error to files, and print the standard out (if any).
// fnCheckStderr is a function to check the standard error.
func SaveCmdOut(dir, filename string, stdout, stderr []byte, fnCheckStderr func([]byte) error,
) (err error) {
	if stderr != nil {
		if fnCheckStderr != nil {
			if err = fnCheckStderr(stderr); err != nil {
				return err
			}
		}
		err = os.WriteFile(filepath.Join(dir, filename+"_stderr.log"), stderr, 0600)
		if err != nil {
			return err
		}
	}
	if stdout != nil {
		fmt.Println(string(stdout))
		err = os.WriteFile(filepath.Join(dir, filename+"_stdout.log"), stdout, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}
