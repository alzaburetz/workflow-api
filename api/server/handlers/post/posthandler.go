package post

import (
	"errors"
	"github.com/alzaburetz/workflow-api/api/server/handlers/comment"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
)

type Post struct {
	Id        string `json:"id" bson:"_id_"`
	GroupID   string `json:"group_id" bson:"group_id"`
	Author    User   `json:"author" bson:"author"`
	Name      string `json:"name" bson:"name"`
	Body      string `json:"body" bson:"body"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Comments []comment.Comment `json:"comments" bson:"-"`
	CommentsCount int `json:"comments_count" bson:"comments"`
	LikesCount int `json:"likes" bson:"likescount"`
	LikedByUser bool `json:"liked" bson:"-"`
	Likes []string `json:"-" bson:"likes"`
	Tags []string `json:"tags" bson:"tags"`
}

func (p Post) HasRequiredFields() error {
	if len(p.Name) == 0 {
		return errors.New("Post name is empty")
	} else if len(p.Body) == 0 {
		return errors.New("Post body is empty")
	}

	return nil
}
