package notification


type Notification struct {
	Uid string `json:"-" bson:"_id_"`
	Seen bool `json:"seen" bson:"seen"`
	Object string `json:"comment" bson:"comment"`
	Created int64 `json:"created" bson:"created"`
	UserEmail string `json:"-" bson:"useremail"`
}
