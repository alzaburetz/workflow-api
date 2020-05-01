package comment

import . "github.com/alzaburetz/workflow-api/api/server/handlers/user"

type Comment struct {
	Uid     string `json:"id" bson:"_id_"`
	PostId  string `json:"postid" bson:"postid"`
	Author  User   `json:"author" bson:"author"`
	At      User   `json:"at" bson:"at"`
	Body    string `json:"body" bson:"body"`
	Created int64  `json:"created" bson:"created"`
	Edited  int64  `json:"edited" bson:"edited"`
	Likes   int    `json:"likes" bson:"likes"`
}

func (c *Comment) BodyNotEmpty() bool {
	return len(c.Body) > 0
}
