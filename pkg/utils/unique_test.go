package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqeID(t *testing.T) {
	id, err := UniqueID()
	assert.NoError(t, err, "uniqeid not erroring")
	assert.True(t, id > 0, "unique id > 0")
}
