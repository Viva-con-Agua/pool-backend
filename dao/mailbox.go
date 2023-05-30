package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcapool"
)

func MailboxGetByID(ctx context.Context, i *models.MailboxParam, token *vcapool.AccessToken) (r *models.Mailbox, err error) {
	filter := i.PermittedFilter(token)
	r = new(models.Mailbox)
	if err = MailboxCollection.AggregateOne(ctx, models.MailboxPipeline().Match(filter).Pipe, r); err != nil {
		return
	}
	return
}
