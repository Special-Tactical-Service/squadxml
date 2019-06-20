package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

const (
	squadxmlFile            = "squad.xml"
	defaultHttpWriteTimeout = 20
	defaultHttpReadTimeout  = 20
)

func main() {
	path := os.Getenv("SQUADXML_PATH")
	logrus.WithField("squadxml_path", path).Info("Starting server...")

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(path, squadxmlFile))
	}))

	if err := http.ListenAndServe(os.Getenv("SQUADXML_HOST"), nil); err != nil {
		logrus.WithError(err).Fatal("Error starting server")
	}
}
