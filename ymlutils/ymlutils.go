// This file is part of go-utils.
//
// Copyright (C) 2016  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package ymlutils - Utilities to read yml files like if using xpath
*/
package ymlutils

import (
	"fmt"
	"github.com/davidgamba/go-utils/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// YML object
type YML struct {
	Filename string
	data     []byte
	Tree     interface{}
}

// New returns a pointer to a YML object
func New(filename string) (*YML, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var tree interface{}
	err = yaml.Unmarshal(data, &tree)
	if err != nil {
		return nil, err
	}
	return &YML{Filename: filename, data: data, Tree: tree}, nil
}

// NavigateTree passthrough to utils.NavigateTree
func NavigateTree(i interface{}, path []string) (interface{}, bool) {
	return utils.NavigateTree(i, path)
}

// YMLQuery returns the yml object designated by path.
// Example:
// YMLQuery("level1/level2/level3")
func (y *YML) YMLQuery(path string) (interface{}, bool) {
	current, ok := utils.NavigateTree(y.Tree, strings.Split(path, "/"))
	return current, ok
}

// YMLGetString returns a string designated by path.
func (y *YML) YMLGetString(path string) (string, error) {
	i, ok := y.YMLQuery(path)
	if ok {
		switch i.(type) {
		case string:
			return i.(string), nil
		default:
			return "", fmt.Errorf("yaml path didn't return a valid string")
		}
	}
	return "", fmt.Errorf("yaml path not found")
}

// Unmarshal will unmarshal to the given structure.
func (y *YML) Unmarshal(out interface{}) error {
	return yaml.Unmarshal(y.data, out)
}
