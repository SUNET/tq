package message

import (
	"fmt"
	"sync"

	"github.com/SUNET/tq/pkg/utils"
	"github.com/sirupsen/logrus"
)

type MessageChannel struct {
	wg   sync.WaitGroup
	c    chan Message
	name string
}

func NewMessageChannel(name string, sz ...int) *MessageChannel {
	size := 1
	if len(sz) > 0 {
		size = sz[0]
	}
	p := MessageChannel{c: make(chan Message), name: name}
	p.wg.Add(size)
	return &p
}

func (channel *MessageChannel) String() string {
	return channel.name
}

func (channel *MessageChannel) Wait() {
	channel.wg.Wait()
}

func (channel *MessageChannel) Done() {
	channel.wg.Done()
}

func (channel *MessageChannel) Close() {
	close(channel.c)
}

func (dst *MessageChannel) Send(o Message) {
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := o.String()
		if err == nil {
			Log.Debugf("[OUT] %s => %s", s, dst.name)
		}
	}
	dst.c <- o
}

func (src *MessageChannel) Sink() {
	for {
		select {
		case <-src.c:
		}
	}
}

func (src *MessageChannel) Recv() Message {
	v := <-src.c
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := v.String()
		if err == nil {
			Log.Debugf("[IN] %s <= %s", src.name, s)
		}
	}
	return v
}

func ProcessChannels(h MessageHandler, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(fmt.Sprintf("handler [%s]", utils.GetFunctionName(h)), len(cs))
	for _, c := range cs {
		go func(in *MessageChannel) {
			for v := range in.c {
				out.Send(h(v))
			}
			out.Done()
		}(c)
	}
	go func() {
		out.Wait()
		out.Close()
	}()
	return out
}
