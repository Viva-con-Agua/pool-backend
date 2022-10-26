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
	i.State = &takingDatabase.State
	for _, v := range i.Sources {
		//create new sources
		if v.ID == "" {
			i.State.Open.Amount += v.Money.Amount
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
			if deleteSource.Status == "open" {
				takingDatabase.State.Open.Amount -= deleteSource.Money.Amount
				if err = SourceCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: deleteSource.ID}}); err != nil {
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
			if databaseSource.Status == "open" {
				if v.Money.Amount != databaseSource.Money.Amount {
					i.State.Open.Amount -= databaseSource.Money.Amount
					i.State.Open.Amount += v.Money.Amount
				}
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
		models.NewTakingsPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
	); err != nil {
		return
	}
	return
}
