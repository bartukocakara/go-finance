package v1

import (
	"encoding/json"
	"net/http"

	"finance-app/internal/config"

	"github.com/sirupsen/logrus"
)

// API for returning version
// When server starts, we set version and than use it if necessary.

type ServerVersion struct {
	Version string `json:"version"`
}

var versionJSON []byte

func init(){
	var err error
	versionJSON, err = json.Marshal(ServerVersion{
		Version: config.Version,
	})
	if err != nil {
		panic(err)
	}
}

// Serves version information
func VersionHandler(w http.ResponseWriter, _ *http.Request){
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(versionJSON); err != nil {
		logrus.WithError(err).Debug("error writing version")
	}
}