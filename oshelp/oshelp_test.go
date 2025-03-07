// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package oshelp

import (
	"reflect"
	"testing"

	"github.com/drone/runner-go/shell/bash"
	"github.com/drone/runner-go/shell/powershell"
)

func Test_join(t *testing.T) {
	tests := []struct {
		os string
		a  []string
		b  string
	}{
		{os: "windows", a: []string{"C:", "Windows", "Temp"}, b: "C:\\Windows\\Temp"},
		{os: "linux", a: []string{"/tmp", "foo", "bar"}, b: "/tmp/foo/bar"},
	}
	for _, test := range tests {
		if got, want := JoinPaths(test.os, test.a...), test.b; got != want {
			t.Errorf("Want %s, got %s", want, got)
		}
	}
}

func Test_getExt(t *testing.T) {
	tests := []struct {
		os string
		a  string
		b  string
	}{
		{os: "windows", a: "clone", b: "clone.ps1"},
		{os: "linux", a: "clone", b: "clone"},
	}
	for _, test := range tests {
		if got, want := GetExt(test.os, test.a), test.b; got != want {
			t.Errorf("Want %s, got %s", want, got)
		}
	}
}

func Test_getNetrc(t *testing.T) {
	tests := []struct {
		os   string
		name string
	}{
		{os: "windows", name: "_netrc"},
		{os: "linux", name: ".netrc"},
		{os: "openbsd", name: ".netrc"},
		{os: "netbsd", name: ".netrc"},
		{os: "freebsd", name: ".netrc"},
	}
	for _, test := range tests {
		if got, want := GetNetrc(test.os), test.name; got != want {
			t.Errorf("Want %s on %s, got %s", want, test.os, got)
		}
	}
}

func Test_getScript(t *testing.T) {
	commands := []string{"go build"}

	a := GenScript("windows", commands)
	b := powershell.Script(commands)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Generated windows linux script")
	}

	a = GenScript("linux", commands)
	b = bash.Script(commands)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Generated invalid linux script")
	}
}
