package kbcdatastudioproc

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type EncodingType string

const (
	EncodingTypeRaw    = "raw"
	EncodingTypeNumber = "number"
	EncodingTypeMap    = "map"
	//EncodingTypeNumericMap = "numericMap"
)

type Dataset struct {
	Columns []*Column `json:"columns"`
	RowsNum int       `json:"rowsNum"`
}

type Column struct {
	Name     string       `json:"name"`
	Encoding EncodingType `json:"encoding"`
	Strings  []string     `json:"strings,omitempty"`
	Numbers  []float32    `json:"numbers,omitempty"`
	Indexes  []int        `json:"indexes,omitempty"`
	Values   []string     `json:"values,omitempty"`
}

func encodeColumn(enc EncodingType, data []string) (*Column, error) {
	switch enc {
	case EncodingTypeRaw:
		return encodeColumnRaw(data)
	case EncodingTypeNumber:
		return encodeColumnNumber(data)
	case EncodingTypeMap:
		return encodeColumnMap(data)
	default:
		return nil, fmt.Errorf("unsupported encoding")
	}
}

func encodeColumnRaw(data []string) (*Column, error) {
	return &Column{
		Encoding: EncodingTypeRaw,
		Strings:  data,
	}, nil
}

func encodeColumnNumber(data []string) (*Column, error) {
	nums := make([]float32, 0)

	for _, str := range data {
		num, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return nil, fmt.Errorf("column contains a non-numeric string: %v", err)
		}

		nums = append(nums, float32(num))
	}

	return &Column{
		Encoding: EncodingTypeNumber,
		Numbers:  nums,
	}, nil
}

func encodeColumnMap(data []string) (*Column, error) {
	i := 0
	mapping := map[string]int{}

	idxs := make([]int, 0)
	vals := make([]string, 0)

	for _, str := range data {
		idx, ok := mapping[str]
		if !ok {
			idx = i
			mapping[str] = i
			vals = append(vals, str)
			i += 1
		}

		idxs = append(idxs, idx)
	}

	return &Column{
		Encoding: EncodingTypeMap,
		Indexes:  idxs,
		Values:   vals,
	}, nil
}

// The lower the better. Calculates gzipped bytes on json encoded column.
func scoreColumn(col *Column) float64 {
	b := bytes.NewBuffer([]byte{})

	gw := gzip.NewWriter(b)
	jw := json.NewEncoder(gw)

	if err := jw.Encode(col); err != nil {
		// TODO: handle more systematically
		_, _ = fmt.Fprintf(os.Stderr, "Cannot json encode column: %v", err)
		os.Exit(1)
	}

	_ = gw.Close()

	return float64(b.Len())
}

// Process a column with all defined encodings and pick the one with the best score.
func autoEncodeColumn(data []string) *Column {
	encodings := []EncodingType{
		EncodingTypeRaw,
		EncodingTypeNumber,
		EncodingTypeMap,
	}

	var col *Column
	var score float64

	for _, enc := range encodings {
		// Ignore errors because not all encodings are valid on a column.
		encCol, err := encodeColumn(enc, data)
		if err != nil {
			continue
		}

		colScore := scoreColumn(encCol)

		//fmt.Printf("Column scored to %f by %s encoder\n", colScore, enc)

		if score == 0.0 || colScore < score {
			col = encCol
			score = colScore
		}
	}

	return col
}

// Reader with csv input, writer for encoded datastudio format.
func EncodeCsv(r io.Reader, w io.Writer) error {
	cr := csv.NewReader(r)

	data, err := cr.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv: %v", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("empty input dataset")
	}

	dataset := Dataset{
		Columns: make([]*Column, 0),
	}

	// expects headers row
	for i, columnName := range data[0] {
		column := make([]string, 0)

		for j := 1; j < len(data); j++ {
			column = append(column, data[j][i])
		}

		col := autoEncodeColumn(column)

		col.Name = columnName

		dataset.Columns = append(dataset.Columns, col)
		dataset.RowsNum = len(data) - 1
	}

	gw, _ := gzip.NewWriterLevel(w, gzip.DefaultCompression)

	if err := json.NewEncoder(gw).Encode(dataset); err != nil {
		return fmt.Errorf("json encode dataset: %v", err)
	}

	if err := gw.Close(); err != nil {
		return fmt.Errorf("gzip encode dataset: %v", err)
	}

	return nil
}
