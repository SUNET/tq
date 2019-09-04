package tq

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	//"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/json"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	_ "github.com/SUNET/tq/monitor"
)

var Log = logrus.New()

func to_json(message []byte) (map[string]interface{}, error) {
	var o map[string]interface{}
	err := json.Unmarshal(message, &o)
	return o, err
}

func from_json(o map[string]interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func handle_message(o map[string]interface{}, cmdline []string) {
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
}

func Pub(url string) {
	var sock mangos.Socket
	var err error
	var data []byte
	if sock, err = pub.NewSocket(); err != nil {
		Log.Fatalf("can't get new pub socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		Log.Fatalf("can't listen on pub socket: %s", err.Error())
	}

	pubHandler := func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		data, err = ioutil.ReadAll(req.Body)
		Log.Printf("got request to publish data: %s", data)
		io.WriteString(w, "ok\n")
		if err = sock.Send(data); err != nil {
			Log.Fatalf("Failed publishing: %s", err.Error())
		}
	}
	http.HandleFunc("/publish", pubHandler)
	Log.Fatal(http.ListenAndServe(":8080", nil))
}

func Sub(url string, cmdline []string) {
	var sock mangos.Socket
	var err error
	var data []byte
	var o map[string]interface{}

	if sock, err = sub.NewSocket(); err != nil {
		Log.Fatalf("can't get new sub socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		Log.Fatalf("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		Log.Fatalf("cannot subscribe: %s", err.Error())
	}
	for {
		if data, err = sock.Recv(); err != nil {
			Log.Printf("Cannot recv: %s", err.Error())
		}
		if o, err = to_json(data); err != nil {
			Log.Printf("Cannot parse json: %s", err.Error())
		}
		go handle_message(o, cmdline)
	}
}
