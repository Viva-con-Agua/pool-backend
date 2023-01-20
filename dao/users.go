package dao

import (
	"context"
	"pool-backend/models"

	"go.mongodb.org/mongo-driver/bson"
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

func UserDelete(ctx context.Context, id string) (err error) {
	user := new(models.User)
	if err = UserCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}, user); err != nil {
		return
	}
	delete := bson.D{{Key: "user_id", Value: id}}
	if err = AddressesCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = ProfilesCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = UserCrewCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = ActiveCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = NVMCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = NVMCollection.TryDeleteMany(ctx, delete); err != nil {
		return
	}
	if err = AvatarCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = MailboxCollection.TryDeleteOne(ctx, bson.D{{Key: "_id", Value: user.MailboxID}}); err != nil {
		return
	}
	if err = MessageCollection.TryDeleteMany(ctx, bson.D{{Key: "mailbox_id", Value: user.MailboxID}}); err != nil {
		return
	}
	if err = UserCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return
	}
	return
}
