package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func ProfileInsert(ctx context.Context, i *models.ProfileCreate, token *models.AccessToken) (result *models.Profile, err error) {
	result = i.Profile(token.ID)
	update := bson.D{{Key: "profile", Value: result}}
	filter := bson.D{{Key: "_id", Value: result.UserID}}
	if err = UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	return
}

func ProfileGetByID(ctx context.Context, i *models.UserParam, token *models.AccessToken) (result *models.User, err error) {
	if err = models.UsersDetailsPermission(token); err != nil {
		return
	}
	if err = UserCollection.AggregateOne(ctx, models.UserPermittedPipeline(token).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func ProfileUpdate(ctx context.Context, i *models.ProfileUpdate, token *models.AccessToken) (result *models.Profile, err error) {
	filter := i.PermittedFilter(token)
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		&user,
	); err != nil {
		return
	}
	result = &user.Profile
	if i.Birthdate == 0 {
		var nvm *models.NVM
		if nvm, err = NVMWithdraw(ctx, token); err == nil {
			NvmSync(nvm)
		}
	}
	return
}

func ProfileSync(i *models.Profile) (result *models.Profile, err error) {
	go func() {
		user := new(models.User)
		userFilter := bson.D{{Key: "_id", Value: i.UserID}}
		if err = UserCollection.FindOne(context.Background(), userFilter, user); err != nil {
			log.Print(err)
			return
		}
		if err = IDjango.Post(user, "/v1/pool/profile/"); err != nil {
			log.Print(err)
		}
	}()
	return
}

func NvmSync(i *models.NVM) (result *models.NVM, err error) {
	go func() {
		if err = IDjango.Post(result, "/v1/pool/profile/nvm/"); err != nil {
			log.Print(err)
		}
	}()
	return
}

func UsersProfileUpdate(ctx context.Context, i *models.ProfileUpdate, token *models.AccessToken) (result *models.Profile, err error) {
	if err = models.UsersEditPermission(token); err != nil {
		return
	}
	user := new(models.User)
	if err = UserCollection.UpdateOne(ctx, i.Match(), vmdb.UpdateSet(i.ToUserUpdate()), &user); err != nil {
		return
	}
	result = &user.Profile
	return
}

func ProfileImport(ctx context.Context, profile *models.ProfileImport) (result *models.Profile, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: profile.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return
	}
	result = profile.Profile(user.ID)
	update := bson.D{{Key: "profile", Value: result}}
	filter := bson.D{{Key: "_id", Value: result.UserID}}
	if err = UserCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	return
}
