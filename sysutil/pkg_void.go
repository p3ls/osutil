// Copyright 2021 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sysutil

// ManagerVoid is the interface to pass a Manager with an error,
// avoiding to have to use a pointer.
type ManagerVoid struct{}

func (m ManagerVoid) setExecPath(string) {}

func (m ManagerVoid) ExecPath() string { return "" }

func (m ManagerVoid) PackageType() string { return "" }

func (m ManagerVoid) Install(name ...string) error { return nil }

func (m ManagerVoid) Remove(name ...string) error { return nil }

func (m ManagerVoid) Purge(name ...string) error { return nil }

func (m ManagerVoid) Update() error { return nil }

func (m ManagerVoid) Upgrade() error { return nil }

func (m ManagerVoid) Clean() error { return nil }
