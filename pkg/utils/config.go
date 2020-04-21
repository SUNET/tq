package utils

import (
	"crypto/tls"
	"os"
)

var cfg *tls.Config

func GetTLSConfig() *tls.Config {
	return cfg
}

func init() {
	crt := os.Getenv("TQ_TLS_CERT")
	key := os.Getenv("TQ_TLS_KEY")
	if len(crt) > 0 && len(key) > 0 {
		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			panic(err)
		}
		cfg = &tls.Config{Certificates: []tls.Certificate{cert}}
	}
}
