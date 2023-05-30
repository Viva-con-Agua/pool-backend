package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO CHECK ALL INSERTS FOR USAGE OF CREATION TYPES WITH MODIFIED
func CrewInsert(ctx context.Context, i *models.CrewCreate, token *vcapool.AccessToken) (result *models.Crew, err error) {

	if err = models.CrewPermission(token); err != nil {
		return
	}
	//create mailbox
	mailbox := models.NewMailboxDatabase("crew")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	result = i.Crew()
	// refer the mailbox.ID
	result.MailboxID = mailbox.ID
	// insert user
	if err = CrewsCollection.InsertOne(ctx, result); err != nil {
		return
	}
	//select user from database
	return
}

func CrewGet(ctx context.Context, i *models.CrewQuery, token *vcapool.AccessToken) (result *[]models.Crew, err error) {
	if err = models.CrewPermission(token); err != nil {
		return
	}
	filter := i.Filter()
	result = new([]models.Crew)
	if err = CrewsCollection.Find(ctx, filter, result); err != nil {
		return
	}

	return
}

func CrewGetByID(ctx context.Context, i *models.CrewParam, token *vcapool.AccessToken) (result *models.Crew, err error) {
	result = new(models.Crew)
	filter := i.PermittedFilter(token)
	if err = CrewsCollection.FindOne(ctx, filter, result); err != nil {
		return
	}
	return
}

func CrewPublicGet(ctx context.Context, i *models.CrewQuery) (result *[]models.CrewPublic, err error) {
	result = new([]models.CrewPublic)
	filter := i.ActiveFilter()
	if err = CrewsCollection.Find(ctx, filter, result); err != nil {
		return
	}

	return
}

func CrewGetAsMember(ctx context.Context, i *models.CrewQuery, token *vcapool.AccessToken) (result *models.Crew, err error) {
	result = new(models.Crew)
	filter := i.PermittedFilter(token)
	if err = CrewsCollection.FindOne(ctx, filter, result); err != nil {
		return
	}
	return
}

// TODO CHECK ALL UPDATE METHODS IN ALL FILES FOR USAGE OF vmdb.UpdateSet()
func CrewUpdate(ctx context.Context, i *models.CrewUpdate, token *vcapool.AccessToken) (result *models.Crew, err error) {
	if err = models.CrewUpdatePermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	result = new(models.Crew)
	if !token.Roles.Validate("employee;admin") {
		if err = CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i.ToCrewUpdateASP()), result); err != nil {
			return
		}
	} else {
		if err = CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), result); err != nil {
			return
		}
	}
	return
}

func CrewDelete(ctx context.Context, i *models.CrewParam, token *vcapool.AccessToken) (err error) {
	if err = models.CrewPermission(token); err != nil {
		return
	}
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, crew); err != nil {
		return
	}
	if err = MailboxCollection.TryDeleteOne(ctx, bson.D{{Key: "_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = MessageCollection.TryDeleteMany(ctx, bson.D{{Key: "mailbox_id", Value: crew.MailboxID}}); err != nil {
		return
	}
	if err = CrewsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: i.ID}}); err != nil {
		return
	}
	return
}

func CrewImport(ctx context.Context, i *models.CrewCreate) (r *models.Crew, err error) {
	//create mailbox
	mailbox := models.NewMailboxDatabase("crew")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	r = i.Crew()
	// refer the mailbox.ID
	r.MailboxID = mailbox.ID
	// insert user
	if err = CrewsCollection.InsertOne(ctx, r); err != nil {
		return
	}
	//select user from database
	return
}
