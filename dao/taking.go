package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func TakingGet(ctx context.Context, query *models.TakingQuery, token *vcapool.AccessToken) (result *[]models.Taking, listSize int64, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	result = new([]models.Taking)
	filter := query.PermittedFilter(token)
	sort := query.Sort()
	pipeline := models.TakingPipeline().SortFields(sort).Match(filter).Sort(sort).Skip(query.Skip, 0).Limit(query.Limit, 100).Pipe
	if err = TakingCollection.Aggregate(
		ctx,
		pipeline,
		result,
	); err != nil {
		return
	}
	opts := options.Count().SetHint("_id_")
	if query.FullCount != "true" {
		opts.SetSkip(query.Skip).SetLimit(query.Limit)
	}
	if cursor, cErr := UserViewCollection.Collection.CountDocuments(ctx, filter, opts); cErr != nil {
		print(cErr)
		listSize = 0
	} else {
		listSize = cursor
	}
	return
}

func TakingGetByID(ctx context.Context, param *vmod.IDParam, token *vcapool.AccessToken) (result *models.Taking, err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	filter := models.TakingPermittedFilter(param, token)
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

func TakingDeletetByID(ctx context.Context, param *vmod.IDParam, token *vcapool.AccessToken) (err error) {
	if err = models.TakingPermission(token); err != nil {
		return
	}
	taking := new(models.Taking)
	filter := models.TakingPermittedFilter(param, token)
	if err = TakingCollection.AggregateOne(
		ctx,
		models.TakingPipeline().Match(filter).Pipe,
		&taking,
	); err != nil {
		return
	}
	if taking.Event.ID != "" {
		return vcago.NewBadRequest("taking", "depending_in_event")
	}
	err = TakingCollection.DeleteOne(ctx, filter)
	return
}
