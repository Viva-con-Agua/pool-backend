package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func NVMConfirm(ctx context.Context, token *models.AccessToken) (result *models.NVM, err error) {
	if err = models.NVMConfirmedPermission(token); err != nil {
		return
	}
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: token.ID}},
		vmdb.UpdateSet(models.NVMConfirm()),
		&user,
	); err != nil {
		return
	}
	result = &user.NVM
	return
}

func NVMConfirmUser(ctx context.Context, i *models.NVMIDParam, token *models.AccessToken) (result *models.NVM, err error) {
	if err = models.NVMPermission(token); err != nil {
		return
	}
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(models.NVMConfirm()),
		&user,
	); err != nil {
		return
	}
	result = &user.NVM

	return
}

func NVMReject(ctx context.Context, i *models.NVMIDParam, token *models.AccessToken) (result *models.NVM, err error) {
	if err = models.NVMPermission(token); err != nil {
		return
	}
	return nvmReject(ctx, i.ID)
}

func nvmReject(ctx context.Context, id string) (result *models.NVM, err error) {
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		vmdb.UpdateSet(models.NVMReject()),
		&user,
	); err != nil {
		return
	}
	result = &user.NVM
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	return
}

func NVMWithdraw(ctx context.Context, token *models.AccessToken) (result *models.NVM, err error) {
	if err = models.NVMPermission(token); err != nil {
		return
	}
	return nvmWithdraw(ctx, token.ID)
}

func nvmWithdraw(ctx context.Context, id string) (result *models.NVM, err error) {
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		vmdb.UpdateSet(models.NVMWithdraw()),
		&user,
	); err != nil {
		return
	}
	result = &user.NVM
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	return
}
