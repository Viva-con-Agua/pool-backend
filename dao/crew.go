package dao

import (
	"context"
	"pool-backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

func CrewInsert(ctx context.Context, i *models.Crew) (r *models.Crew, err error) {
	//create mailbox
	mailbox := models.NewMailboxDatabase("crew")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	// refer the mailbox.ID
	i.MailboxID = mailbox.ID
	// insert user
	if err = CrewsCollection.InsertOne(ctx, i); err != nil {
		return
	}
	// initial r.
	r = i
	//select user from database
	return
}

func CrewDelete(ctx context.Context, id string) (err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}, crew); err != nil {
		return
	}
	if err = MailboxCollection.TryDeleteOne(ctx, bson.D{{Key: "_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = MessageCollection.TryDeleteMany(ctx, bson.D{{Key: "mailbox_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = CrewsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return
	}
	return
}
