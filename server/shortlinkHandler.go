package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/bmotto/go_shortlink/datastore"
	"github.com/bmotto/go_shortlink/util"
)

//shortlink short link creation struct
type shortlink struct {
	Link   string `json:"link"`
	Custom string `json:"custom"`
}

//ShortlinkHandler short link hanlder
func ShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Shortlink handler")
	var slink datastore.StLink
	var shortL shortlink
	var link string
	var parameter string

	// read config.yml
	domainName := viper.GetString("domainName")

	// get body content
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
	}
	if len(body) != 0 {
		err = json.Unmarshal(body, &shortL)
		if err != nil {
			logrus.Error(err)
		}
		link = shortL.Link
		parameter = shortL.Custom
	} else {
		// get link value
		vars := mux.Vars(r)
		link = strings.Split(vars["link"], "&custom=")[0]
		parameter = strings.SplitAfter(vars["link"], "&custom=")[1]
	}

	// initialize data store client
	client := datastore.NewClient()
	// read data store
	listSlink, keys := datastore.ReadDataStore(client)

	// generate shortLink
	gen := genToken(link, parameter, listSlink)

	slink.Origin = link

	url := "http://" + slink.Origin
	resp, err := http.Get(url)
	if err != nil || resp == nil {
		logrus.Error("couldn't get " + slink.Origin + " err :" + err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Error("status KO for", slink.Origin)
		return
	}

	slink.Token = gen

	jsonResp := datastore.WriteDataStore(client, keys, slink)
	listSlink, keys = datastore.ReadDataStore(client)

	// write as json the new slink written in the body
	fmt.Fprintln(w, "shortlink:", domainName+slink.Token)
	fmt.Fprintln(w, "json:", string(jsonResp))
}

//genToken check ten times if the parameter is present on
// the data store and generate an other one
func genToken(link string, parameter string, listSlink []datastore.StLink) (gen string) {
	// read configuration yaml
	shortLinkSize := viper.GetInt("shortlinkSize")

	for i := 0; i < 10; i++ {
		found := false
		if parameter != "" {
			// check parameter already used
			for _, sl := range listSlink {
				if sl.Token == parameter {
					found = true
					break
				}
			}

			if found {
				if len(parameter) < shortLinkSize {
					gen = parameter + util.RandNum(shortLinkSize-len(parameter))
				} else {
					gen = parameter[:shortLinkSize-2] + util.RandNum(2)
				}
			} else {
				gen = parameter
				break
			}
		} else {
			gen = util.RandString(shortLinkSize)
		}

		// check if the generate short link is already
		// in the data store
		for _, sl := range listSlink {
			if sl.Token == gen {
				found = true
			}
		}
		if !found {
			break
		}
	}
	return
}
