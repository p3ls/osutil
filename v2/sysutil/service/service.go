// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package service handles the services at systems Linux, FreeBSD, macOS and Windows.
//
// The information messages are written through 'logShell' configured at 'SetupLogger()'.
package service

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/tredoe/osutil/v2"
	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/internal"
	"github.com/tredoe/osutil/v2/sysutil"
)

// timeKillServ is the time used to wait before of kill a service.
const timeKillServ = 90 * time.Second

var (
	excmd    = executil.NewCommand("", "").TimeKill(timeKillServ).Env([]string{"LANG=C"})
	excmdWin = executil.NewCommand("", "").TimeKill(timeKillServ)
)

// ColumnWin represents the column name at Windows where to find the service name.
type ColumnWin uint8

const (
	ColWinName        = iota // Column 'Name', by default.
	ColWinDisplayname        // Column 'Displayname'
)

func (c ColumnWin) String() string {
	switch c {
	case ColWinName:
		return "Name"
	case ColWinDisplayname:
		return "DisplayName"

	default:
		panic("unimplemented")
	}
}

// Service represents a service.
type Service struct {
	name string
	path string
	sys  sysutil.System
	dis  sysutil.Distro

	// Custom commands to both start and stop.
	start *executil.Command
	stop  *executil.Command
}

// Name returns the service name.
func (s *Service) Name() string { return s.name }

// NewService creates a new service with the given name.
// It checks if it is being run by an administrator.
func NewService(
	sys sysutil.System, dis sysutil.Distro, name string,
) (*Service, error) {
	if name == "" {
		return nil, ErrNoService
	}
	if err := internal.CheckSuperuser(sys); err != nil {
		return nil, err
	}

	return &Service{
		name: name,
		sys:  sys,
		dis:  dis,
	}, nil
}

// NewCustomService creates a new service with the custom commands.
func NewCustomService(
	sys sysutil.System,
	cmdStart string, argsStart []string,
	cmdStop string, argsStop []string,
) *Service {
	s := new(Service)

	if cmdStart != "" {
		s.start = executil.NewCommand(cmdStart, argsStart...).TimeKill(timeKillServ)

		if s.sys != sysutil.Windows {
			s.start.Env([]string{"LANG=C"})
		}
	}
	if cmdStop != "" {
		s.stop = executil.NewCommand(cmdStop, argsStop...).TimeKill(timeKillServ)

		if s.sys != sysutil.Windows {
			s.stop.Env([]string{"LANG=C"})
		}
	}

	return s
}

// * * *

// Start starts the service.
func (srv Service) Start() error {
	var err error
	var stderr []byte

	osutil.Log.Print("Starting service ...")

	if srv.start != nil {
		stderr, err = srv.start.OutputStderr()
		return executil.CheckStderr(stderr, err)
	}

	switch srv.sys {
	case sysutil.Linux:
		stderr, err = excmd.Command(
			"systemctl", "start", srv.name,
		).OutputStderr()

	case sysutil.FreeBSD:
		stderr, err = excmd.Command(
			"service", srv.name, "start",
		).OutputStderr()

	case sysutil.MacOS:
		stderr, err = excmd.Command(
			"launchctl", "load", "-F", srv.name,
		).OutputStderr()

		if err != nil {
			if !bytes.Contains(stderr, []byte("service already loaded")) {
				return fmt.Errorf("%s", stderr)
			}
			//logs.Debug.Printf("%s\n%s", stderr, err)
		}

	case sysutil.Windows:
		stderr, err = excmdWin.Command(
			"net", "start", srv.name,
		).OutputStderr()

		if err != nil {
			if !bytes.Contains(stderr, []byte("already been started")) {
				return fmt.Errorf("%s", stderr)
			}
			//logs.Debug.Printf("%s\n%s", stderr, err)
		}

	default:
		panic("unimplemented: " + srv.sys.String())
	}

	return executil.CheckStderr(stderr, err)
}

