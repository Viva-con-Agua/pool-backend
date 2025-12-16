package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func UsersUserCrewInsert(ctx context.Context, i *models.UsersCrewCreate, token *models.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UsersCrewCreatePermission(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	result = models.NewUserCrew(i.UserID, crew)
	filter := bson.D{{Key: "_id", Value: i.UserID}}
	update := bson.D{{Key: "crew", Value: result}}
	if err = UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	if _, err = activeNew(ctx, i.UserID, crew.ID); err != nil {
		return
	}
	if _, err = nvmNew(ctx, i.UserID); err != nil {
		return
	}
	return
}

func UserCrewInsert(ctx context.Context, i *models.UserCrewCreate, token *models.AccessToken) (result *models.UserCrew, err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	result = models.NewUserCrew(token.ID, crew)
	filter := bson.D{{Key: "_id", Value: token.ID}}
	update := bson.D{{Key: "crew", Value: result}}
	if err = UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	if _, err = activeNew(ctx, token.ID, crew.ID); err != nil {
		return
	}
	if _, err = nvmNew(ctx, token.ID); err != nil {
		return
	}
	return
}

func UserCrewUpdate(ctx context.Context, i *models.UserCrewUpdate, token *models.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UserCrewUpdatePermission(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	i.OrganisationID = crew.OrganisationID
	user := new(models.User)
	if err = UserCollection.UpdateOne(ctx, i.PermittedFilter(token), vmdb.UpdateSet(i), &user); err != nil {
		return
	}
	result = &user.Crew
	//reset active and nvm
	if _, err = activeWithdraw(ctx, i.UserID); err != nil {
		return
	}
	if _, err = nvmWithdraw(ctx, i.UserID); err != nil {
		return
	}
	return
}

func UsersCrewUpdate(ctx context.Context, i *models.UserCrewUpdate, token *models.AccessToken) (result *models.UserCrew, err error) {
	if err = i.UsersCrewUpdatePermission(token); err != nil {
		return
	}

	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.CrewFilter(), crew); err != nil {
		return
	}
	i.OrganisationID = crew.OrganisationID
	user := new(models.User)
	if err = UserCollection.UpdateOne(ctx, i.Match(), vmdb.UpdateSet(i), &user); err != nil {
		return
	}
	result = &user.Crew
	if _, err = activeWithdraw(ctx, i.UserID); err != nil {
		return
	}
	//reset active and nvm
	if _, err = nvmWithdraw(ctx, i.UserID); err != nil {
		return
	}

	return
}

func UserCrewDelete(ctx context.Context, id string) (err error) {
	update := bson.D{{Key: "crew", Value: models.UserCrewClean()}}
	if err = UserCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	if _, err = activeClean(ctx, id); err != nil {
		return
	}
	//delete
	if _, err = nvmClean(ctx, id); err != nil {
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
	/*
		result = models.NewUserCrew(user.ID, crew)
		if err = UserCrewCollection.InsertOne(ctx, result); err != nil {
			return
		}*/
	/*
		active := imp.ToActive(user.ID)
		if err = ActiveCollection.InsertOne(ctx, active); err != nil {
			return
		}*/
	roles := imp.ToRoles(user.ID)
	for _, role := range roles {
		if err = PoolRoleCollection.InsertOne(ctx, &role); err != nil {
			return
		}
	}
	return

}

func UserCrewSync(i models.UserCrew) (result *models.UserCrew, err error) {
	go func() {
		if err = IDjango.Post(i, "/v1/pool/profile/crew/"); err != nil {
			log.Print(err)
		}
	}()
	return
}
