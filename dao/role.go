package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
)

func RoleInsert(ctx context.Context, i *models.RoleRequest, token *vcapool.AccessToken) (result *vmod.Role, err error) {
	filter := i.MatchUser()
	user := new(models.User)
	if err = UserCollection.AggregateOne(
		ctx,
		models.UserPipeline(false).Match(filter).Pipe,
		user,
	); err != nil {
		return
	}
	if result, err = i.NewRole(); err != nil {
		return
	}
	if err = models.RolesPermission(result, user, token); err != nil {
		return
	}
	if err = PoolRoleCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func RoleDelete(ctx context.Context, i *models.RoleRequest, token *vcapool.AccessToken) (result *vmod.Role, err error) {
	filter := i.MatchUser()
	user := new(models.User)
	if err = UserCollection.FindOne(
		ctx,
		filter,
		user,
	); err != nil {
		return
	}
	if err = PoolRoleCollection.FindOne(
		ctx,
		i.Filter(),
		&result,
	); err != nil {
		return
	}
	if err = models.RolesDeletePermission(result, token); err != nil {
		return
	}
	if err = PoolRoleCollection.DeleteOne(ctx, i.Filter()); err != nil {
		return
	}

	return
}
