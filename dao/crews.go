package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Crew struct {
		ID       string         `json:"id,omitempty" bson:"_id"`
		Name     string         `json:"name" bson:"name"`
		Email    string         `json:"email" bson:"email"`
		Cities   []City         `json:"cities" bson:"cities"`
		Modified vcago.Modified `json:"modified" bson:"modified"`
	}
	City struct {
		City        string   `json:"city" bson:"city"`
		Country     string   `json:"country" bson:"country"`
		CountryCode string   `json:"country_code" bson:"country_code"`
		PlaceID     string   `json:"place_id" bson:"place_id"`
		Position    Position `json:"position" bson:"position"`
	}
	Position struct {
		Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lin" bson:"lin"`
	}
)

var CrewsCollection = Database.Collection("crews").CreateIndex("name", true)

func (i *Crew) Create(ctx context.Context) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.Modified = vcago.NewModified()
	err = CrewsCollection.InsertOne(ctx, &i)
	return
}

func (i *Crew) Get(ctx context.Context, id string) (err error) {
	err = CrewsCollection.FindOne(ctx, bson.M{"_id": id}, &id)
	return
}

//Update updates the crew model and all user_crew model email and name.
func (i *Crew) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	update := bson.M{"$set": &i}
	if err = CrewsCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update); err != nil {
		return
	}
	update = bson.M{"$set": bson.M{"email": i.Email, "name": i.Name, "modified.updated": i.Modified.Updated}}
	if err = UserCrewCollection.Update(ctx, bson.M{"crew_id": i.ID}, update); err != nil && !vcago.MongoNoUpdated(err) {
		return
	}
	err = nil
	return
}

func (i *Crew) Delete(ctx context.Context) (err error) {
	err = CrewsCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
	return
}

type CrewList []Crew

type CrewQuery vcapool.CrewQuery

func (i *CrewQuery) Match() *vcago.MongoMatch {
	match := vcago.NewMongoMatch()
	match.LikeString("_id", i.ID)
	match.EqualString("email", i.Email)
	match.ElemMatchList("cities", "city", i.Cities)
	return match
}
func (i *CrewQuery) List(ctx context.Context) (r *CrewList, err error) {
	pipe := vcago.NewMongoPipe()
	pipe.Match(i.Match())
	r = new(CrewList)
	err = CrewsCollection.Aggregate(ctx, pipe.Pipe, r)
	return
}
