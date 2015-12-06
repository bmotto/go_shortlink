package server

import(
	"net/http"
	"fmt"
	//"strings"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"

	"github.com/bmotto/go_shortlink/datastore"
)

func AdminHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Admin handler")
	var slink datastore.StLink
	var count int

	fmt.Fprintln(w, "Bonjour Admin")

	// read configuration yaml
	domainName := viper.GetString("domainName")

	// get link value
	vars := mux.Vars(r)
	//link := vars["link"]
	token := vars["token"]
	link := vars["link"]
	fmt.Println("link ", link)

	// initialize data store client
	client := datastore.NewClient()

	// read data store
	listSlink, _ := datastore.ReadDataStore(client)

  // look for the original link knowing the shortlink
	domainName =  domainName + "/"
	for _, val := range listSlink{
		/*token := strings.SplitAfter(link, domainName)
		if len(token) != 2{
			fmt.Println("***************WARNING***********")
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
		jsonResp, err := json.Marshal(slink)
		if err == nil {
			// write as json the new slink written
			w.Header().Set("JSON",string(jsonResp))
			w.WriteHeader(http.StatusOK)
		}

		fmt.Fprintln(w, "value:", domainName + slink.Token)
	}
}
