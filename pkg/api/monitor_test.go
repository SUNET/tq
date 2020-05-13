package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
)

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(status)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK, "status returns 200")
	o, err := message.ToJson([]byte(rr.Body.String()))
	assert.Contains(t, utils.Keys(o), "number-of-gorutines",
		"status contains number-of-goroutines")
	assert.Contains(t, utils.Keys(o), "version",
		"status contains version")
}

func TestAllChannelHandler(t *testing.T) {
	ch := message.NewMessageChannel("test", 1)
	req, err := http.NewRequest("GET", "/channels/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ch.IsFinal())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(all_channels)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK, "channels all returns 200")
	var o []message.Message = make([]message.Message, 1)
	err = json.Unmarshal([]byte(rr.Body.String()), &o)
	assert.NoError(t, err, "unmarshal response")
	assert.Equal(t, len(o), 1, "return 1 channel")
}

func TestListChannelHandler(t *testing.T) {
	ch := message.NewMessageChannel("test", 1)
	req, err := http.NewRequest("GET", "/channels/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ch.IsFinal())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(list_channels)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK, "channels list returns 200")
	var o []string = make([]string, 1)
	err = json.Unmarshal([]byte(rr.Body.String()), &o)
	assert.NoError(t, err, "unmarshal response")
	assert.Equal(t, len(o), 2, "return 2 channels")
}
