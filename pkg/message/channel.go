package message

import (
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/sunet/tq/pkg/utils"
)

type MessageChannelSource func(Message) (*MessageChannel, error)

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
	Quit   chan bool
	C      chan Message
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
		C:      make(chan Message),
		Quit:   make(chan bool),
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

func (channel *MessageChannel) ID() uint64 {
	return channel.id
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
	inputs := channel.Inputs()
	for id := range inputs {
		inputs[id].Close()
	}
	channel.Quit <- true
	defer close(channel.C)
}

func (channel *MessageChannel) Inputs() []*MessageChannel {
	cs := make([]*MessageChannel, len(channel.inputs))
	for i, id := range channel.inputs {
		cs[i] = Channels[id]
	}
	return cs
}

func (channel *MessageChannel) SentCount() int {
	return channel.nsent
}

func (channel *MessageChannel) RecvCount() int {
	return channel.nrecv
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
	dst.C <- o
}

func (out *MessageChannel) AddInput(in *MessageChannel) {
	out.inputs = append(out.inputs, in.id)
}

func (src *MessageChannel) Sink() {
	for {
		select {
		case <-src.C:
		}
	}
}

func (dst *MessageChannel) Consume(src *MessageChannel) {
	for v := range src.C {
		src.nsent++
		dst.Send(v)
	}
}

func (src *MessageChannel) Recv() Message {
	v := <-src.C
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
	for o := range in.C {
		in.nrecv++
		for _, c := range out {
			c.Send(o)
		}
	}
}

func ConnectChannels(source MessageSource, sink MessageSink, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	go func(out *MessageChannel) {
		o, err := source()
		if err != nil {
			Log.Error(err)
		} else {
			out.Send(o)
		}
	}(out)

	for _, c := range cs {
		go func(in *MessageChannel) {
			out.AddInput(in)
			in.final = false
			for v := range in.C {
				in.nrecv++
				err := sink(v)
				if err != nil {
					Log.Error(err)
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

func ConsumeChannels(h MessageChannelSource, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	for _, c := range cs {
		go func(in *MessageChannel) {
			out.AddInput(in)
			in.final = false
			for v := range in.C {
				in.nrecv++
				src, err := h(v)
				if err != nil {
					Log.Error(err)
				} else {
					out.Consume(src)
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

func ProcessChannels(h MessageHandler, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	for _, c := range cs {
		go func(in *MessageChannel) {
			out.AddInput(in)
			in.final = false
			for v := range in.C {
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
