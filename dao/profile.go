package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func ProfileInsert(ctx context.Context, i *models.ProfileCreate, token *vcapool.AccessToken) (result *models.Profile, err error) {
	result = i.Profile(token.ID)
	if err = ProfilesCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func ProfileUpdate(ctx context.Context, i *models.ProfileUpdate, token *vcapool.AccessToken) (result *models.Profile, err error) {
	filter := i.PermittedFilter(token)
	if err = ProfilesCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		result,
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
