package group

import (. "app/server/handlers/user"
		. "app/server/handlers/post"
		"errors")

type Group struct {
	Creator User `json:"creator" bson:"creator"`
	Id string `json:"id" bson:"_id_"`
	Name string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	UserCount int `json:"usercount" bson:"usercount"`
	Posts []Post `json:"posts" bson:"posts"`
}

func (gr *Group)HasRequiredFields() error {
	if len(gr.Name) == 0 {
		return errors.New("Name is empty")
	} else if len(gr.Description) == 0 {
		return errors.New("Description is empty")
	}
	return nil
}
