package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcapool"
)

func MailboxGetByID(ctx context.Context, i *models.MailboxParam, token *vcapool.AccessToken) (result *models.Mailbox, err error) {
	filter := i.PermittedFilter(token)
	if err = MailboxCollection.AggregateOne(ctx, models.MailboxPipeline(token).Match(filter).Pipe, &result); err != nil {
		return
	}
	return
}
