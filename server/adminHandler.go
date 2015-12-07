package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/bmotto/go_shortlink/datastore"
)

//AdminHandler monitoring service
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Admin handler")
	var slink datastore.StLink
	var count int

	// read configuration yaml
	domainName := viper.GetString("domainName")
	// get token value
	vars := mux.Vars(r)
	token := vars["token"]

	// initialize data store client
	client := datastore.NewClient()
	// read data store
	listSlink, _ := datastore.ReadDataStore(client)

	// look for the original link knowing the shortlink
	domainName = domainName + "/"
	for _, val := range listSlink {
		if val.Token == token {
			slink = val
			break
		}
		count = count + 1
	}

	if count >= len(listSlink) {
		logrus.Warn("Couldn't find token :", token)
	} else {
		jsonResp, err := json.Marshal(slink)
		if err != nil {
			logrus.Error("Couldn't marshal slink, error: " + err.Error())
			return
		}

		// write as json the new slink written in the body
		fmt.Fprintln(w, "json:", string(jsonResp))
	}
}
