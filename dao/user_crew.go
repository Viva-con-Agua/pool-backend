package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func UsersUserCrewInsert(ctx context.Context, i *models.UsersCrewCreate, token *vcapool.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UsersCrewCreatePermission(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	result = models.NewUserCrew(i.UserID, crew)
	if err = UserCrewCollection.InsertOne(ctx, result); err != nil {
		return
	}
	if err = ActiveCollection.InsertOne(ctx, models.NewActive(i.UserID, crew.ID)); err != nil {
		return
	}
	if err = NVMCollection.InsertOne(ctx, models.NewNVM(i.UserID)); err != nil {
		return
	}
	return
}

func UserCrewInsert(ctx context.Context, i *models.UserCrewCreate, token *vcapool.AccessToken) (result *models.UserCrew, err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	result = models.NewUserCrew(token.ID, crew)
	if err = UserCrewCollection.InsertOne(ctx, result); err != nil {
		return
	}
	if err = ActiveCollection.InsertOne(ctx, models.NewActive(token.ID, crew.ID)); err != nil {
		return
	}
	if err = NVMCollection.InsertOne(ctx, models.NewNVM(token.ID)); err != nil {
		return
	}
	return
}

func UserCrewUpdate(ctx context.Context, i *models.UserCrewUpdate, token *vcapool.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UserCrewUpdatePermission(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	i.OrganisationID = crew.OrganisationID
	if err = UserCrewCollection.UpdateOne(ctx, i.PermittedFilter(token), vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	//reset active and nvm
	if err = ActiveCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
		vmdb.UpdateSet(models.ActiveWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	//reject nvm state
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
		vmdb.UpdateSet(models.NVMWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}

	return
}

func UsersCrewUpdate(ctx context.Context, i *models.UserCrewUpdate, token *vcapool.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UsersCrewUpdatePermission(token); err != nil {
		return
	}

	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	i.OrganisationID = crew.OrganisationID
	if err = UserCrewCollection.UpdateOne(ctx, i.Match(), vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	//reset active and nvm
	if err = ActiveCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
		vmdb.UpdateSet(models.ActiveWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}
	//reject nvm state
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
		vmdb.UpdateSet(models.NVMWithdraw()),
		nil,
	); err != nil && vmdb.ErrNoDocuments(err) {
		return
	}

	return
}

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
	if err = NewsletterCollection.TryDeleteOne(
		ctx,
		bson.D{{Key: "user_id", Value: id}, {Key: "value", Value: "regional"}},
	); err != nil {
		return
	}
	return
}

func UserCrewImport(ctx context.Context, imp *models.UserCrewImport) (result *models.UserCrew, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: imp.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return nil, vcago.NewBadRequest(models.UserCollection, err.Error(), userFilter)
	}
	crew := new(models.Crew)
	crewFilter := bson.D{{Key: "_id", Value: imp.CrewID}}
	if err = CrewsCollection.FindOne(ctx, crewFilter, crew); err != nil {
		return nil, vcago.NewBadRequest(models.CrewCollection, err.Error(), crewFilter)
	}
	if crew.Status != "active" {
		return nil, vcago.NewBadRequest(models.CrewCollection, "crew_is_dissolved", nil)
	}
	result = models.NewUserCrew(user.ID, crew)
	if err = UserCrewCollection.InsertOne(ctx, result); err != nil {
		return
	}
	active := imp.ToActive(user.ID)
	if err = ActiveCollection.InsertOne(ctx, active); err != nil {
		return
	}
	nvm := imp.ToNVM(user.ID)
	if err = NVMCollection.InsertOne(ctx, nvm); err != nil {
		return
	}
	roles := imp.ToRoles(user.ID)
	for _, role := range roles {
		if err = PoolRoleCollection.InsertOne(ctx, &role); err != nil {
			return
		}
	}
	return

}
