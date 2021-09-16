package main

import (
	"github.com/spy16/slurp"
	"github.com/spy16/slurp/builtin"
	"github.com/spy16/slurp/core"
	"github.com/sunet/tq/pkg/api"
	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/pipeline"
)

func NewInterpreter() *slurp.Interpreter {
	sl := slurp.New()
	var globals = map[string]core.Any {
		"nil":     		builtin.Nil{},
		"true":    		builtin.Bool(true),
		"false":   		builtin.Bool(false),

		"tq/merge": 	slurp.Func("tq/merge", pipeline.Merge),
		"tq/log": 		slurp.Func("tq/log", pipeline.LogMessages),
		"tq/run": 		slurp.Func("tq/run", pipeline.Run),
		"tq/fork": 		slurp.Func("tq/fork", pipeline.ForkAndMerge),
		"tq/timer": 	slurp.Func("tq/timer", pipeline.MakeTimer),
		"tq/listen":	slurp.Func("tq/listen", api.Listen),
		"tq/load": 		slurp.Func("tq/load", func(files ...string) { readEvalFiles(sl, files... )}),

		"syslog/udp": 	slurp.Func("syslog/udp", pipeline.MakeSyslogUDP),
		"syslog/tcp": 	slurp.Func("syslog/tcp", pipeline.MakeSyslogTCP),
		"syslog/tcptls": slurp.Func("syslog/tcptls", pipeline.MakeSyslogTCPTLS),
		"syslog/unix": 	slurp.Func("syslog/unix", pipeline.MakeSyslogUnix),

		"script/pipeline": slurp.Func("script/pipeline", pipeline.MakeScriptPipeline),
		"script/handler": 	slurp.Func("script/handler", message.ScriptHandler),

		"filter/eq": 	slurp.Func("filter/eq", pipeline.MakeEQFilter),
		"filter/any": 	slurp.Func("filter/any", pipeline.MakeMatchAnyFilter),

		"kazaam/pipeline": slurp.Func("kazaam/pipeline", pipeline.MakeKazaamPipeline),
		"kazaam/handler": slurp.Func("kazaam/handler", message.KazaamHandler),
		"kazaam/parse": slurp.Func("kazaam/parse", message.NewKazaam),

		"nanomsg/pub": slurp.Func("nanomsg/pub", pipeline.MakePublishPipeline),
		"nanomsg/sub": slurp.Func("nanomsg/sub", pipeline.MakeSubscribePipeline),
		"nanomsg/pull": slurp.Func("nanomsg/pull", pipeline.MakePullPipeline),
		"nanomsg/push": slurp.Func("nanomsg/push", pipeline.MakePushPipeline),
		"nanomsg/surveyor": slurp.Func("nanomsg/surveyor", pipeline.MakeSurveyorPipeline),
		"nanomsg/respondent": slurp.Func("nanomsg/respondent", pipeline.MakeRespondentPipeline),

		"http/source": slurp.Func("http/source", pipeline.MakeHTTPSourcePipeline),
		"http/sourceresponse": slurp.Func("http/sourceresponse", pipeline.MakeHTTPSourceResponsePipeline),
		"http/sink": slurp.Func("http/sink", pipeline.MakePOSTPipeline),
		"http/post": slurp.Func("http/post", message.PostHandler),
	}

	if err := sl.Bind(globals); err != nil {
		Log.Panic(err)
	}

	return sl
}
