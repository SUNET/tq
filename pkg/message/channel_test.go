package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeChannel(t *testing.T) {
	ch := NewMessageChannel("test", 1)
	assert.True(t, ch.final, "new channel is final")
	assert.Zero(t, ch.nrecv, "none received")
	assert.Zero(t, ch.nsent, "none sent")

	msg, _ := Jsonf("{\"test\": 1}")

	go func() { ch.Send(msg); ch.Close() }()
	msg2 := ch.Recv()
	assert.Equal(t, ch.nsent, 1, "one sent")
	assert.Equal(t, ch.nrecv, 1, "one recv")
	assert.Equal(t, msg, msg2, "got what I sent")
}

func TestChannelDB(t *testing.T) {
	assert.Equal(t, len(Channels), 1, "one channel created so far")
}
