package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func ActiveWithdraw(ctx context.Context, token *models.AccessToken) (result *models.Active, err error) {
	return activeWithdraw(ctx, token.ID)
}

func activeWithdraw(ctx context.Context, id string) (result *models.Active, err error) {
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		vmdb.UpdateSet(models.ActiveWithdraw()),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	//withdrawn nvm
	if _, err = nvmWithdraw(ctx, id); err != nil {
		return
	}
	//Delete Pool Roles
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: id}},
	); err != nil {
		return
	}
	return
}

func ActiveReject(ctx context.Context, i *models.ActiveParam, token *models.AccessToken) (result *models.Active, err error) {
	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	user := new(models.User)
	filter := i.PermittedFilter(token)
	if err = UserCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(models.ActiveReject()),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	//reject nvm state
	if _, err = nvmReject(ctx, token.ID); err != nil {
		return
	}
	//Delete Pool Roles
	if err = PoolRoleCollection.TryDeleteMany(
		ctx,
		bson.D{{Key: "user_id", Value: i.UserID}},
	); err != nil {
		return
	}
	return
}

func ActiveConfirm(ctx context.Context, i *models.ActiveParam, token *models.AccessToken) (result *models.Active, err error) {
	//check permissions for update an other users active model.
	if err = models.ActivePermission(token); err != nil {
		return
	}
	user := new(models.User)
	filter := i.PermittedFilter(token)
	if err = UserCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(models.ActiveConfirm()),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	return
}

func ActiveRequest(ctx context.Context, token *models.AccessToken) (result *models.Active, err error) {
	//check permissions for active request
	if err = models.ActiveRequestPermission(token); err != nil {
		return
	}
	user := new(models.User)
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: token.ID}},
		vmdb.UpdateSet(models.ActiveRequest()),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	return
}

func activeNew(ctx context.Context, id string, crewID string) (result *models.Active, err error) {
	user := new(models.User)
	update := bson.D{{Key: "nvm", Value: models.NewActive(id, crewID)}}
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		vmdb.UpdateSet(update),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	return
}

func activeClean(ctx context.Context, id string) (result *models.Active, err error) {
	user := new(models.User)
	update := bson.D{{Key: "nvm", Value: models.NVMClean()}}
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		vmdb.UpdateSet(update),
		&user,
	); err != nil {
		return
	}
	result = &user.Active
	return
}

func ActiveNotification(ctx context.Context, i *models.Active, template string) (err error) {
	user := new(models.User)
	if user, err = UsersGetByID(ctx, &models.UserParam{ID: i.UserID}); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.CrewID}}, &crew); err != nil {
		return
	}
	mail := vcago.NewMailData(user.Email, "pool-backend", template, "pool", user.Country)
	mail.AddUser(user.User())
	mail.AddContent(i.ToContent(crew))
	vcago.Nats.Publish("system.mail.job", mail)
	//notification := vcago.NewMNotificationData(result.User.Email, "pool-backend", template, result.User.Country, token.ID)
	//notification.AddUser(result.User.User())
	//notification.AddContent(result.ToContent())
	//vcago.Nats.Publish("system.notification.job", notification)
	return
}
