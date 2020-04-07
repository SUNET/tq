package pipeline

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/SUNET/tq/pkg/message"
	sh "github.com/kballard/go-shellquote"
)

var makeScriptPipeline = func(cmd string) Pipeline {
	cmdline, err := sh.Split(cmd)
	if err != nil {
		Log.Fatalf("Failed to parse string: %s\n", err)
	}
	return NewPipelineFromHandler(func(o message.Message) message.Message {
		cmd := exec.Command(cmdline[0], cmdline[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		env := os.Environ()
		for k, v := range o {
			env = append(env, fmt.Sprintf("DOIT_%v=%v", k, v))
		}
		cmd.Env = env
		err := cmd.Run()
		if err != nil {
			Log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		m := make(message.Message)
		m["ok"] = 1
		return m
	})
}

var makeLogPipeline = func(_ string) Pipeline {
	return NewPipelineFromHandler(func(o message.Message) message.Message {
		Log.Print(o)
		return nil
	})
}

func init() {
	Register("log", makeLogPipeline)
	Register("script", makeScriptPipeline)
}
