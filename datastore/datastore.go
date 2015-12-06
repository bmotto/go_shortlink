package datastore

import(
  "time"
  "encoding/json"

  "github.com/go-redis/redis"
  "github.com/spf13/viper"
  "github.com/Sirupsen/logrus"

  "github.com/bmotto/go_shortlink/util"
)

type StLink struct{
	Creation_timestamp int64
	Origin string
	Token string
	Count int
}

// create a new redis client
func NewClient() *redis.Client{
  addr := viper.GetString("datastore.addr")
  password := viper.GetString("datastore.password")
  db := viper.GetInt("datastore.DB")
  db64 := int64(db)

  client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db64,
  })

  pong, err := client.Ping().Result()
  if err != nil{
    logrus.Error("DataStore new client pong", pong, " error ", err)
    panic(err)
  }
  return client
}

// read the redis datastore
func ReadDataStore(client *redis.Client) ([]StLink, []string){
  var listStLink []StLink
  var deserialized StLink

  keys, err := client.Keys("*").Result()
  if err != nil{
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

// write data in Data Store (redis)
func WriteDataStore(client *redis.Client, keys []string, slink StLink) string{
  // generate a random id
  // and check it doesn't already exists in the data store
  id := ""
  for {
    found := false
    id = util.RandString(20)
    for _,key := range keys{
      if key == id{
        found = true
      }
    }
    if !found{
      break
    }
  }

  time_now := time.Now()
  duration := threeMonthFromDate(time_now)
  slink.Creation_timestamp = time_now.Unix()

  serialized, err := json.Marshal(slink)
  if err == nil {
      status := client.Set(id, string(serialized), duration).Err()
      if status != nil{
        logrus.Error("Write DataStore: set client values error ", err)
        panic(status)
      }
  }

  return string(serialized)
}

// Update data store values
func UpdateDataStore(client *redis.Client, key string, newSLink StLink){
  creation_time := time.Unix(newSLink.Creation_timestamp, 0)
  duration := threeMonthFromDate(creation_time)

  serialized, err := json.Marshal(newSLink)
  if err == nil {
      status := client.Set(key, string(serialized), duration).Err()
      if status != nil{
        logrus.Error("Update DataStore: set client values error ", err)
        panic(status)
      }
  }
}

// compute three month duration from creation time
func threeMonthFromDate(creation_time time.Time) time.Duration{
  var time_3months time.Time
  if int(creation_time.Month()) + 3 > 12 {
    month := int(creation_time.Month()) + 3 -12
    time_3months = time.Date(creation_time.Year() + 1,
                             time.Month(month), creation_time.Day(),
                             creation_time.Hour(),
                             creation_time.Minute(),
                             creation_time.Second(),
                             creation_time.Nanosecond(),
                             creation_time.Location())
  }else{
    time_3months = time.Date(creation_time.Year(),
                             creation_time.Month() + 3,
                             creation_time.Day(),
                             creation_time.Hour(),
                             creation_time.Minute(),
                             creation_time.Second(),
                             creation_time.Nanosecond(),
                             creation_time.Location())
  }
  return time_3months.Sub(creation_time)
}
