// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package service

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tredoe/osutil/v2/sysutil"
	"github.com/tredoe/osutil/v2/userutil"
)

// LookupService returns the service that matchs a pattern using the syntax of Match.
// If 'exclude' is set, it discards the names that contains it.
// When a service name is not found, returns error 'ServNotFoundError'.
func LookupService(
	sys sysutil.System,
	dis sysutil.Distro,
	pattern, exclude string,
	column ColumnWin,
) (*Service, error) {
	err := userutil.CheckSudo(sys)
	if err != nil {
		return nil, err
	}

	switch sys {
	case sysutil.Linux:
		return lookupServiceLinux(dis, pattern, exclude)
	case sysutil.MacOS:
		return lookupServiceMacos(pattern, exclude)
	case sysutil.Windows:
		return lookupServiceWindows(pattern, exclude, column)
	case sysutil.FreeBSD:
		return lookupServiceFreebsd(pattern, exclude)
	}

	panic("unreachable")
}

func lookupServiceFreebsd(pattern, exclude string) (*Service, error) {
	dirs := []string{"/usr/local/etc/rc.d/"}

	for _, dir := range dirs {
		files, err := filepath.Glob(dir + pattern)
		if err != nil {
			return nil, err
		}
		if files == nil {
			continue
		}

		for i := len(files) - 1; i >= 0; i-- {
			pathService := files[i]

			if exclude != "" && strings.Contains(pathService, exclude) {
				continue
			}
			serviceName := filepath.Base(pathService)

			return &Service{
				name: serviceName,
				path: pathService,
				sys:  sysutil.FreeBSD,
			}, nil
		}
	}

	return nil, ServNotFoundError{pattern: pattern, dirs: dirs}
}

func lookupServiceLinux(dis sysutil.Distro, pattern, exclude string) (*Service, error) {
	var dirs []string

	switch dis {
	case sysutil.Debian, sysutil.Ubuntu:
		dirs = []string{
			"/lib/systemd/system/",
		}
	case sysutil.CentOS, sysutil.Fedora:
		dirs = []string{
			"/lib/systemd/system/",
			"/etc/init.d/",
		}
	case sysutil.OpenSUSE:
		dirs = []string{
			"/usr/lib/systemd/system/",
			"/etc/init.d/",
		}

	default:
		dirs = []string{
			"/lib/systemd/system/",
			"/usr/lib/systemd/system/",
			"/etc/init.d/",
		}
	}

	for _, dir := range dirs {
		files, err := filepath.Glob(dir + pattern)
		if err != nil {
			return nil, err
		}
		if files == nil {
			continue
		}

		for i := len(files) - 1; i >= 0; i-- {
			pathService := files[i]

			if exclude != "" && strings.Contains(pathService, exclude) {
				continue
			}
			if strings.Contains(pathService, "@") {
				continue
			}

			//fmt.Println("SERVICE:", pathService) // DEBUG
			serviceName := filepath.Base(pathService)

			// The file could be finished with an extension like '.service' or '.target',
			// and with several dots like 'firebird3.0.service'
			if strings.Contains(serviceName, ".") {
				idLastDot := strings.LastIndex(serviceName, ".")
				part1 := serviceName[:idLastDot]
				part2 := serviceName[idLastDot+1:] // discard dot

				if len(part2) > 2 {
					serviceName = part1
				}
			}

			return &Service{
				name: serviceName,
				path: pathService,
				sys:  sysutil.Linux,
				dis:  dis,
			}, nil
		}
	}

	return nil, ServNotFoundError{pattern: pattern, dirs: dirs}
}

// Handle services installed through HomeBrew:
//
// + brew services list
// + ls -l ~/Library/LaunchAgents/

func lookupServiceMacos(pattern, exclude string) (*Service, error) {
	dirs := []string{ // Installed by
		fmt.Sprintf("/Library/LaunchDaemons/*.%s*", pattern),   // binary installer
		fmt.Sprintf("/usr/local/Cellar/%s/*/*.plist", pattern), // HomeBrew
	}

	for iDir, dir := range dirs {
		files, err := filepath.Glob(dir)
		if err != nil {
			return nil, err
		}
		if files == nil {
			continue
		}

		for i := 0; i < len(files); i++ {
			pathService := files[i]

			serviceName := ""
			switch iDir {
			case 0:
				serviceName = strings.SplitAfter(pathService, "/Library/LaunchDaemons/")[1]
				//if split := strings.SplitN(serviceName, ".plist", 2); len(split) != 1 {
				//	serviceName = split[0]
				//}
				serviceName = strings.SplitN(serviceName, ".plist", 2)[0]
			case 1:
				serviceName = strings.SplitAfter(pathService, "/usr/local/Cellar/")[1]
				serviceName = strings.SplitN(serviceName, "/", 2)[0]
			}

			if exclude != "" && strings.Contains(serviceName, exclude) {
				continue
			}

			return &Service{
				name: serviceName,
				path: pathService,
				sys:  sysutil.MacOS,
			}, nil
		}
	}

	return nil, ServNotFoundError{pattern: pattern, dirs: dirs}
}

func lookupServiceWindows(pattern, exclude string, column ColumnWin) (*Service, error) {
	var out bytes.Buffer
	cmd := exec.Command(
		"powershell.exe",
		fmt.Sprintf("Get-Service -%s %q | Select-Object Name", column, pattern),
	)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	rd := bytes.NewReader(out.Bytes())
	sc := bufio.NewScanner(rd)
	line := ""

	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())

		if line == "" {
			continue
		}
		if line[0] == '-' || strings.HasPrefix(line, "Name") {
			line = ""
			continue
		}

		break
	}
	if err = sc.Err(); err != nil {
		return nil, err
	}
	if line == "" {
		return nil, ServNotFoundError{pattern: pattern}
	}

	return &Service{
		name: line,
		sys:  sysutil.Windows,
	}, nil
}
