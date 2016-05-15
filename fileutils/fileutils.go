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

// StringError is a struct containing the string `String` and error `Error`.
type StringError struct {
	String string
	Error  error
}

// GetFileList returns a channel with each file (`channel.String`) or an error indicating failure (`channel.Error`).
func GetFileList(dirname string, ignoreDirs, recursive bool) <-chan StringError {
	c := make(chan StringError)
	go func() {
		fInfo, err := os.Stat(dirname)
		if err != nil {
			c <- StringError{"", err}
			close(c)
			return
		}
		if fInfo.IsDir() {
			fileSearch := dirname + string(filepath.Separator) + "*"
			fileMatches, err := filepath.Glob(fileSearch)
			if err != nil {
				c <- StringError{"", err}
				close(c)
				return
			}
			for _, file := range fileMatches {
				fInfo, err := os.Stat(file)
				if err != nil {
					c <- StringError{"", err}
					continue
				}
				if fInfo.IsDir() {
					if ignoreDirs == false {
						c <- StringError{file, nil}
					}
					if recursive {
						d := GetFileList(file, ignoreDirs, recursive)
						for dirFile := range d {
							c <- dirFile
						}
					}
				} else {
					c <- StringError{file, nil}
				}
			}
		} else {
			c <- StringError{"", fmt.Errorf("Provided dir is not a dir: '%s'", dirname)}
			close(c)
			return
		}
		close(c)
	}()
	return c
}

// GetDirList returns a channel with each file (`channel.String`) or an error indicating failure (`channel.Error`).
func GetDirList(dir string) <-chan StringError {
	c := make(chan StringError)
	go func() {
		fInfo, err := os.Stat(dir)
		if err != nil {
			c <- StringError{"", err}
			return
		}
		if fInfo.IsDir() {
			c <- StringError{dir, nil}
			fileSearch := dir + string(filepath.Separator) + "*"
			fileMatches, err := filepath.Glob(fileSearch)
			if err != nil {
				c <- StringError{"", err}
				return
			}
			for _, file := range fileMatches {
				if filepath.Base(dir) == filepath.Base(file) {
					continue
				}
				d := GetDirList(file)
				for dirFile := range d {
					if dirFile.Error != nil {
						return
					}
					c <- StringError{dirFile.String, nil}
				}
			}
		}
		close(c)
	}()
	return c
}
