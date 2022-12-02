package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func filterID(id string) bson.D {
	return bson.D{{Key: "_id", Value: id}}
}

func updateWaitTaking(amount int64) bson.D {
	return bson.D{{Key: "$inc", Value: bson.D{{Key: "state.open.amount", Value: -amount}, {Key: "state.wait.amount", Value: amount}}}}
}

func depositPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(DepositUnitCollection.Name, "_id", "deposit_id", "deposit_units")
	pipe.LookupUnwind(TakingCollection.Name, "deposit_units.taking_id", "_id", "deposit_units.taking")
	pipe.Append(bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"}, {Key: "deposit_units", Value: bson.D{
				{Key: "$push", Value: "$deposit_units"},
			}},
		}},
	})
	pipe.LookupUnwind(DepositCollection.Name, "_id", "_id", "deposits")
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "deposits.deposit_units", Value: "$deposit_units"}}}})
	pipe.Append(bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$deposits"}}}})
	pipe.LookupUnwind(CrewsCollection.Name, "crew_id", "_id", "crew")
	return pipe
}

func DepositInsert(ctx context.Context, i *models.DepositCreate, token *vcapool.AccessToken) (r *models.Deposit, err error) {
	deposit, depositUnits := i.DepositDatabase(token)
	taking := new(models.TakingDatabase)
	for _, unit := range depositUnits {
		if err = TakingCollection.FindOne(ctx, filterID(unit.TakingID), taking); err != nil {
			return
		}
	}
	deposit.ReasonForPayment, err = GetNewReasonForPayment(ctx, i.CrewID)
	if err != nil {
		log.Print(err)
		err = nil
	}

	for _, unit := range depositUnits {
		if err = DepositUnitCollection.InsertOne(ctx, unit); err != nil {
			return
		}
	}
	if err = DepositCollection.InsertOne(ctx, deposit); err != nil {
		return
	}
	r = new(models.Deposit)
	if err = DepositCollection.AggregateOne(ctx, depositPipeline().Match(bson.D{{Key: "_id", Value: deposit.ID}}).Pipe, r); err != nil {
		return
	}
	return
}

func DepositUpdate(ctx context.Context, i *models.DepositUpdate, token *vcapool.AccessToken) (result *models.Deposit, err error) {
	depositDatabase := new(models.DepositDatabase)
	if err = DepositCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, depositDatabase); err != nil {
		return
	}
	i.Money = depositDatabase.Money
	if err = DepositCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, vmdb.UpdateSet(i), nil); err != nil {
		return
	}
	result = new(models.Deposit)
	if err = DepositCollection.AggregateOne(
		ctx,
		depositPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func DepositGet(ctx context.Context, i *models.DepositQuery, token *vcapool.AccessToken) (result *[]models.Deposit, err error) {
	result = new([]models.Deposit)
	if err = DepositCollection.Aggregate(
		ctx,
		depositPipeline().Match(i.Filter()).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func DepositGetByID(ctx context.Context, param *models.DepositParam) (result *models.Deposit, err error) {
	result = new(models.Deposit)
	if err = DepositCollection.AggregateOne(
		ctx,
		depositPipeline().Match(bson.D{{Key: "_id", Value: param.ID}}).Pipe,
		result,
	); err != nil {
		return
	}
	return
}
