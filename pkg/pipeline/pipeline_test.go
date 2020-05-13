package pipeline

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sunet/tq/pkg/message"
)

func runPipeline(p Pipeline, duration string) *message.MessageChannel {
	d, err := time.ParseDuration(duration)
	if err != nil {
		Log.Panicf("Unable to parse duration: %s", err.Error())
	}

	var ch *message.MessageChannel
	go func() { ch = p(); SinkChannel(ch) }()
	time.Sleep(d)
	ch.Close()

	return ch
}

func TestTimerPipeline(t *testing.T) {
	p := MakeTimer("1s")
	ch := runPipeline(p, "3s")
	log.Print(ch.SentCount())
	assert.True(t, ch.SentCount() >= 1, "sent at least 1 messages")
}
