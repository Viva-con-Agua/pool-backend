package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcapool"
)

func MailboxGetByID(ctx context.Context, i *models.MailboxParam, token *vcapool.AccessToken) (r *models.Mailbox, err error) {

	r = new(models.Mailbox)

	if err = MailboxCollection.AggregateOne(
		ctx,
		models.MailboxPipeline().Match(i.PermittedFilter(token)).Pipe,
		r,
	); err != nil {
		return
	}
	return
}
