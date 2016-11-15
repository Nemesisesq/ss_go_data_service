package gracenote

import (
	"net/url"
	"os"
	"time"

	"github.com/nemesisesq/ss_data_service/common"
	dbase "github.com/nemesisesq/ss_data_service/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
	"github.com/Sirupsen/logrus"
)

var sesh, err = mgo.Dial(os.Getenv("MONGODB_URI"))

func GetMongoSession() *mgo.Session {
	common.Check(err)
	return sesh.Copy()

}

func (g *Guide) RefreshListings() {

			//logrus.SetFormatter()
	_sesh := GetMongoSession()
	dbname := dbase.GetDatabase()
	db := _sesh.DB(dbname)
	col := db.C("lineups")

	lineups := []Lineup{}

	err := col.Find("").All(lineups)

	common.Check(err)

	for _, val := range lineups[:3] {
		go func(lineup Lineup) {
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
				DB:      0,
			})

			defer redisClient.Close()
			timeout := time.Hour * 5
			redisClient.Set(lineup.LineupId, the_json, timeout)

		}(val)
	}

	//TODO check the geo coordinates for
}
