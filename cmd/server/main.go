package main

import (
	"finance-app/internal/api"
	"finance-app/internal/config"
	"finance-app/internal/database"
	"flag"
	"net"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	host = flag.String("host", "0.0.0.0", "host for listen")
	port = flag.String("port", "8088", "port for listen")
)

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithField("version", config.Version).Debug("Starting server.")
	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database")
	}
	router, err := api.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	var addr = net.JoinHostPort(*host, *port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}

}
