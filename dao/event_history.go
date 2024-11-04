package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmod"
)

func EventStateHistoryInsert(ctx context.Context, i *models.EventStateHistoryCreate, token *models.AccessToken) (err error) {
	if err = EventStateHistoryCollection.InsertOne(ctx, i); err != nil {
		return
	}
	return
}

func EventStateHistoryGet(ctx context.Context, i *models.EventStateHistoryQuery, token *models.AccessToken) (result *[]models.EventStateHistory, list_size int64, err error) {
	result = new([]models.EventStateHistory)
	pipeline := models.EventStatePipeline().Match(i.Filter()).Count().Pipe
	if err = EventStateHistoryCollection.Aggregate(
		ctx,
		pipeline,
		result,
	); err != nil {
		return
	}
	count := vmod.Count{}
	var cErr error
	if cErr = EventStateHistoryCollection.AggregateOne(ctx, pipeline, &count); cErr != nil {
		print(cErr)
		list_size = 1
	} else {
		list_size = int64(count.Total)
	}
	return
}
