// Copyright (c) 2016 Betalo AB
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"os"
	"testing"
)

func TestSplitArgsResourceAndCommand(t *testing.T) {
	os.Args = []string{"./await", "true", "--", "echo"}
	flag.Parse()
	res, cmd := splitArgs(flag.Args())
	if len(res) != 1 || res[0] != "true" {
		t.Errorf("no resources found")
	}
	if len(cmd) != 1 || cmd[0] != "echo" {
		t.Errorf("no command found")
	}
}

func TestSplitArgsResourceAndNoCommand(t *testing.T) {
	os.Args = []string{"./await", "true", "--"}
	flag.Parse()
	res, cmd := splitArgs(flag.Args())
	if len(res) != 1 || res[0] != "true" {
		t.Errorf("no resources found")
	}
	if len(cmd) != 0 {
		t.Errorf("unexpected command found")
	}

	os.Args = []string{"./await", "true"}
	flag.Parse()
	res, cmd = splitArgs(flag.Args())
	if len(res) != 1 || res[0] != "true" {
		t.Errorf("no resources found")
	}
	if len(cmd) != 0 {
		t.Errorf("unexpected command found")
	}
}

func TestSplitArgsNoResourceAndCommand(t *testing.T) {
	os.Args = []string{"./await", "--", "echo"}
	flag.Parse()
	res, cmd := splitArgs(flag.Args())
	if len(res) != 0 {
		t.Errorf("unexpected resource found")
	}
	if len(cmd) != 1 || cmd[0] != "echo" {
		t.Errorf("no command found")
	}
}

func TestSplitArgsNoResourceNoCommand(t *testing.T) {
	os.Args = []string{"./await", "--"}
	flag.Parse()
	res, cmd := splitArgs(flag.Args())
	if len(res) != 0 {
		t.Errorf("unexpected resource found")
	}
	if len(cmd) != 0 {
		t.Errorf("unexpected command found")
	}

	os.Args = []string{"./await"}
	flag.Parse()
	res, cmd = splitArgs(flag.Args())
	if len(res) != 0 {
		t.Errorf("unexpected resource found")
	}
	if len(cmd) != 0 {
		t.Errorf("unexpected command found")
	}
}

func TestParseResourcesSuccess(t *testing.T) {
	ress := []string{
		"http://user:pass@localhost:foo/path?query=val#fragment",
		"https://user:pass@localhost:foo/path?query=val#fragment",

		"ws://user:pass@localhost:42/path?query=val#fragment",
		"wss://user:pass@localhost:42/path?query=val#fragment",

		"tcp://localhost:42",
		"tcp4://localhost:42",
		"tcp6://[::1]:42",

		"file://relative/path/to/file",
		"file://relative/path/to/file#absent",
		"file:///absolute/path/to/file",
		"file://absolute/path/to/file#absent",

		"postgres://user:pass@localhost:5432/dbname?query=val#fragment",

		"mysql://user:pass@localhost:3306/dbname?query=val#fragment",

		"command",
		"command with args",
		"relative/path/to/command",
		"relative/path/to/command with args",
		"/absolute/path/to/command",
		"/absolute/path/to/command with args",

		"",
	}

	actualRess, err := parseResources(ress)
	if err != nil {
		t.Errorf("failed to parse ressources")
	}
	if len(actualRess) != len(ress) {
		t.Errorf("missing parsed ressources")
	}
}

func TestParseResourcesFailure(t *testing.T) {
	_, err := parseResources([]string{"//foo test"})
	if err == nil {
		t.Errorf("expected error parsing invalid ressource")
	}
}

func TestParseTags(t *testing.T) {
	tests := map[string]map[string]string{
		"":              map[string]string{},
		"=":             map[string]string{}, // invalid format, should be skipped
		"=ignore":       map[string]string{}, // invalid format, should be skipped
		"foo":           map[string]string{"foo": ""},
		"foo=":          map[string]string{"foo": ""},
		"foo=bar":       map[string]string{"foo": "bar"},
		"foo=bar&baz":   map[string]string{"foo": "bar", "baz": ""},
		"foo=bar&baz=1": map[string]string{"foo": "bar", "baz": "1"},
	}
	for given, expected := range tests {
		actual := parseTags(given)
		if len(actual) != len(expected) {
			t.Errorf("unexpected parsed tags count %d != %d: %#v",
				len(expected), len(actual), actual)
		}
		for actualK, actualV := range actual {
			if expectedV, ok := expected[actualK]; !ok || actualV != expectedV {
				t.Errorf("unexpected parsed tags k/v pair")
			}
		}
	}
}
