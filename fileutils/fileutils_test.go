package fileutils

import (
	"reflect"
	"testing"
)

func TestGetFileList(t *testing.T) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.A/b/C/d/E",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/A/b/C/d/E",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
			"test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"test_tree/.A/b/C/d/E",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A/b/C/d/E",
			"test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
		},
		},
	}
	for _, c := range cases {
		ch := GetFileList(c.file, c.ignoreDir, c.recursive)
		tree := []string{}
		for e := range ch {
			if e.Error != nil {
				t.Fatalf("Unexpected error: %s\n", e.Error)
			}
			tree = append(tree, e.String)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Fatalf("tree %q != %q", c.result, tree)
		}
	}
}

func BenchmarkGetFileList(b *testing.B) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.A/b/C/d/E",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/A/b/C/d/E",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
			"test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"test_tree/.A/b/C/d/E",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A/b/C/d/E",
			"test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
		},
		},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			ch := GetFileList(c.file, c.ignoreDir, c.recursive)
			tree := []string{}
			for e := range ch {
				tree = append(tree, e.String)
			}
		}
	}
}

func TestListFiles(t *testing.T) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"./test_tree/.A",
			"./test_tree/.A/b",
			"./test_tree/.A/b/C",
			"./test_tree/.A/b/C/d",
			"./test_tree/.A/b/C/d/E",
			"./test_tree/.a",
			"./test_tree/.a/B",
			"./test_tree/.a/B/c",
			"./test_tree/.a/B/c/D",
			"./test_tree/.a/B/c/D/e",
			"./test_tree/.svn",
			"./test_tree/.svn/E",
			"./test_tree/.svn/e",
			"./test_tree/A",
			"./test_tree/A/b",
			"./test_tree/A/b/C",
			"./test_tree/A/b/C/d",
			"./test_tree/A/b/C/d/E",
			"./test_tree/a",
			"./test_tree/a/B",
			"./test_tree/a/B/c",
			"./test_tree/a/B/c/D",
			"./test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"./test_tree/.A/b/C/d/E",
			"./test_tree/.a/B/c/D/e",
			"./test_tree/.svn/E",
			"./test_tree/.svn/e",
			"./test_tree/A/b/C/d/E",
			"./test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"./test_tree/.A",
			"./test_tree/.a",
			"./test_tree/.svn",
			"./test_tree/A",
			"./test_tree/a",
		},
		},
	}
	for _, c := range cases {
		tree, err := ListFiles(c.file, c.ignoreDir, c.recursive)
		if err != nil {
			t.Fatalf("Unexpected error: %s\n", err)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Fatalf("tree %q != %q", c.result, tree)
		}
	}
}

func BenchmarkListFiles(b *testing.B) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.A/b/C/d/E",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/A/b/C/d/E",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
			"test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"test_tree/.A/b/C/d/E",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A/b/C/d/E",
			"test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
		},
		},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			list, _ := ListFiles(c.file, c.ignoreDir, c.recursive)
			tree := []string{}
			for _, e := range list {
				tree = append(tree, e)
			}
		}
	}
}

func TestGetNumSortFileList(t *testing.T) {
	cases := []struct {
		dir       string
		ignoreDir bool
		recursive bool
		reverse   bool
		result    []string
	}{
		{"./test_tree2", false, true, false, []string{
			"test_tree2/1",
			"test_tree2/2",
			"test_tree2/3",
			"test_tree2/10",
			"test_tree2/20",
			"test_tree2/30",
		},
		},
		{"./test_tree2", false, true, true, []string{
			"test_tree2/30",
			"test_tree2/20",
			"test_tree2/10",
			"test_tree2/3",
			"test_tree2/2",
			"test_tree2/1",
		},
		},
	}
	for _, c := range cases {
		ch := GetNumSortFileList(c.dir, c.ignoreDir, c.recursive, c.reverse)
		tree := []string{}
		for e := range ch {
			tree = append(tree, e.String)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Errorf("tree %q != %q", c.result, tree)
		}
	}
}

func BenchmarkGetNumSortFileList(b *testing.B) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.A/b/C/d/E",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/A/b/C/d/E",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
			"test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"test_tree/.A/b/C/d/E",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A/b/C/d/E",
			"test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
		},
		},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			ch := GetNumSortFileList(c.file, c.ignoreDir, c.recursive, false)
			tree := []string{}
			for e := range ch {
				tree = append(tree, e.String)
			}
		}
	}
}

