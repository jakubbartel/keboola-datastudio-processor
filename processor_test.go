package kbcdatastudioproc

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchOutput(t *testing.T) {
	p := matchOutputPath(
		filepath.FromSlash("/data/in/dir/a.csv"), filepath.FromSlash("/data/in/"), filepath.FromSlash("/data/out/"))
	assert.Equal(t, filepath.FromSlash("/data/out/dir/a.datastudio"), p, "absolute paths")
}

func TestListInputs(t *testing.T) {
	inputs, err := listInputs("test/fixtures/map_inputs")
	if assert.NoError(t, err, "unexpected list inputs error") {
		assert.ElementsMatch(t, inputs, []string{
			filepath.FromSlash("test/fixtures/map_inputs/dir/c.csv"),
			filepath.FromSlash("test/fixtures/map_inputs/a.csv"),
			filepath.FromSlash("test/fixtures/map_inputs/b.csv"),
		}, "all csv files presented")
	}
}

func TestProcessDir(t *testing.T) {
	err := processDir(
		filepath.FromSlash("test/fixtures/data_1/in/files"), filepath.FromSlash("test/fixtures/data_1/out/files"))
	if assert.NoError(t, err, "unexpected process dir error") {
		assert.FileExists(t, "test/fixtures/data_1/out/files/a.datastudio", "csv processed")
		assert.FileExists(t, "test/fixtures/data_1/out/files/dir/b.datastudio", "nested csv processed")
	}
}

func clearDir(dirPath string) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic("read test dir for cleanup")
	}

	for _, d := range dir {
		err := os.RemoveAll(path.Join(dirPath, d.Name()))
		if err != nil {
			panic("clear test dir")
		}
	}
}

func TestMain(m *testing.M) {
	clearDir("test/fixtures/data_1/out/files")

	code := m.Run()

	clearDir("test/fixtures/data_1/out/files")

	os.Exit(code)
}
