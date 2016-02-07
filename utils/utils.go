// This file is part of go-utils.
//
// Copyright (C) 2016  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package utils provides generic utils I can't think of another place to put into.
*/
package utils

import (
// "fmt"
// "reflect"
)

// NavigateTree allows you to define a path string to traverse a tree composed of maps and arrays.
func NavigateTree(m interface{}, p []string) (interface{}, bool) {
	// fmt.Printf("type: %v, path: %v\n", reflect.TypeOf(m), p)
	if len(p) <= 0 {
		return m, true
	}
	switch m.(type) {
	case map[interface{}]interface{}:
		switch t := m.(map[interface{}]interface{})[p[0]].(type) {
		case string:
			if len(p) == 1 {
				return t, true
			}
			return t, false
		case map[interface{}]interface{}:
			return NavigateTree(t, p[1:])
		case []interface{}:
			return NavigateTree(t, p[1:])
		default:
			return t, false
		}
	case []interface{}:
		var i interface{}
		var b bool
		for _, v := range m.([]interface{}) {
			i, b = NavigateTree(v, p)
			if b {
				return i, b
			}
		}
		return m, false
	default:
		return m, false
	}
}
