package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func AddressInsert(ctx context.Context, i *models.AddressCreate, token *vcapool.AccessToken) (result *models.Address, err error) {
	result = i.Address(token.ID)
	if err = AddressesCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func AddressGet(ctx context.Context, i *models.AddressQuery, token *vcapool.AccessToken) (result *[]models.Address, err error) {
	filter := i.PermittedFilter(token)
	result = new([]models.Address)
	if err = AddressesCollection.Find(ctx, filter, result); err != nil {
		return
	}
	return
}

func AddressGetByID(ctx context.Context, i *models.AddressParam, token *vcapool.AccessToken) (result *models.Address, err error) {
	filter := i.PermittedFilter(token)
	if err = AddressesCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func AddressUpdate(ctx context.Context, i *models.AddressUpdate, token *vcapool.AccessToken) (result *models.Address, err error) {
	filter := i.PermittedFilter(token)
	if err = AddressesCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	return
}

func AddressDelete(ctx context.Context, i *models.AddressParam, token *vcapool.AccessToken) (result *models.NVM, err error) {
	filter := i.PermittedFilter(token)
	if err = AddressesCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	if result, err = NVMWithdraw(ctx, token); err != nil {
		return
	}
	return
}

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
