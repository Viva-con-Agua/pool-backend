package dao

import "github.com/Viva-con-Agua/vcago"

type Roles struct {
	ID       string         `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Label    string         `json:"label" bson:"label"`
	Root     string         `json:"root" bson:"root"`
	Modified vcago.Modified `json:"modified" bson:"modified"`
}
