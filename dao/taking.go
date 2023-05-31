package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	TakingCreatedActivity = &models.ActivityDatabase{ModelType: "taking", Comment: "Successfully created", Status: "created"}
	TakingUpdatedActivity = &models.ActivityDatabase{ModelType: "taking", Comment: "Successfully updated", Status: "updated"}
)

func TakingInsert(ctx context.Context, i *models.TakingCreate, token *vcapool.AccessToken) (result *models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	//create taking model form i.
	taking := i.TakingDatabase()
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	if i.NewSource != nil {
		sources := i.SourceList(taking.ID)
		if err = SourceCollection.InsertMany(ctx, sources.InsertMany()); err != nil {
			return
		}
	}
	if err = ActivityCollection.InsertOne(ctx, TakingCreatedActivity.New(token.ID, taking.ID)); err != nil {
		return
	}
	if err = TakingCollection.AggregateOne(
		ctx,
		models.TakingPipeline().Match(bson.D{{Key: "_id", Value: taking.ID}}).Pipe,
		&result,
	); err != nil {
		return
	}

	return
}

func TakingUpdate(ctx context.Context, i *models.TakingUpdate, token *vcapool.AccessToken) (result *models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	takingDatabase := new(models.TakingDatabase)
	if err = TakingCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, takingDatabase); err != nil {
		return
	}
	i.State = &takingDatabase.State
	for _, v := range i.Sources {
		//create new sources
		if v.ID == "" {
			v.TakingID = i.ID
			newSource := v.Source()
			if err = SourceCollection.InsertOne(ctx, newSource); err != nil {
				return
			}
		}
		if v.UpdateState == "deleted" {
			deleteSource := new(models.Source)
			if err = SourceCollection.FindOne(ctx, bson.D{{Key: "_id", Value: v.ID}}, deleteSource); err != nil {
				return
			}
		}
		if v.UpdateState == "updated" {
			databaseSource := new(models.Source)
			if err = SourceCollection.FindOne(
				ctx,
				bson.D{{Key: "_id", Value: v.ID}},
				databaseSource,
			); err != nil {
				return
			}
			if err = SourceCollection.UpdateOne(
				ctx,
				bson.D{{Key: "_id", Value: v.ID}},
				vmdb.UpdateSet(v),
				nil,
			); err != nil {
				return
			}
		}
	}
	if err = ActivityCollection.InsertOne(ctx, TakingUpdatedActivity.New(token.ID, takingDatabase.ID)); err != nil {
		return
	}
	if err = TakingCollection.UpdateOneAggregate(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(i),
		&result,
		models.TakingPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
	); err != nil {
		return
	}
	/*
		event := new(models.EventUpdate)
		if err = EventCollection.FindOne(
			ctx,
			bson.D{{Key: "taking_id", Value: i.ID}},
			event,
		); event != nil {
			event.EventState.State = "finished"
			result := new(models.Event)
			if err = EventCollection.UpdateOneAggregate(
				ctx,
				event.Filter(),
				vmdb.UpdateSet(event),
				result,
				models.EventPipeline(token).Match(event.Match()).Pipe,
			); err != nil {
				return
			}
		}*/

	return
}

func TakingGet(ctx context.Context, query *models.TakingQuery, token *vcapool.AccessToken) (result *[]models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	result = new([]models.Taking)
	if err = TakingCollection.Aggregate(
		ctx,
		models.TakingPipeline().Match(query.PermittedFilter(token)).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func TakingGetByID(ctx context.Context, param *models.TakingParam, token *vcapool.AccessToken) (result *models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	filter := param.PermittedFilter(token)
	if err = TakingCollection.AggregateOne(
		ctx,
		models.TakingPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func TakingDeletetByID(ctx context.Context, param *models.TakingParam, token *vcapool.AccessToken) (err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	err = TakingCollection.DeleteOne(ctx, param.PermittedFilter(token))
	return
}
