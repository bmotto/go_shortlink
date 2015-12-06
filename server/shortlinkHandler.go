package server

import(
	"net/http"
	"fmt"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"

	"github.com/bmotto/go_shortlink/datastore"
	"github.com/bmotto/go_shortlink/util"
)

func ShortlinkHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Shortlink handler")
	var slink datastore.StLink

	// read configuration yaml
	shortLinkSize := viper.GetInt("shortlinkSize")
	domainName := viper.GetString("domainName")

	// get link value
	vars := mux.Vars(r)
	link := vars["link"]
	fmt.Println("link ", link)

	// initialize data store client
	client := datastore.NewClient()
	// read data store
	listSlink, keys := datastore.ReadDataStore(client)
	found := true

	// generate shortLink
	gen := ""
	// check ten times if the parameter is present on
	// the data store and generate an other one
	for i:=0; i<10; i++ {
		found = false
		if strings.Contains(link, "&custom="){
			parameter := strings.SplitAfter(link, "&custom=")

			// check parameter already used
			for _, sl := range listSlink{
				if sl.Token == parameter[1]{
					found = true
					break
				}
			}

			if found {
				if len(parameter[1]) < shortLinkSize {
					gen = parameter[1] + util.RandNum(shortLinkSize - len(parameter[1]))
				}else {
					gen = parameter[1][:shortLinkSize - 2 ] + util.RandNum(2)
				}
			}else{
				gen = parameter[1]
				break
			}
		}else{
			gen = util.RandString(shortLinkSize)
		}

		// check if the generate short link is already
		// in the data store
		for _, sl := range listSlink{
			if sl.Token == gen{
				found = true;
			}
		}
		if !found{
			break
		}
	}

	slink.Origin = strings.Split(link, "&custom=")[0]

	url := "http://" + slink.Origin
	resp, err := http.Get(url)
	if err != nil || resp == nil{
		logrus.Error("couldn't get ", slink.Origin)
		return
	}
	if resp.StatusCode != http.StatusOK{
		logrus.Error("status KO for", slink.Origin)
		return
	}

	slink.Token = gen

	jsonResp := datastore.WriteDataStore(client, keys, slink)
	listSlink, keys = datastore.ReadDataStore(client)

	// write as json the new slink written
	w.Header().Set("JSON",string(jsonResp))
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "value:", domainName + gen)

}
