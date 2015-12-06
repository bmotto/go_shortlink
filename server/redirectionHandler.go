package server

import(
	"net/http"
	"fmt"
	//"strings"
	"encoding/json"

	"github.com/gorilla/mux"
	//"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/datastore"
)

func RedirectionHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Redirection handler")
	var slink datastore.StLink
	var count int

	// get link value
	vars := mux.Vars(r)
	link := vars["link"]
	token := vars["token"]
	fmt.Println("token ", token)
	fmt.Println("link ", link)


	// initialize data store client
	client := datastore.NewClient()
	// read data store
	listSlink, keys := datastore.ReadDataStore(client)

  // look for the original link knowing the shortlink
	for _, val := range listSlink{
		/*domainName =  domainName + "/"
		token := strings.SplitAfter(link, domainName)
		if len(token) != 2{
			logrus.Warn("init OK")
			break
		}*/
		//if val.Token == token[1]{
		if val.Token == token{
			slink = val
			break
		}
		count = count + 1
	}

	if count >= len(listSlink){
		logrus.Warn("Couldn't find token :", token)
	}else{
		slink.Count = slink.Count + 1

		// log counter
		logrus.WithFields(logrus.Fields{
			"Creation_timestamp": slink.Creation_timestamp,
			"Origin": slink.Origin,
			"Token": slink.Token,
			"Count": slink.Count,
	  }).Info("increment counter for token ", slink.Token)

		datastore.UpdateDataStore(client, keys[count], slink)
	}

	jsonResp, err := json.Marshal(slink)
  if err == nil {
		// write as json the new slink written
		w.Header().Set("JSON",string(jsonResp))
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprintln(w, "value:", slink.Origin)
}
