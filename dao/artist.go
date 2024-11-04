package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
)

func ArtistInsert(ctx context.Context, i *models.ArtistCreate, token *models.AccessToken) (result *models.Artist, err error) {
	if err = models.ArtistPermission(token); err != nil {
		return
	}
	result = i.Artist()
	if err = ArtistCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func ArtistGet(ctx context.Context, i *models.ArtistQuery) (result *[]models.Artist, err error) {
	filter := i.Filter()
	result = new([]models.Artist)
	if err = ArtistCollection.Find(ctx, filter, result); err != nil {
		return
	}
	return
}

func ArtistGetByID(ctx context.Context, i *models.ArtistParam) (result *models.Artist, err error) {
	filter := i.Match()
	if err = ArtistCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func ArtistUpdate(ctx context.Context, i *models.ArtistUpdate, token *models.AccessToken) (result *models.Artist, err error) {
	if err = models.ArtistPermission(token); err != nil {
		return
	}
	filter := i.Match()
	if err = ArtistCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	return
}

func ArtistDelete(ctx context.Context, i *models.ArtistParam, token *models.AccessToken) (err error) {
	if err = models.ArtistDeletePermission(token); err != nil {
		return
	}
	filter := i.Match()
	if err = ArtistCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}
