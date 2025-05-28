package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

func validateDepositUnits(ctx context.Context, takingID string, amount int64, crewID string, token *models.AccessToken) (err error) {
	taking := new(models.Taking)
	takingPipeline := models.TakingPipeline()
	if err = TakingCollection.AggregateOne(
		ctx,
		takingPipeline.Match(models.Match(takingID)).Pipe,
		taking,
	); err != nil {
		return
	}
	if amount > taking.Money.Amount {
		return vcago.NewBadRequest(models.DepositCollection, "taking_amount_failure", nil)
	}
	if (!token.Roles.Validate("admin;employee;pool_employee") && crewID != token.CrewID) || taking.CrewID != crewID {
		return vcago.NewBadRequest(models.DepositCollection, "taking_crew_failure", nil)
	}
	return
}

func DepositInsert(ctx context.Context, i *models.DepositCreate, token *models.AccessToken) (result *models.Deposit, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	deposit, depositUnits := i.DepositDatabase(token)
	for _, unit := range depositUnits {
		if err = validateDepositUnits(ctx, unit.TakingID, unit.Money.Amount, deposit.CrewID, token); err != nil {
			return
		}
	}
	deposit.ReasonForPayment, err = GetNewReasonForPayment(ctx, i.CrewID)
	if err != nil {
		return
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

func DepositUpdate(ctx context.Context, i *models.DepositUpdate, token *models.AccessToken) (result *models.Deposit, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	deposit := new(models.Deposit)
	filter := bson.D{{Key: "_id", Value: i.ID}}
	if err = DepositCollection.AggregateOne(
		ctx,
		models.DepositPipeline().Match(filter).Pipe,
		deposit,
	); err != nil {
		return
	}
	i.Money = deposit.Money
	if deposit.Status == "confirmed" && !token.Roles.Validate("admin;employee;pool_employee") {
		return nil, vcago.NewBadRequest("deposit", "deposit_confirmed_failure", nil)
	}
	depositUpdate, depositUnitCreate, depositUnitUpdate, depositUnitDelete := i.DepositDatabase(deposit)
	for _, unit := range i.DepositUnit {
		if err = validateDepositUnits(ctx, unit.TakingID, unit.Money.Amount, depositUpdate.CrewID, token); err != nil {
			return
		}
	}

	for _, unit := range depositUnitCreate {
		if err = DepositUnitCollection.InsertOne(ctx, unit); err != nil {
			return
		}
	}
	for _, unit := range depositUnitUpdate {
		updateFilter := bson.D{{Key: "_id", Value: unit.ID}}
		if err = DepositUnitCollection.UpdateOne(ctx, updateFilter, vmdb.UpdateSet(unit), nil); err != nil {
			return
		}
	}
	for _, unit := range depositUnitDelete {
		deleteFilter := bson.D{{Key: "_id", Value: unit.ID}}
		if err = DepositUnitCollection.DeleteOne(ctx, deleteFilter); err != nil {
			return
		}
	}
	if err = DepositCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, vmdb.UpdateSet(depositUpdate), nil); err != nil {
		return
	}
	if err = DepositCollection.AggregateOne(
		ctx,
		models.DepositPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
		&result,
	); err != nil {
		return
	}

	ctx = context.Background()
	if i.Status == "confirmed" {
		go func() {
			ctxAsync := context.Background()
			for _, unit := range i.DepositUnit {
				event := new(models.EventUpdate)
				if err = EventCollection.FindOne(
					ctxAsync,
					bson.D{{Key: "taking_id", Value: unit.TakingID}},
					event,
				); err != nil {
					if !vmdb.ErrNoDocuments(err) {
						return
					}
					err = nil
				}
				if event.ID != "" {
					event.EventState.State = "closed"
					e := new(models.Event)
					if err = EventCollection.UpdateOneAggregate(
						ctxAsync,
						event.Match(),
						vmdb.UpdateSet(event),
						e,
						models.EventPipeline(token).Match(event.Match()).Pipe,
					); err != nil {
						return
					}

					// Update CRM event
					if err = IDjango.Post(e, "/v1/pool/event/update/"); err != nil {
						log.Print(err)
					}

					// Add takings to CRM
					var taking *models.Taking
					if taking, err = TakingGetByID(ctx, &vmod.IDParam{ID: unit.TakingID}, token); err != nil {
						log.Print(err)
					}

					taking.EditorID = token.ID
					if err = IDjango.Post(taking, "/v1/pool/taking/create/"); err != nil {
						log.Print(err)
					}

					// Add participations to event
					participations := new([]models.Participation)

					if err = ParticipationCollection.Aggregate(
						ctxAsync,
						models.ParticipationPipeline().Match(bson.D{{Key: "event_id", Value: e.ID}}).Pipe,
						participations,
					); err != nil {
						return
					}

					if err = IDjango.Post(participations, "/v1/pool/participations/create/"); err != nil {
						log.Print(err)
						err = nil
					}

				}
			}
		}()

	}

	return
}

func DepositSync(ctx context.Context, i *models.DepositParam, token *models.AccessToken) (result *models.Deposit, err error) {

	filter := bson.D{{Key: "_id", Value: i.ID}}
	if err = DepositCollection.AggregateOne(
		ctx,
		models.DepositPipeline().Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	if result.Status != "confirmed" {
		return nil, vcago.NewBadRequest("deposit", "deposit_confirmed_failure", nil)
	}
	ctx = context.Background()
	go func() {
		for _, unit := range result.DepositUnit {
			event := new(models.EventUpdate)
			if err = EventCollection.FindOne(
				ctx,
				bson.D{{Key: "taking_id", Value: unit.TakingID}},
				event,
			); err != nil {
				if !vmdb.ErrNoDocuments(err) {
					return
				}
				err = nil
			}
			if event.ID != "" {
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

				// Update CRM event
				if err = IDjango.Post(e, "/v1/pool/event/update/"); err != nil {
					log.Print(err)
				}

				// Add takings to CRM
				var taking *models.Taking
				if taking, err = TakingGetByID(ctx, &vmod.IDParam{ID: unit.TakingID}, token); err != nil {
					log.Print(err)
				}
				taking.EditorID = token.ID
				if err = IDjango.Post(taking, "/v1/pool/taking/create/"); err != nil {
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
					err = nil
				}

			}
		}
	}()

	return
}

func DepositGet(ctx context.Context, query *models.DepositQuery, token *models.AccessToken) (result *[]models.Deposit, listSize int64, err error) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	filter := query.PermittedFilter(token)
	sort := query.Sort()
	result = &[]models.Deposit{}
	if err = DepositCollection.Aggregate(
		ctx,
		models.DepositPipelineList().SortFields(sort).Match(filter).Sort(sort).Skip(query.Skip, 0).Limit(query.Limit, 100).Pipe,
		result,
	); err != nil {
		return
	}
	count := new([]Count)
	if err = DepositCollection.Aggregate(
		context.Background(),
		models.DepositPipelineCount(filter).Pipe,
		count,
	); err != nil {
		return
	}
	if len(*count) == 0 {
		listSize = 0
	} else {
		listSize = (*count)[0].ListSize
	}
	return
}

func DepositGetByID(ctx context.Context, i *models.DepositParam, token *models.AccessToken) (result *models.Deposit, err error) {
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
