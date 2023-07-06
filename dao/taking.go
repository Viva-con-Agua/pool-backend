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
	takingDatabase := new(models.Taking)
	filter := i.PermittedFilter(token)
	if err = TakingCollection.AggregateOne(ctx, models.TakingPipeline().Match(filter).Pipe, takingDatabase); err != nil {
		return
	}
	if err = takingDatabase.UpdatePermission(token); err != nil {
		return
	}
	//i.State = takingDatabase.State
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
			if models.SourceDeletePermission(takingDatabase, token) {
				if err = SourceCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: v.ID}}); err != nil {
					return
				}
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
	return
}

func TakingGet(ctx context.Context, query *models.TakingQuery, token *vcapool.AccessToken) (result *[]models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	result = new([]models.Taking)
	filter := query.PermittedFilter(token)
	pipeline := models.TakingPipeline().Match(filter).Pipe
	if err = TakingCollection.Aggregate(
		ctx,
		pipeline,
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

func TakingGetByIDSystem(ctx context.Context, id string) (result *models.Taking, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	if err = TakingCollection.AggregateOne(
		ctx,
		models.TakingPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func TakingDeletetByIDSystem(ctx context.Context, id string) (err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = TakingCollection.DeleteOne(ctx, filter)
	return
}

func TakingDeletetByID(ctx context.Context, param *models.TakingParam, token *vcapool.AccessToken) (err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	err = TakingCollection.DeleteOne(ctx, param.PermittedFilter(token))
	return
}

func TakingUpdateEventStatus(ctx context.Context, state string, id string) (err error) {
	switch state {
	case "canceled":
		filter := bson.D{{Key: "_id", Value: id}}
		update := bson.D{{Key: "state.no_income", Value: true}}
		if err = TakingCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(update), nil); err != nil {
			return
		}
	}
	return
}
