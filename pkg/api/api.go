package api

import (
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sunet/tq/pkg/utils"
)

var Log = logrus.New()
var Handler *mux.Router

func Listen(hostPort string) {
	go func() {
		http.Handle("/", Handler)
		cfg := utils.GetTLSConfig()
		if cfg == nil {
			Log.Fatal(http.ListenAndServe(hostPort, nil))
		} else {
			Log.Debug("Enabling TLS...")
			listener, err := tls.Listen("tcp", hostPort, cfg)
			if err != nil {
				Log.Fatal(err)
			}
			Log.Fatal(http.Serve(listener, nil))
		}
	}()
}

func init() {
	Handler = mux.NewRouter()
	Handler.StrictSlash(true)

	Handler.HandleFunc("/api/status", status).Methods("GET")
	Handler.HandleFunc("/api/channels/all", allChannels).Methods("GET")
	Handler.HandleFunc("/api/channels/list", listChannels).Methods("GET")
	Handler.HandleFunc("/api/channel/{id}", showChannel).Methods("GET")
}
