package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	m := map[string]interface{}{"x": 1, "y": 2}
	k := Keys(m)
	assert.Contains(t, k, "x", "keys contains expected values")
	assert.Contains(t, k, "y", "keys contains expected values")
}

func TestFuncionName(t *testing.T) {
	assert.Contains(t, GetFunctionName(GetFunctionName), "GetFunctionName")
}
