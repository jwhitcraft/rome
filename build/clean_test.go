package build

import (
	"testing"
	"errors"
)

type MockTargetDir struct {
	Path string
	Err error
}

func (t MockTargetDir) Clean(path string) error {
	return t.Err
}

func (t MockTargetDir) Dir() string {
	return t.Path
}

func TestCleanBuild(t *testing.T) {
	dir := MockTargetDir{Path: "/tmp", Err: nil}

	actual := CleanBuild(dir)
	if actual != nil {
		t.Errorf("Expected Nil to Be Returned")
	}

	// Test if an error is returned
	errDir := MockTargetDir{Path: "/tmp", Err: errors.New("Invalid Dir")}
	actual = CleanBuild(errDir)
	if actual == nil {
		t.Errorf("Expected An Error to Be Returned")
	}
}