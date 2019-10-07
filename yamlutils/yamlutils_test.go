// This file is part of go-utils.
//
// Copyright (C) 2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package yamlutils

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		input    string
		expected string
		err      error
	}{
		{"simple", []string{}, "hello", "hello", nil},
		{"extra elements in path", []string{"hello"}, "hello", "hello", ErrExtraElementsInPath},
		{"simple", []string{"hello"}, "hello: world", "world", nil},
		{"simple", []string{"x"}, "hello: world", "hello: world\n", ErrMapKeyNotFound},
		{"simple", []string{"hello", "1"}, `hello:
  - one
  - two
  - three`, "two", nil},
		{"simple", []string{"hello", "3"}, `hello:
  - one
  - two
  - three`, "- one\n- two\n- three\n", ErrInvalidIndex},
		{"simple", []string{"hello", "-1"}, `hello:
  - one
  - two
  - three`, "- one\n- two\n- three\n", ErrInvalidIndex},
		{"simple", []string{"hello", "1", "world"}, `hello:
  - one
  - world: hola
  - three`, "hola", nil},
		{"simple", []string{"hello", "1", "world"}, `hello:
  - one
  - world: true
  - three`, "true", nil},
		{"simple", []string{"hello", "1", "world"}, `hello:
  - one
  - world: 123
  - three`, "123", nil},
		{"simple", []string{"hello", "1", "world"}, `hello:
  - one
  - world: 123.123
  - three`, "123.123", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := ""
			buf := bytes.NewBufferString(s)
			Logger.SetOutput(buf)
			yml, err := NewFromString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %s\n", err)
			}
			output, err := yml.GetString(test.path)
			if !errors.Is(err, test.err) {
				t.Errorf("Unexpected error: %s\n", err)
			}
			if !reflect.DeepEqual(output, test.expected) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v\n", test.expected, output)
			}
			t.Log(buf.String())
		})
	}

}

func TestNavigateTree(t *testing.T) {
	tests := []struct {
		name         string
		path         []string
		input        interface{}
		expected     interface{}
		expectedPath []string
		err          error
	}{
		{"string config, no path",
			[]string{}, "hola", "hola", []string{}, nil},
		{"string config, bad path",
			[]string{"extra", "elements"}, "hola", "hola", []string{"extra", "elements"}, ErrExtraElementsInPath},
		{"array config, no path",
			[]string{}, []interface{}{"one", "two", "three"}, []interface{}{"one", "two", "three"}, []string{}, nil},
		{"array config, path",
			[]string{"1"}, []interface{}{"one", "two", "three"}, "two", []string{}, nil},
		{"array config, bad path",
			[]string{"[1]"}, []interface{}{"one", "two", "three"}, []interface{}{"one", "two", "three"}, []string{"[1]"}, ErrNotAnIndex},
		{"map config, no path",
			[]string{},
			map[interface{}]interface{}{"map": []string{"one", "two", "three"}},
			map[interface{}]interface{}{"map": []string{"one", "two", "three"}},
			[]string{}, nil},
		{"map config, path",
			[]string{"map"},
			map[interface{}]interface{}{"map": []string{"one", "two", "three"}},
			[]string{"one", "two", "three"},
			[]string{}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := ""
			buf := bytes.NewBufferString(s)
			Logger.SetOutput(buf)
			output, path, err := NavigateTree(test.input, test.path)
			if !errors.Is(err, test.err) {
				t.Errorf("Unexpected error: %s\n", err)
			}
			if !reflect.DeepEqual(path, test.expectedPath) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v\n", test.expectedPath, path)
			}
			if !reflect.DeepEqual(output, test.expected) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v\n", test.expected, output)
			}
			t.Log(buf.String())
		})
	}
}
