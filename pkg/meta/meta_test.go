package meta

import (
	"strings"
	"testing"
)

func TestMeta(t *testing.T) {
	got := Name()
	if len(got) <= 0 {
		t.Errorf("empty name")
	}
	if !strings.Contains(got, "tq") {
		t.Errorf("%s does not contain 'tq'", got)
	}
}
