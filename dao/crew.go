package dao

import (
	"context"
	"pool-user/models"
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
