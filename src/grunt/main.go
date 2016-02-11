package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	graceful "gopkg.in/tylerb/graceful.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"time"
)

type SMTP struct {
	Username string
	Password string
	Server   string
}

type Config struct {
	Services   []*Service          `json:"services"`
	ServiceMap map[string]*Service `json:omit`
	Mail       SMTP
}

var config Config

func main() {
	var port int
	flag.IntVar(&port, "p", 9901, "specify port to use.  defaults to 9901.")

	config.ServiceMap = make(map[string]*Service)
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Usage: grunt gruntfile.yml")
	}
	gruntfile := flag.Arg(0)
	data, err := ioutil.ReadFile(gruntfile)
	if err != nil {
		log.Fatal("Error reading %v: %v", gruntfile, err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error in YML parsing: %v", err)
	}
	for _, service := range config.Services {
		config.ServiceMap[service.EndPoint] = service
	}

	log.Infof("SMTP: %+v", config.Mail)

	// Expose the endpoints
	r := mux.NewRouter()
	r.HandleFunc("/rest/service", GetServices).Methods("GET")
	r.HandleFunc("/rest/service/{id}", GetService).Methods("GET")
	r.HandleFunc("/rest/service/{id}", StartService).Methods("POST")
	r.HandleFunc("/rest/job/{id}", GetJob).Methods("GET")
	r.HandleFunc("/rest/job/{id}/file/{filename}", GetJobFile).Methods("GET")

	r.HandleFunc("/help.html", Help).Methods("GET")
	r.HandleFunc("/jobs.html", Jobs).Methods("GET")
	r.HandleFunc("/job/{id}", JobDetail).Methods("GET")
	r.HandleFunc("/services.html", Services).Methods("GET")
	r.HandleFunc("/submit/{id}.html", Submit).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

	http.Handle("/", r)
	log.Infof("Starting grunt on http://localhost:%v", port)
	graceful.Run(fmt.Sprintf(":%d", port), 1*time.Second, nil)

}
