package post

import (. "app/server/handlers/user"
		"errors")

type Post struct {
	Id int `json:"id" bson:"_id_:`
	GroupID int `json:"group_id" bson:"group_id"`
	Author User `json:"author" bson:"author"`
	Name string `json:"name" bson:"name"`
	Body string `json:"body" bson:"body"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

func (p Post)HasRequiredFields() error  {
	if len(p.Name) == 0 {
		return errors.New("Post name is empty")
	} else if len(p.Body) == 0 {
		return errors.New("Post body is empty")
	}

	return nil
}