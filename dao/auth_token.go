package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

//AuthToken represents the database model for an vpool.AuthToken
type AuthToken struct {
	ID        string            `json:"id" bson:"_id"`
	Token     vcapool.AuthToken `json:"token" bson:"token"`
	UserID    string            `json:"user_id" bson:"user_id"`
	ExpiresAt int64             `json:"expires_at" bson:"expires_at"`
}

//AuthTokenCollection represents the database collection auth_token
var AuthTokenCollection = Database.Collection("auth_token")

//Create store an AuthToken into the database.
func (i *AuthToken) Create(ctx context.Context) (err error) {
	err = AuthTokenCollection.InsertOne(ctx, &i)
	return
}

//Get selects an AuthToken from database based on the bson.M filter.
func (i *AuthToken) Get(ctx context.Context, filter bson.M) (err error) {
	err = AuthTokenCollection.FindOne(ctx, filter, &i)
	return
}

//Delete deletes the AuthToken from database.
func (i *AuthToken) Delete(ctx context.Context) (err error) {
	err = AuthTokenCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
	return
}
