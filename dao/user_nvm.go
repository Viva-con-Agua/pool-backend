package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func NVMConfirm(ctx context.Context, token *vcapool.AccessToken) (result *models.NVM, err error) {
	result = new(models.NVM)
	if err = models.NVMConfirmedPermission(token); err != nil {
		return
	}
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.UpdateSet(models.NVMConfirm()),
		result,
	); err != nil {
		return
	}

	return
}

func NVMConfirmUser(ctx context.Context, i *models.NVMIDParam, token *vcapool.AccessToken) (result *models.NVM, err error) {
	result = new(models.NVM)
	if err = models.NVMRejectPermission(token); err != nil {
		return nil, vcago.NewPermissionDenied("nvm", nil)
	}
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.ID}},
		vmdb.UpdateSet(models.NVMConfirm()),
		result,
	); err != nil {
		return
	}

	return
}

func NVMRejectUser(ctx context.Context, i *models.NVMIDParam, token *vcapool.AccessToken) (result *models.NVM, err error) {
	result = new(models.NVM)
	if err = models.NVMRejectPermission(token); err != nil {
		return nil, vcago.NewPermissionDenied("nvm", nil)
	}

	if err = NVMCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: i.ID}}, result); err != nil {
		return
	}
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.ID}},
		vmdb.UpdateSet(models.NVMReject()),
		result,
	); err != nil {
		return
	}
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: i.ID}},
	); err != nil {
		return
	}

	return
}

func NVMWithdraw(ctx context.Context, token *vcapool.AccessToken) (result *models.NVM, err error) {
	result = new(models.NVM)
	if err = NVMCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.UpdateSet(models.NVMWithdraw()),
		result,
	); err != nil {
		return
	}
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
	); err != nil {
		return
	}
	return
}
