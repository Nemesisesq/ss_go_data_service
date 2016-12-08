package gracenote

import (
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nemesisesq/ss_data_service/common"
	dbase "github.com/nemesisesq/ss_data_service/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
)



//func GetMongoSession() *mgo.Session {
//	common.Check(err)
//	return sesh.Copy()
//
//}

func RefreshListings() {

	var _sesh, err = mgo.Dial(os.Getenv("MONGODB_URI"))

	log.SetFormatter(&log.JSONFormatter{})
	//_sesh := GetMongoSession()
	dbname := dbase.GetDatabase()
	db := _sesh.DB(dbname)
	col := db.C("lineups")

	lineups := []Lineup{}

	err = col.Find(nil).All(&lineups)
	log.Println(len(lineups))
	common.Check(err)

	for _, val := range lineups {
		if IsRightLinup(val) {
			go func(lineup Lineup) {
				log.WithField("lineup", lineup.LineupId).Info("refreshing lineup")

				the_json := lineup.GetFreshTVListingsGrid()
				redis_url := os.Getenv("REDISCLOUD_URL")

				u, err := url.Parse(redis_url)

				common.Check(err)

				pass, b := u.User.Password()

				if !b {
					pass = ""
				}
				redisClient := redis.NewClient(&redis.Options{
					Addr:     u.Host,
					Password: pass,
					DB:       0,
				})

				defer redisClient.Close()
				timeout := time.Hour * 5
				redisClient.Set(lineup.LineupId, the_json, timeout)

			}(val)
		}
	}

	//TODO check the geo coordinates for
}
