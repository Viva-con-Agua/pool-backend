package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type UserCrew vcapool.UserCrew

var UserCrewCollection = Database.Collection("user_crew").CreateIndex("user_id", true)

type UserCrewCreateRequest struct {
	CrewID string `json:"crew_id"`
}

func (i *UserCrewCreateRequest) Create(ctx context.Context, userID string) (r *UserCrew, err error) {
	crew := new(Crew)
	if err = CrewsCollection.FindOne(ctx, bson.M{"_id": i.CrewID}, crew); err != nil {
		return
	}
	userCrew := vcapool.NewUserCrew(userID, i.CrewID, crew.Name)
	r = (*UserCrew)(userCrew)
	err = UserCrewCollection.InsertOne(ctx, r)
	return
}
