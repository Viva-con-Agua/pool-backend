package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UserCrewDelete(ctx context.Context, id string) (err error) {
	if err = UserCrewCollection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: id}}); err != nil {
		return
	}
	if err = ActiveCollection.TryDeleteOne(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	//reject nvm state
	if err = NVMCollection.TryDeleteOne(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	return
}
