// This file is part of go-utils.
//
// Copyright (C) 2016-2019  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package yamlutils - Utilities to read yml files like if using xpath
*/
package yamlutils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	// "reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Logger - Custom lib logger
var Logger = log.New(ioutil.Discard, "yamlutils ", log.LstdFlags)

// YML object
type YML struct {
	Tree interface{}
}

// NewFromFile returns a pointer to a YML object from a file.
func NewFromFile(filename string) (*YML, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var tree interface{}
	err = yaml.Unmarshal(data, &tree)
	if err != nil {
		return nil, err
	}
	return &YML{Tree: tree}, nil
}

// NewFromReader returns a pointer to a YML object from an io.Reader.
func NewFromReader(reader io.Reader) (*YML, error) {
	var tree interface{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf.Bytes(), &tree)
	if err != nil {
		return nil, err
	}
	return &YML{Tree: tree}, nil
}

// NewFromString - returns a pointer to a YML object from a string.
func NewFromString(str string) (*YML, error) {
	var tree interface{}
	err := yaml.Unmarshal([]byte(str), &tree)
	if err != nil {
		return nil, err
	}
	return &YML{Tree: tree}, nil
}

// GetString returns a string designated by path.
// Path is a string with elements separated by /.
// Array indexes are given as a number.
// For example: "level1/level2/3/level4"
func (y *YML) GetString(keys []string) (string, error) {
	path := strings.Join(keys, ",")
	target, _, errPath := NavigateTree(y.Tree, keys)
	out, err := yaml.Marshal(target)
	if errPath != nil {
		return string(out), fmt.Errorf("yaml path '%s' didn't return a valid string: %w", path, errPath)
	}
	if err != nil {
		return string(out), fmt.Errorf("failed to Marshal output: %w", err)
	}
	return string(out), nil
}

// ErrExtraElementsInPath - Indicates when there is a final match and there are remaining path elements.
var ErrExtraElementsInPath = fmt.Errorf("extra elements in path")

// ErrMapKeyNotFound - Key not in config.
var ErrMapKeyNotFound = fmt.Errorf("map key not found")

// ErrNotAnIndex - The given path is not a numerical index and the element is of type slice/array.
var ErrNotAnIndex = fmt.Errorf("not an index")

// ErrInvalidIndex - The given index is invalid.
var ErrInvalidIndex = fmt.Errorf("invalid index")

// NavigateTree allows you to define a path string to traverse a tree composed of maps and arrays.
// To navigate through slices/arrays use a numerical index, for example: [path to array 1]
func NavigateTree(m interface{}, p []string) (interface{}, []string, error) {
	// Logger.Printf("type: %v, path: %v\n", reflect.TypeOf(m), p)
	path := strings.Join(p, "/")
	Logger.Printf("NavigateTree: Input path: %s", path)
	if len(p) <= 0 {
		return m, p, nil
	}
	switch m.(type) {
	case map[interface{}]interface{}:
		Logger.Printf("NavigateTree: map type")
		t, ok := m.(map[interface{}]interface{})[p[0]]
		if !ok {
			return m, p, fmt.Errorf("%w: %s", ErrMapKeyNotFound, p[0])
		}
		return NavigateTree(t, p[1:])
	case []interface{}:
		Logger.Printf("NavigateTree: slice/array type")

		index, err := strconv.Atoi(p[0])
		if err != nil {
			return m, p, fmt.Errorf("%w: %s", ErrNotAnIndex, p[0])
		}
		if index < 0 || len(m.([]interface{})) <= index {
			return m, p, fmt.Errorf("%w: %s", ErrInvalidIndex, p[0])
		}
		return NavigateTree(m.([]interface{})[index], p[1:])
	default:
		Logger.Printf("NavigateTree: single element type")
		return m, p, fmt.Errorf("%w: %s", ErrExtraElementsInPath, strings.Join(p, "/"))
	}
}
