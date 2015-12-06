package main

import(
	"net/http"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"

  "github.com/bmotto/go_shortlink/server"
)

func main(){
	// load configuration yaml
	err := LoadConfig()
	if err != nil{
		panic(err);
	}

	f := initLogFile()
	//  close the file when quit the main
  defer f.Close()

	logrus.Info("init OK")

	// read configuration
	port := ":" + viper.GetString("port")

	r := mux.NewRouter()
	r.HandleFunc("/shortlink/{link}", server.ShortlinkHandler)
	r.HandleFunc("/admin/{link}/{token}", server.AdminHandler)
	r.HandleFunc("/{link}/{token}", server.RedirectionHandler)
	http.ListenAndServe(port, r)
}

func initLogFile() *os.File {
	// open file
	f, err := os.OpenFile("File.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
  if err != nil {
      fmt.Printf("error opening file: %v", err)
  }

  // Log with the default ASCII formatter.
  logrus.SetFormatter(&logrus.TextFormatter{})

  // Output to stderr instead of stdout, could also be a file.
  logrus.SetOutput(f)

  // Only log the warning severity or above.
  logrus.SetLevel(logrus.InfoLevel)

	return f
}
