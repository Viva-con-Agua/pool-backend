package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func TakingsPipeline() *vmdb.Pipeline {
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
	return pipe
}
func TakingInsert(ctx context.Context, i *models.TakingCreate, token *vcapool.AccessToken) (r *models.Taking, err error) {
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
	if err = TakingCollection.AggregateOne(
		ctx,
		TakingsPipeline().Match(bson.D{{Key: "_id", Value: taking.ID}}).Pipe,
		r,
	); err != nil {
		return
	}
	return
}

func TakingUpdate(ctx context.Context, i *models.TakingUpdate, token *vcapool.AccessToken) (r *models.Taking, err error) {
	taking := new(models.Taking)

	//permission
	filter := bson.D{{Key: "_id", Value: "cant find"}}
	if token.Roles.Validate("employee;admin") {
		filter = bson.D{{Key: "_id", Value: i.ID}}
	}
	if token.PoolRoles.Validate("finance") {
		filter = bson.D{{Key: "_id", Value: i.ID}, {Key: "crew_id", Value: token.CrewID}}
	}
	pipeline := TakingsPipeline().Match(filter).Pipe
	if err = TakingCollection.AggregateOne(ctx, pipeline, taking); err != nil {
		return
	}
	//update
	if !(taking.State.Open.Amount == 0 && taking.State.Wait.Amount == 0) || true {
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
				if err = SourceCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: v.ID}}); err != nil {
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
	}
	r = new(models.Taking)
	update := vmdb.UpdateSet(i)
	if err = TakingCollection.UpdateOneAggregate(ctx, filter, update, r, pipeline); err != nil {
		return
	}
	return
}

func TakingGet(ctx context.Context, query *models.TakingQuery, token *vcapool.AccessToken) (result *[]models.Taking, err error) {
	filter := bson.D{{Key: "_id", Value: "cant find"}}
	if token.Roles.Validate("employee") {
		filter = query.Filter()
	} else if token.PoolRoles.Validate("finance") {
		query.CrewID = []string{token.CrewID}
		filter = query.Filter()
	}
	pipeline := TakingsPipeline().Match(filter).Pipe
	result = new([]models.Taking)
	if err = TakingCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	return
}

func TakingGetByID(ctx context.Context, param *vmod.IDParam, token *vcapool.AccessToken) (result *models.Taking, err error) {
	filter := bson.D{{Key: "_id", Value: "cant find"}}
	if token.Roles.Validate("employee;admin") {
		filter = bson.D{{Key: "_id", Value: param.ID}}
	}
	if token.PoolRoles.Validate("finance") {
		filter = bson.D{{Key: "_id", Value: param.ID}, {Key: "crew_id", Value: token.CrewID}}
	}
	pipeline := TakingsPipeline().Match(filter).Pipe
	result = new(models.Taking)
	if err = TakingCollection.AggregateOne(ctx, pipeline, result); err != nil {
		return
	}
	return
}

func TakingDeletetByID(ctx context.Context, param *vmod.IDParam, token *vcapool.AccessToken) (result *vmod.DeletedResponse, err error) {
	filter := bson.D{{Key: "_id", Value: "cant find"}}
	if token.Roles.Validate("employee;admin") {
		filter = bson.D{{Key: "_id", Value: param.ID}}
	}
	if token.PoolRoles.Validate("finance") {
		filter = bson.D{{Key: "_id", Value: param.ID}, {Key: "crew_id", Value: token.CrewID}}
	}

	err = TakingCollection.DeleteOne(ctx, filter)
	return
}