func TestListFilesNumSort(t *testing.T) {
	cases := []struct {
		dir       string
		ignoreDir bool
		recursive bool
		reverse   bool
		result    []string
	}{
		{"./test_tree2", false, true, false, []string{
			"./test_tree2/1",
			"./test_tree2/2",
			"./test_tree2/3",
			"./test_tree2/10",
			"./test_tree2/20",
			"./test_tree2/30",
		},
		},
		{"./test_tree2", false, true, true, []string{
			"./test_tree2/30",
			"./test_tree2/20",
			"./test_tree2/10",
			"./test_tree2/3",
			"./test_tree2/2",
			"./test_tree2/1",
		},
		},
	}
	for _, c := range cases {
		tree, err := ListFilesNumSort(c.dir, c.ignoreDir, c.recursive, c.reverse)
		if err != nil {
			t.Fatalf("Unexpected error: %s\n", err)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Errorf("tree %q != %q", c.result, tree)
		}
	}
}

func BenchmarkListFilesNumSort(b *testing.B) {
	cases := []struct {
		file      string
		ignoreDir bool
		recursive bool
		result    []string
	}{
		{"./test_tree", false, true, []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.A/b/C/d/E",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/A/b/C/d/E",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
			"test_tree/a/B/c/D/e"},
		},
		{"./test_tree", true, true, []string{
			"test_tree/.A/b/C/d/E",
			"test_tree/.a/B/c/D/e",
			"test_tree/.svn/E",
			"test_tree/.svn/e",
			"test_tree/A/b/C/d/E",
			"test_tree/a/B/c/D/e",
		},
		},
		{"./test_tree", true, false, []string{}},
		{"./test_tree", false, false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
		},
		},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			list, _ := ListFilesNumSort(c.file, c.ignoreDir, c.recursive, false)
			tree := []string{}
			for _, e := range list {
				tree = append(tree, e)
			}
		}
	}
}

func TestGetDirList(t *testing.T) {
	cases := []struct {
		dir    string
		result []string
	}{
		{"./test_tree", []string{
			"test_tree/.A",
			"test_tree/.A/b",
			"test_tree/.A/b/C",
			"test_tree/.A/b/C/d",
			"test_tree/.a",
			"test_tree/.a/B",
			"test_tree/.a/B/c",
			"test_tree/.a/B/c/D",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/A/b",
			"test_tree/A/b/C",
			"test_tree/A/b/C/d",
			"test_tree/a",
			"test_tree/a/B",
			"test_tree/a/B/c",
			"test_tree/a/B/c/D",
		},
		},
	}
	for _, c := range cases {
		ch := GetDirList(c.dir)
		tree := []string{}
		for e := range ch {
			tree = append(tree, e.String)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Errorf("tree %q != %q", c.result, tree)
		}
	}
}

func TestGetNumSortDirList(t *testing.T) {
	cases := []struct {
		dir     string
		reverse bool
		result  []string
	}{
		{"./test_tree2", false, []string{
			"test_tree2/1",
			"test_tree2/2",
			"test_tree2/3",
			"test_tree2/10",
			"test_tree2/20",
			"test_tree2/30",
		},
		},
		{"./test_tree2", true, []string{
			"test_tree2/30",
			"test_tree2/20",
			"test_tree2/10",
			"test_tree2/3",
			"test_tree2/2",
			"test_tree2/1",
		},
		},
	}
	for _, c := range cases {
		ch := GetNumSortDirList(c.dir, c.reverse)
		tree := []string{}
		for e := range ch {
			tree = append(tree, e.String)
		}
		if !reflect.DeepEqual(tree, c.result) {
			t.Errorf("tree %q != %q", c.result, tree)
		}
	}
}

func TestStringReplace(t *testing.T) {
	n, err := StringReplace("test_tree/A/b/C/d/E", "lorem", "hello", -1, 1024)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if n != 2 {
		t.Fatalf("Unexpected amount of lines changed: %d\n", n)
	}
	n, err = StringReplace("test_tree/A/b/C/d/E", "hello", "lorem", -1, 1024)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if n != 2 {
		t.Fatalf("Unexpected amount of lines changed: %d\n", n)
	}
}
