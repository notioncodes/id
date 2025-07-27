package id

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValidDashed(t *testing.T) {
	parsed, err := NewIDParser(NewNoOpCache()).Parse(TestIDs["valid-dashed"])
	assert.NoError(t, err)
	assert.Equal(t, parsed, TestIDs["valid-dashed"])
}

func TestParseValidUndashed(t *testing.T) {
	parsed, err := NewIDParser(NewNoOpCache()).Parse(TestIDs["valid-undashed"])
	assert.NoError(t, err)
	assert.Equal(t, parsed, TestIDs["valid-dashed"])
}

func TestParseInvalidChar(t *testing.T) {
	_, err := NewIDParser(NewNoOpCache()).Parse(TestIDs["invalid-char"])
	assert.Error(t, err)
}
