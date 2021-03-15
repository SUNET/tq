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

func (channel *MessageChannel) Send(o Message) {
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := o.String()
		if err == nil {
			Log.Debugf("[OUT] %v => %v", s, channel)
		}
	}
	channel.nsent++
	channel.C <- o
}

func (channel *MessageChannel) AddInput(in *MessageChannel) {
	channel.inputs = append(channel.inputs, in.id)
}

func (channel *MessageChannel) Sink() {
	for {
		select {
		case <-channel.C:
		}
	}
}

func (channel *MessageChannel) Consume(src *MessageChannel) {
	for v := range src.C {
		src.nsent++
		channel.Send(v)
	}
}

func (channel *MessageChannel) Recv() Message {
	v := <-channel.C
	channel.nrecv++
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		s, err := v.String()
		if err == nil {
			Log.Debugf("[IN] %v <= %v", channel, s)
		}
	}
	return v
}

func ForkChannel(in *MessageChannel, out ...*MessageChannel) {
	for o := range in.C {
		in.nrecv++
		for _, c := range out {
			c.Send(o)
		}
	}
}

func (channel *MessageChannel) Process(sendToChannel MessageSender, cs ...*MessageChannel) *MessageChannel {
	for _, c := range cs {
		go func(in *MessageChannel) {
			channel.AddInput(in)
			in.final = false
			for v := range in.C {
				in.nrecv++
				sendToChannel(channel, v)
			}
			channel.Done()
		}(c)
	}
	go func() {
		channel.Wait()
		channel.Close()
	}()
	return channel
}

func ConsumeChannels(h MessageChannelSource, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	return out.Process(func(out *MessageChannel, m Message) {
		src, err := h(m)
		if err != nil {
			Log.Error(err)
		} else {
			out.Consume(src)
		}
	}, cs...)
}

func ProcessChannels(h MessageHandler, name string, cs ...*MessageChannel) *MessageChannel {
	out := NewMessageChannel(name, len(cs))
	return out.Process(func(out *MessageChannel, m Message) {
		o, err := h(m)
		if err != nil {
			Log.Error(err)
		} else {
			out.Send(o)
		}
	}, cs...)
}

func init() {
	Channels = make(ChannelDB)
}
