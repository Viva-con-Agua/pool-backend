package dao

import (
	"context"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"go.mongodb.org/mongo-driver/bson"
)

type Role vcago.Role

var PoolRoleCollection = Database.Collection("pool_roles").CreateMultiIndex(bson.D{{Key: "name", Value: 1}, {Key: "user_id", Value: 1}}, true)

func (i *Role) Create(ctx context.Context) (err error) {
	if err = PoolRoleCollection.InsertOne(ctx, i); err != nil {
		return
	}
	update := bson.M{"last_update": time.Now().Format(time.RFC3339), "modified.updated": time.Now().Unix()}
	err = UserCollection.UpdateOne(ctx, bson.M{"_id": i.UserID}, update)
	return
}

func (i *Role) Get(ctx context.Context, filter bson.M) (err error) {
	err = PoolRoleCollection.FindOne(ctx, filter, i)
	return
}

func (i *Role) Delete(ctx context.Context) (err error) {
	if err = PoolRoleCollection.DeleteOne(ctx, bson.M{"_id": i.ID}); err != nil {
		return
	}
	update := bson.M{"last_update": time.Now().Format(time.RFC3339), "modified.updated": time.Now().Unix()}
	err = UserCollection.UpdateOne(ctx, bson.M{"_id": i.UserID}, update)
	return
}

type RoleRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
