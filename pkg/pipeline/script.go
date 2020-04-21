package pipeline

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/SUNET/tq/pkg/message"
	sh "github.com/kballard/go-shellquote"
)

func MakeScriptPipeline(cmd string) Pipeline {
	cmdline, err := sh.Split(cmd)
	if err != nil {
		Log.Fatalf("Failed to parse string: %s\n", err)
	}

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
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
				Log.Errorf("cmd.Run() failed with %s\n", err)
				return nil, err
			}
			//TODO - encode output
			m := make(message.Message)
			m["ok"] = 1
			return m, nil
		}, "script", cs...)
	}
}
