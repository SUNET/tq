package message

import (
	"encoding/json"
	"sync"

	"github.com/SUNET/tq/pkg/utils"
	"github.com/sirupsen/logrus"
)

type ChannelDB map[uint64]*MessageChannel

func AllFinalChannels() []*MessageChannel {
	nch := 0
	for _, ch := range Channels {
		if ch.IsFinal() {
			nch++
		}
	}
	chs := make([]*MessageChannel, nch)
	i := 0
	for _, ch := range Channels {
		if ch.IsFinal() {
			chs[i] = ch
			i++
		}
	}
	return chs
}

var Channels ChannelDB

type MessageChannel struct {
	id     uint64
	wg     sync.WaitGroup
	c      chan Message
	nrecv  int
	nsent  int
	name   string
	final  bool
	inputs []uint64
}

func NewMessageChannel(name string, sz ...int) *MessageChannel {
	size := 1
	if len(sz) > 0 {
		size = sz[0]
	}
	id, err := utils.UniqueID()
	if err != nil {
		Log.Panic(err.Error())
	}
	p := MessageChannel{
		c:      make(chan Message),
		inputs: make([]uint64, 0, 3),
		id:     id,
		final:  true,
		name:   name,
	}
	p.wg.Add(size)
	Channels[p.id] = &p
	return &p
}

func (channel *MessageChannel) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ID       uint64   `json:"id"`
		Name     string   `json:"name"`
		Final    bool     `json:"final"`
		Received int      `json:"received,omitempty"`
		Sent     int      `json:"sent,omitempty"`
		Inputs   []uint64 `json:"inputs,omitempty"`
	}{
		ID:       channel.id,
		Name:     channel.name,
		Final:    channel.final,
		Received: channel.nrecv,
		Sent:     channel.nsent,
		Inputs:   channel.inputs,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
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

func (channel *MessageChannel) IsFinal() bool {
	return channel.final
}

func (dst *MessageChannel) Send(o Message) {
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := o.String()
		if err == nil {
			Log.Debugf("[OUT] %s => %s", s, dst.name)
		}
	}
	dst.nsent++
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
	src.nrecv++
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := v.String()
		if err == nil {
			Log.Debugf("[IN] %s <= %s", src.name, s)
		}
	}
	return v
}

func Fork(in *MessageChannel, out ...*MessageChannel) {
	for o := range in.c {
		in.nrecv++
		for _, c := range out {
			c.Send(o)
		}
	}
}

func ProcessChannels(h MessageHandler, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	for _, c := range cs {
		go func(in *MessageChannel) {
			out.inputs = append(out.inputs, in.id)
			in.final = false
			for v := range in.c {
				in.nrecv++
				o, err := h(v)
				if err != nil {
					Log.Error(err)
				} else {
					out.Send(o)
				}
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

func init() {
	Channels = make(ChannelDB)
}
