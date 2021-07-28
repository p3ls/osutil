// Copyright 2019 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var listWinShell = []WinShell{
	Cmd,
	Db2,
}

// WinShell represents a Windows shell.
type WinShell uint8

// List of Windows shells.
const (
	_ WinShell = iota
	Cmd
	Db2
)

func (sh WinShell) String() string {
	switch sh {
	case Cmd:
		return "cmd.exe"
	case Db2:
		return "db2cmd.exe"
	}
	panic("unreachable")
}

var cmdOut = filepath.Join(os.TempDir(), "cmd_out.txt")

// ExecWinshell executes a command into a Windows' shell called from Powershell.
// Returns the command output.
func ExecWinshell(sh WinShell, closeWindow bool, cmd string) (out []byte, err error) {
	Log.Print(cmd)

	argClose := ""
	if closeWindow {
		argClose = "/C"
	} else if sh == Cmd {
		argClose = "/K"
	}

	defer func() {
		if err2 := os.Remove(cmdOut); err2 != nil {
			Log.Print(err2)
		}
	}()
	// The command output is saved into a file with the exit status code.
	arg := fmt.Sprintf(
		`start %s '%s echo Wait... & %s > %s & echo "exit status: %cerrorlevel%c" >> %s' -v runAs -wait`,
		sh, argClose, cmd, cmdOut, '%', '%', cmdOut,
	)

	if err = exec.Command("powershell.exe", arg).Run(); err != nil {
		return nil, err
	}

	return os.ReadFile(cmdOut)
}
