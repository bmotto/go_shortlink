package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/bmotto/go_shortlink/datastore"
)

//RedirectionHandler redirection service
func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Redirection handler")
	var slink datastore.StLink
	var count int

	// get link value
	vars := mux.Vars(r)
	token := vars["token"]
	fmt.Println("token ", token)

	// initialize data store client
	client := datastore.NewClient()
	// read data store
	listSlink, keys := datastore.ReadDataStore(client)

	// look for the original link knowing the shortlink
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
		slink.Count = slink.Count + 1

		// log counter
		logrus.WithFields(logrus.Fields{
			"CreationTimestamp": slink.CreationTimestamp,
			"Origin":            slink.Origin,
			"Token":             slink.Token,
			"Count":             slink.Count,
		}).Info("increment counter for token ", slink.Token)

		// update redirection counter for this token
		datastore.UpdateDataStore(client, keys[count], slink)
	}

	// redirection to the original url
	redirect(slink.Origin, w, r)
}

//redirect redirection to the url
func redirect(url string, w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(url, "www.") {
		url = "www." + url
	}
	if !strings.Contains(url, "http://") {
		url = "http://" + url
	}
	http.Redirect(w, r, url, 301)
}
