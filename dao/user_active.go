package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func ActiveWithdraw(ctx context.Context, token *vcapool.AccessToken) (result *models.Active, err error) {
	result = new(models.Active)
	if err = ActiveCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.UpdateSet(models.NVMWithdraw()),
		result,
	); err != nil {
		return
	}
	//withdrawn nvm
	if err = NVMCollection.TryUpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.UpdateSet(models.NVMWithdraw()),
	); err != nil {
		return
	}
	//Delete Pool Roles
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
	); err != nil {
		return
	}
	return
}

func ActiveReject(ctx context.Context, i *models.ActiveParam, token *vcapool.AccessToken) (result *models.Active, err error) {
	result = new(models.Active)

	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	//update active model.
	if err = ActiveCollection.UpdateOne(
		ctx,
		i.Filter(token),
		vmdb.UpdateSet(models.ActiveReject()),
		result,
	); err != nil {
		return
	}
	//reject nvm state
	if err = NVMCollection.TryUpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
		vmdb.UpdateSet(models.NVMReject()),
	); err != nil {
		return
	}
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
	); err != nil {
		return
	}
	return
}

func ActiveConfirm(ctx context.Context, i *models.ActiveParam, token *vcapool.AccessToken) (result *models.Active, err error) {
	result = new(models.Active)
	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	//update active model.
	if err = ActiveCollection.UpdateOne(
		ctx,
		i.Filter(token),
		vmdb.UpdateSet(models.ActiveConfirm()),
		result,
	); err != nil {
		return
	}
	return
}

func ActiveRequest(ctx context.Context, token *vcapool.AccessToken) (result *models.Active, err error) {
	result = new(models.Active)
	//check permissions for active request
	if err = models.ActiveRequestPermission(token); err != nil {
		return
	}
	if err = ActiveCollection.UpdateOne(
		ctx,
		bson.D{{Key: "user_id", Value: token.ID}},
		vmdb.UpdateSet(models.ActiveRequest()),
		result,
	); err != nil {
		return
	}
	return
}
