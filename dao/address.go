package dao

import (
	"context"
	"pool-backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

func AddressImport(ctx context.Context, address *models.AddressImport) (result *models.Address, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: address.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return
	}
	result = address.Address(user.ID)
	if err = AddressesCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}
