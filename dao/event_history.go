package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcapool"
)

func EventStateHistoryInsert(ctx context.Context, i *models.EventStateHistoryCreate, token *vcapool.AccessToken) (err error) {
	if err = EventStateHistoryCollection.InsertOne(ctx, i); err != nil {
		return
	}
	return
}

func EventStateHistoryGet(ctx context.Context, i *models.EventStateHistoryQuery, token *vcapool.AccessToken) (result *[]models.EventStateHistory, list_size int64, err error) {
	result = new([]models.EventStateHistory)
	if err = EventStateHistoryCollection.Aggregate(
		ctx,
		models.EventStatePipeline().Match(i.Filter()).Pipe,
		result,
	); err != nil {
		return
	}
	list_size = int64(len(*result))
	return
}
