package helputils

import (
	"path/filepath"
	"testing"
)

func TestFindFilesByPattern(t *testing.T) {
	workingDir, err := GetWd()
	if err != nil {
		t.Fatal("Failed to get working directory. Err: " + err.Error())
	}
	root := filepath.Join(workingDir, "..", "..", "..", "..", "..")
	pattern := "(vc-15|vc15).*.toml"
	t.Log("looking for", pattern, "under", root, "...")
	files, err := FindFilesByPattern(root, pattern)
	if err != nil {
		t.Errorf("Error finding files: %v", err)
	}
	t.Log("there are", len(files), pattern, "files under", root)
}

func TestWdRel(t *testing.T) {
	path, err := WdRel("/home/")
	if err != nil {
		t.Error("Path is: " + path)
	}
}
