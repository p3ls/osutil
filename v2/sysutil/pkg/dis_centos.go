// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import (
	"io"
	"os"

	"github.com/tredoe/osutil/v2"
	"github.com/tredoe/osutil/v2/executil"
	"github.com/tredoe/osutil/v2/internal"
	"github.com/tredoe/osutil/v2/sysutil"
)

const (
	fileDnf = "dnf" // Preferable to YUM
	pathDnf = "/usr/bin/dnf"

	fileYum    = "yum"
	pathYum    = "/usr/bin/yum"
	pathYumCfg = "/usr/bin/yum-config-manager"

	// RPM is used to install/uninstall local packages.
	fileRpm = "rpm"
	pathRpm = "/usr/bin/rpm"
)

// == Dnf
//

// ManagerDnf is the interface to handle the package manager DNG of Linux systems
// based at Red Hat.
type ManagerDnf struct {
	pathExec string
	cmd      *executil.Command

	rpm ManagerRpm
}

// NewManagerDnf returns the DNF package manager.
// It checks if it is being run by an administrator.
func NewManagerDnf() (ManagerDnf, error) {
	if err := internal.CheckSuperuser(sysutil.Linux); err != nil {
		return ManagerDnf{}, err
	}

	managerRpm, _ := NewManagerRpm()

	return ManagerDnf{
		pathExec: pathDnf,
		cmd: cmd.Command("", "").
			// https://dnf.readthedocs.io/en/latest/command_ref.html
			BadExitCodes([]int{1, 3, 200}),
		rpm: managerRpm,
	}, nil
}

func (m ManagerDnf) setPathExec(p string) { m.pathExec = p }

func (m ManagerDnf) Cmd() *executil.Command { return m.cmd }

func (m ManagerDnf) PackageType() string { return Dnf.String() }

func (m ManagerDnf) PathExec() string { return m.pathExec }

func (m ManagerDnf) PreUsage() error {
	// Required to use: "dnf config-manager"
	return m.Install("dnf-plugins-core")
}

func (m ManagerDnf) SetStdout(out io.Writer) { m.cmd.Stdout(out) }

// * * *

func (m ManagerDnf) Install(name ...string) error {
	args := append([]string{"install", "-y"}, name...)

	_, err := m.cmd.Command(pathDnf, args...).Run()
	return err
}

func (m ManagerDnf) Remove(name ...string) error {
	args := append([]string{"remove", "-y"}, name...)

	_, err := m.cmd.Command(pathDnf, args...).Run()
	return err
}

func (m ManagerDnf) Purge(name ...string) error {
	return m.Remove(name...)
}

func (m ManagerDnf) UpdateIndex() error {
	// check-update does not update else it checks the updating.
	return nil

	// check-update returns exit value of 100 if there are packages available for an update.
	// Also returns a list of the packages to be updated in list format.
	// Returns 0 if no packages are available for update.
	// Returns 1 if an error occurred.
	/*err := m.cmd.Command(pathDnf, "check-update")
	if err != nil {
		// Check the exit code
	}
	return err*/
}

func (m ManagerDnf) Update() error {
	_, err := m.cmd.Command(pathDnf, "update", "-y").Run()
	return err
}

func (m ManagerDnf) Clean() error {
	_, err := m.cmd.Command(pathDnf, "autoremove", "-y").Run()
	if err != nil {
		return err
	}
	_, err = m.cmd.Command(pathDnf, "clean", "all").Run()
	return err
}

// * * *

func (m ManagerDnf) ImportKey(alias, keyUrl string) error {
	return m.rpm.ImportKey("", keyUrl)
}

func (m ManagerDnf) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrManagCmd
}

func (m ManagerDnf) RemoveKey(alias string) error {
	return ErrManagCmd
}

// https://docs.fedoraproject.org/en-US/quick-docs/adding-or-removing-software-repositories-in-fedora/

func (m ManagerDnf) AddRepo(alias string, url ...string) error {
	/*pathRepo := m.repository(alias)

	err := fileutil.CreateFromString(pathRepo, url[0]+"\n")
	if err != nil {
		return err
	}*/

	stderr, err := m.cmd.Command(
		pathDnf, "config-manager", "--add-repo", url[0],
	).OutputStderr()

	return executil.CheckStderr(stderr, err)
}

func (m ManagerDnf) RemoveRepo(alias string) error {
	return os.Remove(m.repository(alias))
}

// == Yum
//

// ManagerYum is the interface to handle the package manager YUM of Linux systems
// based at Red Hat.
type ManagerYum struct {
	pathExec string
	cmd      *executil.Command

	rpm ManagerRpm
}

// NewManagerYum returns the YUM package manager.
// It checks if it is being run by an administrator.
func NewManagerYum() (ManagerYum, error) {
	if err := internal.CheckSuperuser(sysutil.Linux); err != nil {
		return ManagerYum{}, err
	}
	managerRpm, _ := NewManagerRpm()

	return ManagerYum{
		pathExec: pathYum,
		cmd: cmd.Command("", "").
			BadExitCodes([]int{1, 2, 3, 16}),
		rpm: managerRpm,
	}, nil
}

func (m ManagerYum) setPathExec(p string) { m.pathExec = p }

func (m ManagerYum) Cmd() *executil.Command { return m.cmd }

