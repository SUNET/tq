package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	_, err := Jsonf("{\"test: 1}")
	assert.Error(t, err, "error from Jsonf on malformed json")

	o, err := Jsonf("{\"test\": 1}")
	assert.Equal(t, err, nil, "no error from Jsonf")
	b, err := FromJson(o)
	assert.NoError(t, err, "marshall json works")
	assert.NotEmpty(t, b, "marshal returns nonemtpty value")
	o2, err := ToJson(b)
	assert.NoError(t, err, "unmarshal(marshal) json works")
	assert.NotNil(t, o2, "unmarshal returns value")
	b2, _ := FromJson(o2)
	assert.JSONEq(t, string(b2), string(b), "unmarshal(marshal) == id")
}
