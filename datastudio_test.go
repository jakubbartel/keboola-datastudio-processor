package kbcdatastudioproc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleInfinityRow(t *testing.T) {
	data := []string{"INFINITY"}

	_, err := encodeColumnNumber(data)
	assert.Error(t, err, "Number column can not contain infinity floats")
}
