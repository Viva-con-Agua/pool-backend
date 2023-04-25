package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
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
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
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

	if i.Status == "confirmed" {
		for _, unit := range i.DepositUnit {
			event := new(models.EventUpdate)
			if err = EventCollection.FindOne(
				ctx,
				bson.D{{Key: "taking_id", Value: unit.TakingID}},
				event,
			); event != nil {
				event.EventState.State = "closed"
				e := new(models.Event)
				if err = EventCollection.UpdateOneAggregate(
					ctx,
					event.Filter(),
					vmdb.UpdateSet(event),
					e,
					models.EventPipeline(token).Match(event.Match()).Pipe,
				); err != nil {
					return
				}

				// Add takings to CRM
				var taking *models.Taking
				if taking, err = TakingGetByID(ctx, &models.TakingParam{ID: unit.TakingID}, token); err != nil {
					log.Print(err)
				}

				taking.EditorID = token.ID
				if err = IDjango.Post(taking, "/v1/pool/taking/create/"); err != nil {
					log.Print(err)
				}

				// Update CRM event
				if err = IDjango.Post(e, "/v1/pool/event/update/"); err != nil {
					log.Print(err)
				}

				// Add participations to event
				participations := new([]models.Participation)

				query := new(models.ParticipationQuery)
				query.EventID = []string{e.ID}
				if err = ParticipationCollection.Aggregate(
					ctx,
					models.ParticipationPipeline().Match(query.Match()).Pipe,
					participations,
				); err != nil {
					return
				}

				if err = IDjango.Post(participations, "/v1/pool/participations/create/"); err != nil {
					log.Print(err)
				}

			}
		}

	}

	return
}

func DepositGet(ctx context.Context, i *models.DepositQuery, token *vcapool.AccessToken) (result *[]models.Deposit, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
	result = new([]models.Deposit)
	if err = DepositCollection.Aggregate(
		ctx,
		depositPipeline().Match(i.Filter(token)).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func DepositGetByID(ctx context.Context, param *models.DepositParam, token *vcapool.AccessToken) (result *models.Deposit, err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return nil, vcago.NewPermissionDenied("takings")
	}
	result = new(models.Deposit)
	var filter bson.D
	if !token.Roles.Validate("admin;employee") {
		filter = bson.D{{Key: "_id", Value: param.ID}, {Key: "crew_id", Value: token.CrewID}}
	} else {
		filter = bson.D{{Key: "_id", Value: param.ID}}
	}
	if err = DepositCollection.AggregateOne(
		ctx,
		depositPipeline().Match(filter).Pipe,
		result,
	); err != nil {
		return
	}
	return
}
