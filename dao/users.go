package dao

import (
	"context"
	"pool-user/models"
)

func UserInsert(ctx context.Context, i *models.UserDatabase) (r *models.User, err error) {
	//create mailbox
	mailbox := models.NewMailboxDatabase("user")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	// refer the mailbox.ID
	i.MailboxID = mailbox.ID
	// insert user
	if err = UserCollection.InsertOne(ctx, i); err != nil {
		return
	}
	// initial r.
	r = new(models.User)
	//select user from database
	if err = UserCollection.AggregateOne(
		ctx,
		models.UserPipeline().Match(models.UserMatch(i.ID)).Pipe,
		r,
	); err != nil {
		return
	}
	return
}
