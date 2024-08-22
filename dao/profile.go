package dao

import (
	"context"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func ProfileInsert(ctx context.Context, i *models.ProfileCreate, token *vcapool.AccessToken) (result *models.Profile, err error) {
	result = i.Profile(token.ID)
	if err = ProfileCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func ProfileGetByID(ctx context.Context, i *models.UserParam, token *vcapool.AccessToken) (result *models.User, err error) {
	if err = models.UsersDetailsPermission(token); err != nil {
		return
	}
	if err = UserCollection.AggregateOne(ctx, models.UserPermittedPipeline(token).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func ProfileUpdate(ctx context.Context, i *models.ProfileUpdate, token *vcapool.AccessToken) (result *models.Profile, err error) {
	filter := i.PermittedFilter(token)
	birthdate := time.Unix(i.Birthdate, 0)
	if i.Birthdate != 0 {
		i.BirthdateDatetime = birthdate.Format("2006-01-02")
	} else {
		i.BirthdateDatetime = ""
	}
	if err = ProfileCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		&result,
	); err != nil {
		return
	}
	if i.Birthdate == 0 {
		var nvm *models.NVM
		if nvm, err = NVMWithdraw(ctx, token); err == nil {
			go func() {
				if err = IDjango.Post(nvm, "/v1/pool/profile/nvm/"); err != nil {
					log.Print(err)
				}
			}()
		}
	}
	return
}

func ProfileSync(ctx context.Context, i models.Profile, token *vcapool.AccessToken) (result *models.Profile, err error) {
	go func() {
		if err = IDjango.Post(i, "/v1/pool/profile/"); err != nil {
			log.Print(err)
		}
	}()
	return
}

func UsersProfileUpdate(ctx context.Context, i *models.ProfileUpdate, token *vcapool.AccessToken) (result *models.Profile, err error) {
	if err = models.UsersEditPermission(token); err != nil {
		return
	}
	if err = ProfileCollection.UpdateOne(ctx, i.Match(), vmdb.UpdateSet(i.ToUserUpdate()), &result); err != nil {
		return
	}
	return
}

func ProfileImport(ctx context.Context, profile *models.ProfileImport) (result *models.Profile, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: profile.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return
	}
	result = profile.Profile(user.ID)
	if err = ProfileCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}
