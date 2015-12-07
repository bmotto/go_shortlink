package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/bmotto/go_shortlink/server"
)

func main() {
	// load configuration yaml
	err := LoadConfig()
	if err != nil {
		fmt.Println("error loading config.yml: %v", err)
		panic(err)
	}

	f := initLogFile()
	//  close the file when quit the main
	defer f.Close()

	// read configuration
	port := viper.GetString("port")
	ip := viper.GetString("ip")

	r := mux.NewRouter()
	r.HandleFunc("/shortlink/{link}", server.ShortlinkHandler).Methods("POST")
	r.HandleFunc("/admin/{token}", server.AdminHandler).Methods("GET")
	r.HandleFunc("/{token}", server.RedirectionHandler).Methods("GET")
	http.ListenAndServe(ip+":"+port, r)
}

//initLogFile initialize log file
func initLogFile() *os.File {
	// open file and/or create file
	f, err := os.OpenFile("File.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// Log with the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(f)
	// Only log the info severity or above.
	logrus.SetLevel(logrus.InfoLevel)

	return f
}
