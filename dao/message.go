package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func MessageFilter(id string, crew *models.Crew, token *vcapool.AccessToken) bson.D {
	return bson.D{
		{Key: "_id", Value: id},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "mailbox_id", Value: token.MailboxID}},
			bson.D{{Key: "mailbox_id", Value: crew.MailboxID}},
		}}}
}

func MesseageCrewUser(ctx context.Context, i *models.RecipientGroup, token *vcapool.AccessToken) (result []models.TOData, err error) {
	if !token.PoolRoles.Validate(models.ASPRole) && !token.Roles.Validate("admin;employee") {
		return nil, vcago.NewBadRequest("message", "not allowed to send message")
	}
	filter := vmdb.NewFilter()
	filter.EqualString("crew.crew_id", i.CrewID)
	filter.EqualStringList("active.status", i.Active)
	filter.EqualStringList("nvm.status", i.NVM)
	if !i.IgnoreNewsletter {
		filter.EqualString("newsletter.value", "global")
	}
	userList := new([]models.User)
	if err = UserCollection.Aggregate(ctx, models.UserPipeline(false).Match(filter.Bson()).Pipe, userList); err != nil {
		return
	}
	result = []models.TOData{}
	for _, value := range *userList {
		result = append(result, models.TOData{UserID: value.ID, MailboxID: value.MailboxID, Email: value.Email})
	}
	return
}

func MessageEventUser(ctx context.Context, i *models.RecipientGroup, token *vcapool.AccessToken) (result []models.TOData, err error) {
	event := new(models.Event)
	filter := bson.D{{Key: "_id", Value: i.EventID}}
	if err = EventCollection.FindOne(ctx, filter, event); err != nil {
		return
	}
	if !token.PoolRoles.Validate(models.ASPRole) && !token.Roles.Validate("admin;employee") && event.EventASPID != token.ID {
		return nil, vcago.NewBadRequest("message", "not allowed to send message")
	}
	pFilter := bson.D{{Key: "event_id", Value: event.ID}}
	pPipeline := models.ParticipationPipeline().Match(pFilter)
	participations := new([]models.Participation)
	if err = ParticipationCollection.Aggregate(ctx, pPipeline.Pipe, participations); err != nil {
		return
	}
	result = []models.TOData{}
	for _, value := range *participations {
		result = append(result, models.TOData{UserID: value.User.ID, MailboxID: value.User.MailboxID, Email: value.User.Email})
	}
	return
}

func MessageSendCycular(ctx context.Context, i *vmod.IDParam, token *vcapool.AccessToken) (result *models.Message, mail *vcago.CycularMail, err error) {
	// get message via filter by mailbox and message ids
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
		return
	}
	filter := MessageFilter(i.ID, crew, token)
	result = new(models.Message)
	if err = MessageCollection.FindOne(ctx, filter, result); err != nil {
		return
	}
	//select TOData based on the recipientgroup
	if result.RecipientGroup.Type == "crew" {

		if result.To, err = MesseageCrewUser(ctx, &result.RecipientGroup, token); err != nil {
			return
		}
	} else if result.RecipientGroup.Type == "event" {
		if result.To, err = MessageEventUser(ctx, &result.RecipientGroup, token); err != nil {
			return
		}
	} else {
		return nil, nil, vcago.NewBadRequest("message", "type is not supported", result.RecipientGroup)
	}

	//create new cycular mail
	mail = vcago.NewCycularMail(result.From, result.ToEmails(), result.Subject, result.Message)
	if err = MessageCollection.InsertMany(ctx, *result.Inbox()); err != nil {
		return
	}
	result.Type = "outbox"
	if err = MessageCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: result.ID}},
		vmdb.UpdateSet(result.MessageUpdate()),
		result,
	); err != nil {
		return
	}
	return
}

func MessageUpdate(ctx context.Context, i *models.MessageUpdate, token *vcapool.AccessToken) (result *models.Message, err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: token.CrewID}}, crew); err != nil {
		return
	}
	filter := MessageFilter(i.ID, crew, token)
	result = new(models.Message)
	if err = MessageCollection.UpdateOne(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		result,
	); err != nil {
		return
	}
	return
}
