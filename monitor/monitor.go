package monitor

import (
	/* "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "os"

        "github.com/sirupsen/logrus"
	"nanomsg.org/go/mangos/v2"
        "nanomsg.org/go/mangos/v2/protocol/surveyor"
        _ "nanomsg.org/go/mangos/v2/transport/all"
        */
)

type Status struct {
	Mode	string		`json:"mode"`
	Peers	[]string	`json:"peers"`
}

func NewStatus() *Status {
	status := &Status{}
	return status
}
