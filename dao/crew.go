package dao

import (
	"context"
	"pool-backend/models"
	"sort"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	opt := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	opt.Collation = &options.Collation{Locale: "en", Strength: 2}

	if err = CrewsCollection.Aggregate(ctx, models.CrewPipeline().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func CrewGetByID(ctx context.Context, i *models.CrewParam, token *vcapool.AccessToken) (result *models.Crew, err error) {
	filter := i.PermittedFilter(token)
	if err = CrewsCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func CrewPublicGet(ctx context.Context, i *models.CrewQuery) (result *[]models.CrewPublic, err error) {
	filter := i.ActiveFilter()
	result = new([]models.CrewPublic)
	if err = CrewsCollection.Find(ctx, filter, result); err != nil {
		return
	}
	return
}

func CrewGetAsMember(ctx context.Context, i *models.CrewQuery, token *vcapool.AccessToken) (result *models.Crew, err error) {
	filter := i.PermittedFilter(token)
	if err = CrewsCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	return
}

func CrewUpdate(ctx context.Context, i *models.CrewUpdate, token *vcapool.AccessToken) (result *models.Crew, err error) {
	if err = models.CrewUpdatePermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, filter, &crew); err != nil {
		return
	}
	// Its not allowed to set the asp_selection to "selected" manually
	if crew.AspSelection != "selected" && i.AspSelection == "selected" {
		return nil, vcago.NewBadRequest(models.CrewCollection, "It is not allowed to set the asp selection state to selected manually!")
	}
	strings := []string{"active", "inactive"}
	sort.Strings(strings)
	match := sort.SearchStrings(strings, i.AspSelection)
	if crew.AspSelection == "selected" && match < len(strings) && strings[match] == i.AspSelection {
		RoleHistoryDelete(ctx, &models.RoleHistoryRequest{CrewID: i.ID, Confirmed: false}, token)
	}
	if !token.Roles.Validate("employee;admin") {
		if err = CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i.ToCrewUpdateASP()), &result); err != nil {
			return
		}
	} else {
		if err = CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(i), &result); err != nil {
			return
		}
	}
	if crew.Email != i.Email || crew.Name != i.Name {
		filter := bson.D{{Key: "crew_id", Value: i.ID}}
		update := bson.D{{Key: "email", Value: i.Email}, {Key: "name", Value: i.Name}}
		if err = UserCrewCollection.UpdateMany(ctx, filter, vmdb.UpdateSet(update)); err != nil {
			return
		}

	}
	return
}

func CrewUpdateAspSelection(ctx context.Context, i *models.CrewParam, value string, token *vcapool.AccessToken) (result *models.Crew, err error) {
	if err = models.CrewUpdatePermission(token); err != nil {
		return
	}
	filter := i.PermittedFilter(token)
	crew := new(models.CrewUpdate)
	if err = CrewsCollection.FindOne(ctx, filter, &crew); err != nil {
		return
	}
	crew.AspSelection = value
	if err = CrewsCollection.UpdateOne(ctx, filter, vmdb.UpdateSet(crew), &result); err != nil {
		return
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

func CrewImport(ctx context.Context, i *models.CrewCreate) (result *models.Crew, err error) {
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
