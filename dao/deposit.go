package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func DepositInsert(ctx context.Context, i *models.DepositCreate, token *vcapool.AccessToken) (result *models.Deposit, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	deposit, depositUnits := i.DepositDatabase(token)
	taking := new(models.TakingDatabase)
	for _, unit := range depositUnits {
		if err = TakingCollection.FindOne(ctx, models.Match(unit.TakingID), taking); err != nil {
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
	if err = DepositCollection.AggregateOne(ctx, models.DepositPipeline().Match(bson.D{{Key: "_id", Value: deposit.ID}}).Pipe, &result); err != nil {
		return
	}
	return
}

func DepositUpdate(ctx context.Context, i *models.DepositUpdate, token *vcapool.AccessToken) (result *models.Deposit, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	depositDatabase := new(models.DepositDatabase)
	if err = DepositCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, depositDatabase); err != nil {
		return
	}
	i.Money = depositDatabase.Money
	if err = DepositCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, vmdb.UpdateSet(i), nil); err != nil {
		return
	}
	if err = DepositCollection.AggregateOne(
		ctx,
		models.DepositPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
		&result,
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
					event.Match(),
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

				if err = ParticipationCollection.Aggregate(
					ctx,
					models.ParticipationPipeline().Match(bson.D{{Key: "event_id", Value: e.ID}}).Pipe,
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
	if err = models.DepositPermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	result = new([]models.Deposit)
	if err = DepositCollection.Aggregate(
		ctx,
		models.DepositPipeline().Match(filter).Pipe,
		result,
	); err != nil {
		return
	}
	return
}

func DepositGetByID(ctx context.Context, i *models.DepositParam, token *vcapool.AccessToken) (result *models.Deposit, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	if err = DepositCollection.AggregateOne(
		ctx,
		models.DepositPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}
