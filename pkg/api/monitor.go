package api

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

	"github.com/SUNET/tq/pkg/message"
	"github.com/SUNET/tq/pkg/meta"
	"github.com/gorilla/mux"
)

func status(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"number-of-gorutines": runtime.NumGoroutine(),
		"version":             meta.Version(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func all_channels(w http.ResponseWriter, r *http.Request) {
	result := make([]*message.MessageChannel, len(message.Channels))
	i := 0
	for _, ch := range message.Channels {
		result[i] = ch
		i++
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func list_channels(w http.ResponseWriter, r *http.Request) {
	result := make([]string, len(message.Channels))
	i := 0
	for id, _ := range message.Channels {
		result[i] = strconv.FormatUint(id, 10)
		i++
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func show_channel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
	} else {
		result := message.Channels[id]
		if result == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	}
}
