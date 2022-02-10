package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/google/uuid"
)

type Crew struct {
	ID       string         `json:"id,omitempty" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	City     string         `json:"city" bson:"city"`
	Country  string         `json:"country" bson:"country"`
	Modified vcago.Modified `json:"modified" bson:"modified"`
}

var CrewCollection = Database.Collection("crews").CreateIndex("name", true)

func (i *Crew) Create(ctx context.Context) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.Modified = vcago.NewModified()
	err = UserCollection.InsertOne(ctx, &i)
	return
}
