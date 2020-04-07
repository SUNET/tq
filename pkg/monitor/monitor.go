package monitor

/* "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "os"

        "github.com/sirupsen/logrus"
	"go.nanomsg.org/mangos/v3"
        "go.nanomsg.org/mangos/v3/protocol/surveyor"
        _ "go.nanomsg.org/mangos/v3/transport/all"
*/

type Status struct {
	Mode  string   `json:"mode"`
	Peers []string `json:"peers"`
}

func NewStatus() *Status {
	status := &Status{}
	return status
}
