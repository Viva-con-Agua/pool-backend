package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func newTakingsPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.Lookup("deposit_unit_taking", "_id", "taking_id", "deposit_units")
	pipe.LookupMatch("deposit_unit_taking", "_id", "taking_id", "wait", bson.D{{Key: "deposit.status", Value: bson.D{{Key: "$in", Value: bson.A{"wait", "open"}}}}})
	pipe.LookupMatch("deposit_unit_taking", "_id", "taking_id", "confirmed", bson.D{{Key: "deposit.status", Value: "confirmed"}})
	pipe.Lookup("sources", "_id", "taking_id", "sources")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	pipe.LookupUnwind("events", "_id", "taking_id", "event")
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.wait.amount", Value: bson.D{{Key: "$sum", Value: "$wait.money.amount"}}},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.confirmed.amount", Value: bson.D{{Key: "$sum", Value: "$confirmed.money.amount"}}},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "money.amount", Value: bson.D{{Key: "$sum", Value: "$sources.money.amount"}}}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "state.open.amount", Value: bson.D{
		{Key: "$subtract", Value: bson.A{"$money.amount", bson.D{{Key: "$add", Value: bson.A{"$state.wait.amount", "$state.confirmed.amount"}}}}},
	}}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.wait.currency", Value: "$currency"},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.confirmed.currency", Value: "$currency"},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "money.currency", Value: "$currency"}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "state.open.currency", Value: "$currency"}}}})
	pipe.Lookup("activity_user", "_id", "model_id", "activities")
	return pipe
}
func TakingInsert(ctx context.Context, i *models.TakingCreate, token *vcapool.AccessToken) (r *models.Taking, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
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
	r = new(models.Taking)
	activity := models.NewActivityDB(token.ID, "taking", taking.ID, "Successfully created", "created")
	if err = ActivityCollection.InsertOne(ctx, activity); err != nil {
		return
	}
	if err = TakingCollection.AggregateOne(
		ctx,
		newTakingsPipeline().Match(bson.D{{Key: "_id", Value: taking.ID}}).Pipe,
		r,
	); err != nil {
		return
	}

	return
}

func TakingUpdate(ctx context.Context, i *models.TakingUpdate, token *vcapool.AccessToken) (r *models.Taking, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
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
	r = new(models.Taking)
	if err = TakingCollection.UpdateOneAggregate(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(i),
		r,
		newTakingsPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
	); err != nil {
		return
	}
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
	}

	return
}

func TakingGet(ctx context.Context, query *models.TakingQuery, token *vcapool.AccessToken) (result *[]models.Taking, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
	result = new([]models.Taking)
	if err = TakingCollection.Aggregate(
		ctx,
		newTakingsPipeline().Match(query.Filter(token)).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func TakingGetByID(ctx context.Context, param *models.TakingParam, token *vcapool.AccessToken) (result *models.Taking, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
	result = new(models.Taking)
	if err = TakingCollection.AggregateOne(
		ctx,
		newTakingsPipeline().Match(param.Filter(token)).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func TakingDeletetByID(ctx context.Context, param *models.TakingParam, token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return vcago.NewPermissionDenied("takings")
	}
	err = TakingCollection.DeleteOne(ctx, param.Filter(token))
	return
}
