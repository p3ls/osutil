// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import (
	"io"

	"github.com/tredoe/osutil/v2/executil"
)

// ManagerVoid is the interface to pass a Manager with an error,
// avoiding to have to use a pointer.
type ManagerVoid struct{}

func (m ManagerVoid) setPathExec(string) {}

func (m ManagerVoid) Cmd() *executil.Command { return nil }

func (m ManagerVoid) PackageType() string { return "" }

func (m ManagerVoid) PathExec() string { return "" }

func (m ManagerVoid) PreUsage() error { return nil }

func (m ManagerVoid) SetStdout(out io.Writer) {}

// * * *

func (m ManagerVoid) Install(name ...string) error { return nil }

func (m ManagerVoid) Remove(name ...string) error { return nil }

func (m ManagerVoid) Purge(name ...string) error { return nil }

func (m ManagerVoid) UpdateIndex() error { return nil }

func (m ManagerVoid) Update() error { return nil }

func (m ManagerVoid) Clean() error { return nil }

// * * *

func (m ManagerVoid) ImportKey(alias, keyUrl string) error { return nil }

func (m ManagerVoid) ImportKeyFromServer(alias, keyServer, key string) error { return nil }

func (m ManagerVoid) RemoveKey(alias string) error { return nil }

func (m ManagerVoid) AddRepo(alias string, url ...string) error { return nil }

func (m ManagerVoid) RemoveRepo(string) error { return nil }
