package streamsavvy

import ()

//import "golang.org/x/tools/go/ssa/interp"

type Content struct {
	Title               string                   `json:"title"`
	GuideboxData        map[string]interface{}   `json:"guidebox_data"`
	OnNetflix           bool                     `json:"on_netflix"`
	Channel             []map[string]interface{} `json:"channel"`
	CurrPopScore        float32                  `json:"curr_pop_score"`
	ChannelsLastChecked string                   `json:"channels_last_checked"`
	Modified            string                   `json:"modified"`
}

type Favorites struct {
	//TODO Technical Debt write setter functions to ensure the UUID of the user is the same as User.UUID property in favorites
	User        User      `json:"user" bson:"user"`
	UserUUID    int16     `json:"user_uuid" bson:"user_uuid"`
	ContentList []Content `json:"content_list" bson:"content_list"`
}

type User struct {
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	UserId   string `json:"user_id" bson:"user_id"`
}
