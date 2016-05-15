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

func TestGetDirList(t *testing.T) {
	cases := []struct {
		dir    string
		result []string
	}{
		{"./test_tree", []string{
			"./test_tree",
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
