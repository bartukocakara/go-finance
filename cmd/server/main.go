package main

import (
	"finance-app/internal/database"
	"flag"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	host = flag.String("host", "0.0.0.0", "host for listen")
	port = flag.String("port", "8088", "port for listen")
)

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	_, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database")
	}
	var addr = net.JoinHostPort(*host, *port)
	router := mux.NewRouter()
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}

}
