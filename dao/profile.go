package dao

import (
	"context"
	"pool-backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ProfileImport(ctx context.Context, profile *models.ProfileImport) (result *models.Profile, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: profile.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return
	}
	result = profile.Profile(user.ID)
	if err = ProfilesCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}
