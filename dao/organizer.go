package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
)

func OrganizerInsert(ctx context.Context, i *models.OrganizerCreate, token *vcapool.AccessToken) (result *models.Organizer, err error) {
	if err = models.OrganizerPermission(token); err != nil {
		return
	}
	result = i.Organizer()
	if err = OrganizerCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func OrganizerGet(ctx context.Context, i *models.OrganizerQuery) (result *[]models.Organizer, err error) {
	filter := i.Filter()
	result = new([]models.Organizer)
	if err = OrganizerCollection.Find(ctx, filter, result); err != nil {
		return
	}
	return
}

func OrganizerGetByID(ctx context.Context, i *models.OrganizerParam) (result *models.Organizer, err error) {
	filter := i.Match()
	if err = OrganizerCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func OrganizerUpdate(ctx context.Context, i *models.OrganizerUpdate, token *vcapool.AccessToken) (result *models.Organizer, err error) {
	if err = models.OrganizerPermission(token); err != nil {
		return
	}
	filter := i.Match()
	if err = OrganizerCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	return
}

func OrganizerDelete(ctx context.Context, i *models.OrganizerParam, token *vcapool.AccessToken) (err error) {
	if err = models.OrganizerDeletePermission(token); err != nil {
		return
	}
	filter := i.Match()
	if err = OrganizerCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}
