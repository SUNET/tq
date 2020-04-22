package pipeline

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"

	"github.com/sunet/tq/pkg/message"
)

func MakeScriptPipeline(cmdline ...string) Pipeline {
	cmd_str := strings.Join(cmdline, " ")
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			Log.Debugf("about to run script %s", cmd_str)
			cmd := exec.Command(cmdline[0], cmdline[1:]...)
			var wg sync.WaitGroup
			wg.Add(3)

			stdin, err := cmd.StdinPipe()
			if err != nil {
				Log.Errorf("script %s: unable to connect to stdin for %s:", cmd_str, err.Error())
				return nil, err
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				Log.Errorf("script %s: unable to connect to stderr: %s", cmd_str, err.Error())
				return nil, err
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				Log.Errorf("script %s: unable to connect to stdout: %s", cmd_str, err.Error())
				return nil, err
			}

			err = cmd.Start()
			go func() {
				defer stdin.Close()
				data, err := message.FromJson(o)
				if err != nil {
					Log.Errorf("script %s: unable to serialize json: %s", cmd_str, err.Error())
				} else {
					Log.Debug(string(data))
					stdin.Write(data)
				}
				wg.Done()
			}()

			go func() {
				defer stderr.Close()
				e := bufio.NewReader(stderr)
				for {
					str, err := e.ReadString('\n')
					if err != nil {
						break
					}
					Log.Error(str)
				}
				wg.Done()
			}()

			var out []byte
			go func() {
				defer stdout.Close()
				e := bufio.NewReader(stdout)
				out, err = ioutil.ReadAll(e)
				if err != nil {
					Log.Errorf("script %s: unable to read stdout: %s", cmd_str, err.Error())
				}
				wg.Done()
			}()

			err = cmd.Wait()
			if err != nil {
				Log.Errorf("script %s: failed with %s\n", cmd_str, err.Error())
				return nil, err
			}

			wg.Wait()

			return message.ToJson(out)
		}, fmt.Sprintf("script %s", cmd_str), cs...)
	}
}
