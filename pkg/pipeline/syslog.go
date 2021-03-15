package pipeline

import (
	"fmt"
	"github.com/sunet/tq/pkg/utils"
	"log"

	"github.com/sunet/tq/pkg/message"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func MakeSyslogUDP(hostport string) Pipeline {
	return MakeSyslog(hostport, func(server *syslog.Server) {
		server.ListenUDP(hostport)
	})
}

func MakeSyslogTCP(hostport string) Pipeline {
	return MakeSyslog(hostport, func(server *syslog.Server) {
		server.ListenTCP(hostport)
	})
}

func MakeSyslogTCPTLS(hostport string) Pipeline {
	return MakeSyslog(hostport, func(server *syslog.Server) {
		cfg := utils.GetTLSConfig()
		server.ListenTCPTLS(hostport, cfg)
	})
}

func MakeSyslogUnix(hostport string) Pipeline {
	return MakeSyslog(hostport, func(server *syslog.Server) {
		server.ListenUnixgram(hostport)
	})
}

func MakeSyslog(hostport string, setup func (server *syslog.Server)) Pipeline {
	out := message.NewMessageChannel(fmt.Sprintf("syslog"))
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	server.SetHandler(handler)

	setup(server)

	server.SetFormat(syslog.Automatic)

	server.Boot()

	return func(...*message.MessageChannel) *message.MessageChannel {

		go func(channel syslog.LogPartsChannel) {
			defer out.Close()
			defer server.Kill()

			for logParts := range channel {
				log.Printf("%s", logParts)
				out.Send(message.Message(logParts))
			}
		}(channel)

		return out
	}
}
