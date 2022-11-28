package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func MailboxGetByID(ctx context.Context, id string, token *vcapool.AccessToken) (r *models.Mailbox, err error) {
	r = new(models.Mailbox)
	if err = MailboxCollection.AggregateOne(
		ctx,
		models.MailboxPipeline().Match(bson.D{{Key: "_id", Value: id}}).Pipe,
		r,
	); err != nil {
		return
	}
	if r.Type == "crew" {
		crew := new(models.Crew)
		if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
			return
		}
		if crew.MailboxID != id {
			return nil, vcago.NewPermissionDenied("mailbox", id)
		}
		return
	} else {
		if token.MailboxID != id {
			return nil, vcago.NewPermissionDenied("mailbox", id)
		}
	}
	return
}
