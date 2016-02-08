// This file is part of go-utils.
//
// Copyright (C) 2015-2016  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package fileutils - file related utilities
*/
package fileutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string) error {
	fmt.Printf("Copy: %s %s\n", src, dst)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}

// ChannelString is a struct containing the string `Data` and error `Error`.
type ChannelString struct {
	Data  string
	Error error
}

// GetFileList returns a channel with each file (`channel.string`) or an error indicating failure (`channel.error`).
func GetFileList(filename string, ignoreDirs bool) <-chan ChannelString {
	c := make(chan ChannelString)
	go func() {
		fInfo, err := os.Stat(filename)
		if err != nil {
			c <- ChannelString{"", err}
			return
		}
		if fInfo.IsDir() {
			if ignoreDirs == false {
				c <- ChannelString{filename, nil}
			}
			fileSearch := filename + string(filepath.Separator) + "*"
			fileMatches, err := filepath.Glob(fileSearch)
			if err != nil {
				c <- ChannelString{"", err}
				return
			}
			for _, file := range fileMatches {
				if filepath.Base(filename) == filepath.Base(file) {
					continue
				}
				d := GetFileList(file, ignoreDirs)
				for dirFile := range d {
					if dirFile.Error != nil {
						return
					}
					c <- ChannelString{dirFile.Data, nil}
				}
			}
		} else {
			c <- ChannelString{filename, nil}
		}
		close(c)
	}()
	return c
}