func (m ManagerYum) PackageType() string { return Yum.String() }

func (m ManagerYum) PathExec() string { return m.pathExec }

func (m ManagerYum) PreUsage() error {
	// Required to use: "yum-config-manager"
	return m.Install("yum-utils")
}

func (m ManagerYum) SetStdout(out io.Writer) { m.cmd.Stdout(out) }

// * * *

func (m ManagerYum) Install(name ...string) error {
	osutil.Log.Print(taskInstall)
	args := append([]string{"install", "-y"}, name...)

	_, err := m.cmd.Command(pathYum, args...).Run()
	return err
}

func (m ManagerYum) Remove(name ...string) error {
	osutil.Log.Print(taskRemove)
	args := append([]string{"remove", "-y"}, name...)

	_, err := m.cmd.Command(pathYum, args...).Run()
	return err
}

func (m ManagerYum) Purge(name ...string) error {
	osutil.Log.Print(taskPurge)
	return m.Remove(name...)
}

func (m ManagerYum) UpdateIndex() error {
	// check-update does not update else it checks the updating.
	return ErrManagCmd
}

func (m ManagerYum) Update() error {
	osutil.Log.Print(taskUpgrade)
	_, err := m.cmd.Command(pathYum, "update", "-y").Run()
	return err
}

func (m ManagerYum) Clean() error {
	osutil.Log.Print(taskClean)
	_, err := m.cmd.Command(pathYum, "clean", "packages").Run()
	return err
}

// * * *

func (m ManagerYum) ImportKey(alias, keyUrl string) error {
	osutil.Log.Print(taskImportKey)
	return m.rpm.ImportKey("", keyUrl)
}

func (m ManagerYum) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrManagCmd
}

func (m ManagerYum) RemoveKey(alias string) error {
	return ErrManagCmd
}

// https://docs.fedoraproject.org/en-US/Fedora/16/html/System_Administrators_Guide/sec-Managing_Yum_Repositories.html

func (m ManagerYum) AddRepo(alias string, url ...string) error {
	osutil.Log.Print(taskAddRepo)
	stderr, err := m.cmd.Command(
		pathYumCfg, "--add-repo", url[0],
	).OutputStderr()

	return executil.CheckStderr(stderr, err)
}

func (m ManagerYum) RemoveRepo(alias string) error {
	osutil.Log.Print(taskRemoveRepo)
	return os.Remove(m.repository(alias))
}

// == RPM
//

// ManagerRpm is the interface to handle the package manager RPM of Linux systems
// based at Red Hat.
type ManagerRpm struct {
	pathExec string
	cmd      *executil.Command
}

// NewManagerRpm returns the RPM package manager.
// It checks if it is being run by an administrator.
func NewManagerRpm() (ManagerRpm, error) {
	if err := internal.CheckSuperuser(sysutil.Linux); err != nil {
		return ManagerRpm{}, err
	}

	return ManagerRpm{
		pathExec: pathRpm,
		cmd:      cmd.Command("", ""),
		//BadExitCodes([]int{}),
	}, nil
}

func (m ManagerRpm) setPathExec(p string) { m.pathExec = p }

func (m ManagerRpm) Cmd() *executil.Command { return m.cmd }

func (m ManagerRpm) PackageType() string { return Rpm.String() }

func (m ManagerRpm) PathExec() string { return m.pathExec }

func (m ManagerRpm) PreUsage() error { return nil }

func (m ManagerRpm) SetStdout(out io.Writer) { m.cmd.Stdout(out) }

// * * *

func (m ManagerRpm) Install(name ...string) error {
	osutil.Log.Print(taskInstall)
	args := append([]string{"-i"}, name...)

	_, err := m.cmd.Command(pathRpm, args...).Run()
	return err
}

func (m ManagerRpm) Remove(name ...string) error {
	osutil.Log.Print(taskRemove)
	args := append([]string{"-e"}, name...)

	_, err := m.cmd.Command(pathRpm, args...).Run()
	return err
}

func (m ManagerRpm) Purge(name ...string) error {
	osutil.Log.Print(taskPurge)
	return m.Remove(name...)
}

func (m ManagerRpm) UpdateIndex() error {
	return ErrManagCmd
}

func (m ManagerRpm) Update() error {
	return ErrManagCmd
}

func (m ManagerRpm) Clean() error {
	return ErrManagCmd
}

// * * *

func (m ManagerRpm) ImportKey(alias, keyUrl string) error {
	osutil.Log.Print(taskImportKey)
	stderr, err := m.cmd.Command(pathRpm, "--import", keyUrl).OutputStderr()

	err = executil.CheckStderr(stderr, err)
	return err
}

func (m ManagerRpm) ImportKeyFromServer(alias, keyServer, key string) error {
	return ErrManagCmd
}

func (m ManagerRpm) RemoveKey(alias string) error {
	return ErrManagCmd
}

func (m ManagerRpm) AddRepo(alias string, url ...string) error {
	return ErrManagCmd
}

func (m ManagerRpm) RemoveRepo(r string) error {
	return ErrManagCmd
}

// == Utility
//

func (m ManagerDnf) repository(alias string) string {
	return "/etc/yum.repos.d/" + alias + ".repo"
}

func (m ManagerYum) repository(alias string) string {
	return "/etc/yum.repos.d/" + alias + ".repo"
}
