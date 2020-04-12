package kbcdatastudioproc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jakubbartel/kbcdatastudioproc/keboola"
)

var ErrUser = errors.New("user error")

// Generate filesystem path to the output file. Assumes input file is in the input dir and has .csv suffix.
func matchOutputPath(path, inputDir, outputDir string) string {
	path = strings.TrimPrefix(path, inputDir)
	path = strings.TrimSuffix(path, ".csv")

	outputDir = strings.TrimRight(outputDir, string(os.PathSeparator)) + string(os.PathSeparator)

	return outputDir + path + ".datastudio"
}

func listInputs(dirPath string) ([]string, error) {
	matches := make([]string, 0)

	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list directory: %w", err)
	}

	return matches, nil
}

func processFile(path, inputDir, outputDir string) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer r.Close()

	outputPath := matchOutputPath(path, inputDir, outputDir)

	dir := filepath.Dir(outputPath)

	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return fmt.Errorf("make output dir: %w", err)
	}

	w, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create outut file: %w", err)
	}
	defer w.Close()

	if err := EncodeCsv(r, w); err != nil {
		return fmt.Errorf("encode input file to output: %w", err)
	}

	return nil
}

func processDir(dirPath, outputDir string) error {
	dirPath = filepath.FromSlash(dirPath)
	outputDir = filepath.FromSlash(outputDir)

	files, err := listInputs(dirPath)
	if err != nil {
		return fmt.Errorf("list inputs: %w", err)
	}

	for _, file := range files {
		if err := processFile(file, dirPath, outputDir); err != nil {
			return fmt.Errorf("process file: %w", err)
		}
	}

	return nil
}

func RunE() error {
	if err := processDir(keboola.InFilesDir, keboola.OutFilesDir); err != nil {
		return fmt.Errorf("process files: %w", err)
	}

	if err := processDir(keboola.InTablesDir, keboola.OutTablesDir); err != nil {
		return fmt.Errorf("process tables: %w", err)
	}

	return nil
}
