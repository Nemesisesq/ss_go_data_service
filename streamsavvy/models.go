package streamsavvy


import (
	"github.com/nemesisesq/ss_data_service/gracenote"
	"gopkg.in/mgo.v2/bson"
)

//import "golang.org/x/tools/go/ssa/interp"

type Content struct {
	Id                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title               string                   `json:"title" bson:"title"`
	GuideboxData        map[string]interface{}   `json:"guidebox_data" bson:"guidebox_data"`
	OnNetflix           bool                     `json:"on_netflix" bson:"on_netflix"`
	Channel             []map[string]interface{} `json:"channel" bson:"channel"`
	CurrPopScore        float32                  `json:"curr_pop_score" bson:"curr_pop_score"`
	ChannelsLastChecked string                   `json:"channels_last_checked" bson:"channels_last_checked"`
	Modified            string                   `json:"modified" bson:"modified"`
}


type Favorites struct {
	//TODO Technical Debt write setter functions to ensure the UUID of the user is the same as User.UUID property in favorites
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User        User      `json:"user" bson:"user"`
	UserUUID    int16     `json:"user_uuid" bson:"user_uuid"`
	ContentList []Content `json:"content_list" bson:"content_list"`
	SportList   []gracenote.Sport `json:"sport_list" bson:"sport_list"`
	TeamList    []gracenote.Team `json:"team_list" bson:"team_list"`
}

type User struct {
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	UserId   string `json:"user_id" bson:"user_id"`
}
