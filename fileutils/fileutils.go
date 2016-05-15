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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
func GetDirList(dirname string) <-chan StringError {
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
					c <- StringError{file, nil}
					d := GetDirList(file)
					for dirFile := range d {
						c <- dirFile
						// if dirFile.Error != nil {
						// 	c <- StringError{"", dirFile.Error}
						// 	continue
						// }
						// c <- StringError{dirFile.String, nil}
					}
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

// StringReplace - Runs strings.Replace on each line of the file.
// The file is read line by line to account for large files.
// The changes are first written to a tmp copy is saved before overwriting the
// original. The original is only changed if linesChanged > 0.
func StringReplace(file, old, new string, n, bufferSize int) (int, error) {
	var tmpFile *os.File
	linesChanged := 0
	tmpFile, err := ioutil.TempFile("", filepath.Base(file)+"-")
	if err != nil {
		return 0, fmt.Errorf("cannot open '%s': %s\n", tmpFile.Name(), err)
	}
	defer tmpFile.Close()
	for d := range ReadLines(file, bufferSize) {
		if d.Error != nil {
			return 0, fmt.Errorf("Error reading file '%s': %s\n", file, d.Error)
		}
		line := strings.Replace(d.String, old, new, n)
		if d.String != line {
			linesChanged++
		}
		tmpFile.WriteString(line + "\n")
	}
	tmpFile.Close()
	if linesChanged > 0 {
		err = CopyFile(tmpFile.Name(), file)
		if err != nil {
			return 0, fmt.Errorf("Couldn't update file: %s. '%s'\n", file, err)
		}
	}
	return linesChanged, nil
}

// ReadLines - returns a channel of type StringError with each line of a file.
func ReadLines(filename string, bufferSize int) <-chan StringError {
	c := make(chan StringError)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			c <- StringError{"", fmt.Errorf("Couldn't open file '%s': %s\n", filename, err)}
			close(c)
			return
		}
		defer file.Close()

		reader := bufio.NewReaderSize(file, bufferSize)
		// line number
		n := 0
		for {
			n++
			line, isPrefix, err := reader.ReadLine()
			if isPrefix {
				c <- StringError{"", fmt.Errorf("%s: buffer size too small\n", filename)}
				break
			}
			// stop reading file
			if err != nil {
				if err != io.EOF {
					c <- StringError{"", fmt.Errorf("Read error '%s': %s\n", filename, err)}
				}
				break
			}
			c <- StringError{string(line), nil}
		}
		close(c)
	}()
	return c
}
