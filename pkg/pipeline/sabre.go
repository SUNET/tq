package pipeline

import (
	"runtime"

	"github.com/spy16/sabre"
)

func SabreScope() sabre.Scope {
	scope := sabre.New()

	_ = scope.BindGo("merge", Merge)
	_ = scope.BindGo("log", LogMessages)
	_ = scope.BindGo("wait", WaitForAll)
	_ = scope.BindGo("print", func(o interface{}) { Log.Printf("%s", o) })
	_ = scope.BindGo("script", MakeScriptPipeline)
	_ = scope.BindGo("pub", MakePublishPipeline)
	_ = scope.BindGo("sub", MakeSubscribePipeline)
	_ = scope.BindGo("rest", MakeRESTPipeline)
	_ = scope.BindGo("timer", MakeTimer)
	_ = scope.BindGo("sink", SinkChannel)
	_ = scope.BindGo("recv", RecvMessage)
	_ = scope.BindGo("ngr", func() int { return runtime.NumGoroutine() })

	return scope
}
