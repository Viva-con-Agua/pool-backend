package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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

func UserCrewImport(ctx context.Context, imp *models.UserCrewImport) (result *models.UserCrew, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: imp.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return nil, vcago.NewBadRequest("user", err.Error(), userFilter)
	}
	crew := new(models.Crew)
	crewFilter := bson.D{{Key: "_id", Value: imp.CrewID}}
	if err = CrewsCollection.FindOne(ctx, crewFilter, crew); err != nil {
		return nil, vcago.NewBadRequest("crew", err.Error(), crewFilter)
	}
	result = models.NewUserCrew(user.ID, crew.ID, crew.Name, crew.Email, crew.MailboxID)
	if err = UserCrewCollection.InsertOne(ctx, result); err != nil {
		return
	}
	if err = ActiveCollection.InsertOne(ctx, models.NewActive(user.ID, crew.ID)); err != nil {
		return
	}
	if err = NVMCollection.InsertOne(ctx, models.NewNVM(user.ID)); err != nil {
		return
	}
	return

}