// Stop stops the service.
func (srv Service) Stop() error {
	osutil.Log.Print("Stopping service ...")

	if srv.stop != nil {
		stderr, err := srv.stop.OutputStderr()
		return executil.CheckStderr(stderr, err)
	}

	switch srv.sys {
	case sysutil.Linux:
		stdout, stderr, err := excmd.Command(
			"systemctl", "is-active", srv.name,
		).OutputCombined()

		if err = executil.CheckStderr(stderr, err); err != nil {
			return err
		}

		if bytes.HasPrefix(stdout, []byte("active")) {
			stderr, err := excmd.Command(
				"systemctl", "stop", srv.name,
			).OutputStderr()

			if err = executil.CheckStderr(stderr, err); err != nil {
				return err
			}
		}

	case sysutil.FreeBSD:
		stderr, err := excmd.Command(
			"service", srv.name, "stop",
		).OutputStderr()

		if err = executil.CheckStderr(stderr, err); err != nil {
			return err
		}

	case sysutil.MacOS:
		stderr, err := excmd.Command(
			"launchctl", "unload", "-F", srv.name,
		).OutputStderr()

		if stderr != nil {
			if bytes.Contains(stderr, []byte("Operation now in progress")) {
				return nil
			}
			if !bytes.Contains(stderr, []byte("Could not find specified service")) {
				return fmt.Errorf("%s", stderr)
			}
			//logs.Debug.Printf("%s\n%s", stderr, err)
		}
		if err != nil {
			return err
		}

	case sysutil.Windows:
		stderr, err := excmdWin.Command(
			"net", "stop", srv.name,
		).OutputStderr()

		if stderr != nil {
			if !bytes.Contains(stderr, []byte("is not started")) {
				return fmt.Errorf("%s", stderr)
			}
			//logs.Debug.Printf("%s\n%s", stderr, err)
		}
		if err != nil {
			return err
		}

	default:
		panic("unimplemented: " + srv.sys.String())
	}

	return nil
}

// Restart stops and starts the service.
func (srv Service) Restart() error {
	switch srv.sys {
	case sysutil.Linux:
		osutil.Log.Print("Re-starting service ...")

		stderr, err := excmd.Command(
			"systemctl", "restart", srv.name,
		).OutputStderr()

		if err = executil.CheckStderr(stderr, err); err != nil {
			return err
		}

		// Wait to restart the service.
		time.Sleep(1 * time.Second)

		return nil

	case sysutil.FreeBSD:
		stderr, err := excmd.Command(
			"service", srv.name, "restart",
		).OutputStderr()

		if err = executil.CheckStderr(stderr, err); err != nil {
			return err
		}

		return nil

	default:
		if err := srv.Stop(); err != nil {
			return err
		}
		return srv.Start()
	}
}

// Enable enables the service.
func (srv Service) Enable() error {
	var err error
	var stderr []byte

	osutil.Log.Print("Enabling service ...")

	switch srv.sys {
	case sysutil.Linux:
		args := []string{"systemctl", "enable", srv.name}

		switch srv.dis {
		case sysutil.CentOS:
			ver := 7

			if ver < 7 {
				args = []string{"chkconfig", srv.name, "on"}
			}
		}

		stderr, err = excmd.Command(args[0], args[1:]...).OutputStderr()

	//case sysutil.FreeBSD:
	//"sysrc sshd_enable='YES'"

	case sysutil.MacOS:
		stderr, err = excmd.Command(
			"launchctl", "enable", srv.name,
		).OutputStderr()

	case sysutil.Windows:
		stderr, err = excmdWin.Command(
			"sc", "config", srv.name, "start= demand",
		).OutputStderr()

	default:
		panic("unimplemented: " + srv.sys.String())
	}

	return executil.CheckStderr(stderr, err)
}

// Disable disables the service.
func (srv Service) Disable() error {
	var err error
	var stderr []byte

	osutil.Log.Print("Disabling service ...")

	switch srv.sys {
	case sysutil.Linux:
		args := []string{"systemctl", "disable", srv.name}

		switch srv.dis {
		case sysutil.CentOS:
			ver := 7

			if ver < 7 {
				args = []string{"chkconfig", srv.name, "off"}
			}
		}

		stderr, err = excmd.Command(args[0], args[1:]...).OutputStderr()

	//case sysutil.FreeBSD:
	//"sysrc sshd_enable='YES'"

	case sysutil.MacOS:
		stderr, err = excmd.Command(
			"launchctl", "disable", srv.name,
		).OutputStderr()

	case sysutil.Windows:
		stderr, err = excmdWin.Command(
			"sc", "config", srv.name, "start= disabled",
		).OutputStderr()

	default:
		panic("unimplemented: " + srv.sys.String())
	}

	return executil.CheckStderr(stderr, err)
}

// == Errors
//
// ErrNoService represents an error
var ErrNoService = errors.New("no service name")

// ServNotFoundError indicates whether a service is not found.
type ServNotFoundError struct {
	pattern string
	dirs    []string
}

func (e ServNotFoundError) Error() string {
	return fmt.Sprintf(
		"failed at searching the service by pattern %q at directories %v",
		e.pattern, e.dirs,
	)
}
