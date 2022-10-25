package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func TakingInsert(ctx context.Context, i *models.TakingCreate, token *vcapool.AccessToken) (r *models.Taking, err error) {
	taking := i.TakingDatabase()
	/*event := models.EventDatabase{
		ID:           uuid.NewString(),
		Name:         i.Name,
		CrewID:       i.CrewID,
		TypeOfEvent:  "automatically",
		StartAt:      time.Now().Unix(),
		EndAt:        time.Now().Unix(),
		TakingID:     taking.ID,
		EventASPID:   token.ID,
		InteralASPID: token.ID,
		CreatorID:    token.ID,
		ArtistIDs:    []string{},
		EventTools:   models.EventTools{},
	}*/
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}

	if i.NewSource != nil {
		sources := i.SourceList(taking.ID)
		if err = SourceCollection.InsertMany(ctx, sources.InsertMany()); err != nil {
			return
		}
	}
	/*
		if err = EventCollection.InsertOne(ctx, event); err != nil {
			return
		}*/
	r = new(models.Taking)
	if err = TakingCollection.AggregateOne(
		ctx,
		models.NewTakingsPipeline().Match(bson.D{{Key: "_id", Value: taking.ID}}).Pipe,
		r,
	); err != nil {
		return
	}
	return
}

func TakingUpdate(ctx context.Context, i *models.TakingUpdate) (r *models.Taking, err error) {
	takingDatabase := new(models.TakingDatabase)
	if err = TakingCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, takingDatabase); err != nil {
		return
	}
	i.State = takingDatabase.State
	for n := range i.NewSources {
		i.State.Open.Amount += i.NewSources[n].Money.Amount
	}
	sources := i.SourceList(takingDatabase.ID)
	if err = SourceCollection.InsertMany(ctx, sources.InsertMany()); err != nil {
		return
	}
	for m := range i.DeleteSources {
		deleteSource := new(models.Source)
		if err = SourceCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.DeleteSources[m]}}, deleteSource); err != nil {
			return
		}
		if deleteSource.Status == "open" {
			takingDatabase.State.Open.Amount -= deleteSource.Money.Amount
			if err = SourceCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: deleteSource.ID}}); err != nil {
				return
			}
		} else {
			continue
		}
	}
	for _, updateSource := range i.UpdateSource {
		databaseSource := new(models.Source)
		if err = SourceCollection.FindOne(
			ctx,
			bson.D{{Key: "_id", Value: updateSource.ID}},
			databaseSource,
		); err != nil {
			return
		}
		if databaseSource.Status == "open" {
			if updateSource.Money.Amount != databaseSource.Money.Amount {
				i.State.Open.Amount -= databaseSource.Money.Amount
				i.State.Open.Amount += updateSource.Money.Amount
			}
		}
		if err = SourceCollection.UpdateOne(
			ctx,
			bson.D{{Key: "_id", Value: updateSource.ID}},
			vmdb.UpdateSet(updateSource),
			nil,
		); err != nil {
			return
		}
	}
	r = new(models.Taking)
	if err = TakingCollection.UpdateOneAggregate(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(i),
		r,
		models.NewTakingsPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
	); err != nil {
		return
	}
	return
}
