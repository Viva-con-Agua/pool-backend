package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func CrewInsert(ctx context.Context, i *models.CrewCreate) (r *models.Crew, err error) {
	//create mailbox
	mailbox := models.NewMailboxDatabase("crew")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	r = i.Crew()
	// refer the mailbox.ID
	r.MailboxID = mailbox.ID
	// insert user
	if err = CrewsCollection.InsertOne(ctx, r); err != nil {
		return
	}
	//select user from database
	return
}

func CrewUpdate(ctx context.Context, i *models.CrewUpdate, token *vcapool.AccessToken) (result *models.Crew, err error) {
	if err = models.CrewUpdatePermission(token); err != nil {
		return
	}
	if !token.Roles.Validate("employee;admin") {
		if err = CrewsCollection.UpdateOne(ctx, i.PermittedFilter(token), vmdb.UpdateSet(i.ToCrewUpdateASP()), token); err != nil {
			return
		}
	} else {
		if err = CrewsCollection.UpdateOne(ctx, i.PermittedFilter(token), vmdb.UpdateSet(i), token); err != nil {
			return
		}
	}
	return
}

func CrewDelete(ctx context.Context, i *models.CrewParam) (err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, crew); err != nil {
		return
	}
	if err = MailboxCollection.TryDeleteOne(ctx, bson.D{{Key: "_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = MessageCollection.TryDeleteMany(ctx, bson.D{{Key: "mailbox_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = CrewsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: i.ID}}); err != nil {
		return
	}
	return
}

func CrewGetByID(ctx context.Context, i *models.CrewParam, token *vcapool.AccessToken) (result *models.Crew, err error) {
	result = new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, i.PermittedFilter(token), result); err != nil {
		return
	}
	return
}
