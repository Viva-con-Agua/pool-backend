package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
)

func OrganisationInsert(ctx context.Context, i *models.OrganisationCreate, token *models.AccessToken) (result *models.Organisation, err error) {
	if err = token.AccessPermission(); err != nil {
		return
	}
	result = i.Organisation()
	if err = OrganisationCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func OrganisationGet(ctx context.Context, i *models.OrganisationQuery) (result *[]models.Organisation, err error) {
	filter := i.Filter()
	result = new([]models.Organisation)
	if err = OrganisationCollection.Find(ctx, filter, result); err != nil {
		return
	}
	return
}

func OrganisationGetByID(ctx context.Context, i *models.OrganisationParam) (result *models.Organisation, err error) {
	filter := i.Match()
	if err = OrganisationCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func OrganisationUpdate(ctx context.Context, i *models.OrganisationUpdate, token *models.AccessToken) (result *models.Organisation, err error) {
	if err = token.AccessPermission(); err != nil {
		return
	}
	filter := i.Match()
	if err = OrganisationCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), &result); err != nil {
		return
	}
	return
}

func OrganisationDelete(ctx context.Context, i *models.OrganisationParam, token *models.AccessToken) (err error) {
	if err = token.AccessPermission(); err != nil {
		return
	}
	filter := i.Match()
	if err = OrganisationCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}
