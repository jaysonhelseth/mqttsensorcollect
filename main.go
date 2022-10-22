package main

import (
	"MqttSensorCollect/io"
	"MqttSensorCollect/models"
	"embed"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/fs"
	"log"
	"net/http"
	"time"
)

var (
	//go:embed static
	resources embed.FS
)

func sensorData(writer http.ResponseWriter, request *http.Request, airTemp *models.Temp) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, airTemp.Read())
}

func main() {
	airHumDevicePath := flag.String("ah", "", "The air temperature and humidity device path.")
	devMode := flag.Bool("d", false, "Enable dev mode.")
	flag.Parse()

	if *airHumDevicePath == "" {
		log.Fatalln("Missing flag")
	}

	var airtemp = models.Temp{}
	go io.ReadFromSerial(*airHumDevicePath, &airtemp)

	var files http.Handler
	if *devMode {
		files = http.FileServer(http.Dir("./static"))
	} else {
		// The fs.Sub line removes the static folder to match how devMode makes paths.
		serverRoot, _ := fs.Sub(resources, "static")
		files = http.FileServer(http.FS(serverRoot))
	}

	router := mux.NewRouter()
	router.Handle("/", files)
	router.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		sensorData(w, r, &airtemp)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
