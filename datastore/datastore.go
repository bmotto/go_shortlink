package datastore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/go-redis/redis"
	"github.com/bmotto/go_shortlink/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/bmotto/go_shortlink/util"
)

//StLink Shortlink struct
type StLink struct {
	CreationTimestamp int64
	Origin            string
	Token             string
	Count             int
}

//NewClient create a new redis client
func NewClient() *redis.Client {
	addr := viper.GetString("datastore.addr")
	password := viper.GetString("datastore.password")
	db := viper.GetInt("datastore.DB")
	db64 := int64(db)

	fmt.Println("add ", addr)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db64,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		logrus.Error("DataStore new client pong", pong, " error ", err)
		panic(err)
	}
	return client
}

//ReadDataStore read the redis datastore
func ReadDataStore(client *redis.Client) ([]StLink, []string) {
	var listStLink []StLink
	var deserialized StLink

	keys, err := client.Keys("*").Result()
	if err != nil {
		logrus.Error("Read DataStore: get client keys error ", err)
		panic(err)
	}

	for _, val := range keys {
		serialized, err := client.Get(val).Bytes()
		if err != nil {
			logrus.Error("Read DataStore : serialize values error ", err)
			panic(err)
		}

		err = json.Unmarshal(serialized, &deserialized)
		if err != nil {
			logrus.Error("Read DataStore : json unmarshal error ", err)
			panic(err)
		}

		listStLink = append(listStLink, deserialized)
	}

	return listStLink, keys
}

//WriteDataStore write data in Data Store (redis)
func WriteDataStore(client *redis.Client, keys []string, slink StLink) string {
	// generate a random id
	// and check it doesn't already exists in the data store
	id := ""
	for {
		found := false
		id = util.RandString(20)
		for _, key := range keys {
			if key == id {
				found = true
			}
		}
		if !found {
			break
		}
	}

	timeNow := time.Now()
	duration := threeMonthFromDate(timeNow)
	slink.CreationTimestamp = timeNow.Unix()

	serialized, err := json.Marshal(slink)
	if err == nil {
		status := client.Set(id, string(serialized), duration).Err()
		if status != nil {
			logrus.Error("Write DataStore: set client values error ", err)
			panic(status)
		}
	}

	return string(serialized)
}

//UpdateDataStore Update data store values
func UpdateDataStore(client *redis.Client, key string, newSLink StLink) {
	creationTime := time.Unix(newSLink.CreationTimestamp, 0)
	duration := threeMonthFromDate(creationTime)

	serialized, err := json.Marshal(newSLink)
	if err == nil {
		status := client.Set(key, string(serialized), duration).Err()
		if status != nil {
			logrus.Error("Update DataStore: set client values error ", err)
			panic(status)
		}
	}
}

//threeMonthFromDate compute three month duration from creation time
func threeMonthFromDate(creationTime time.Time) time.Duration {
	var time3months time.Time
	if int(creationTime.Month())+3 > 12 {
		month := int(creationTime.Month()) + 3 - 12
		time3months = time.Date(creationTime.Year()+1,
			time.Month(month), creationTime.Day(),
			creationTime.Hour(),
			creationTime.Minute(),
			creationTime.Second(),
			creationTime.Nanosecond(),
			creationTime.Location())
	} else {
		time3months = time.Date(creationTime.Year(),
			creationTime.Month()+3,
			creationTime.Day(),
			creationTime.Hour(),
			creationTime.Minute(),
			creationTime.Second(),
			creationTime.Nanosecond(),
			creationTime.Location())
	}
	return time3months.Sub(creationTime)
}
