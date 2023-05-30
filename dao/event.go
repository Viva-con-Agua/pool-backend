package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func EventInsert(ctx context.Context, i *models.EventDatabase, token *vcapool.AccessToken) (result *models.Event, err error) {
	if err = models.EventPermission(token); err != nil {
		return
	}
	taking := models.TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "automatically",
		Modified: vmod.NewModified(),
	}
	i.TakingID = taking.ID
	if err = EventCollection.InsertOne(ctx, i); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivityDB(i.CreatorID, "event", i.ID, "Event created", "created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivityDB(i.CreatorID, "taking", taking.ID, "Taking automatically created for Event", "auto_created")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}
	if err = EventCollection.AggregateOne(ctx, models.EventPipeline(token).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func EventGet(ctx context.Context, i *models.EventQuery, token *vcapool.AccessToken) (result *[]models.ListEvent, err error) {
	filter := i.PermittedFilter(token)
	result = new([]models.ListEvent)
	if err = EventCollection.Aggregate(ctx, models.EventPipeline(token).Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func EventGetByID(ctx context.Context, i *models.EventParam, token *vcapool.AccessToken) (result *models.Event, err error) {
	filter := i.PermittedFilter(token)

	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipeline(token).Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventAspGetByID(ctx context.Context, i *models.EventParam, token *vcapool.AccessToken) (result *models.Event, err error) {
	filter := i.PermittedFilter(token)
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipeline(token).Match(filter).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventViewGetByID(ctx context.Context, i *models.EventParam) (result *models.EventPublic, err error) {
	if err = EventCollection.AggregateOne(
		ctx,
		models.EventPipelinePublic().Match(i.Match()).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func EventGetPublic(ctx context.Context, i *models.EventQuery) (result *[]models.EventPublic, err error) {
	filter := i.Match()
	result = new([]models.EventPublic)
	if err = EventCollection.Aggregate(ctx, models.EventPipelinePublic().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func EventsGetReceiverEvents(ctx context.Context, i *models.EventQuery, token *vcapool.AccessToken) (result *[]models.EventPublic, err error) {
	filter := i.FilterEmailEvents(token)
	result = new([]models.EventPublic)

	if err = EventCollection.Aggregate(ctx, models.EventPipelinePublic().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func EventGetAps(ctx context.Context, i *models.EventQuery, token *vcapool.AccessToken) (result *[]models.ListDetailsEvent, err error) {
	filter := i.FilterAsp(token)
	result = new([]models.ListDetailsEvent)
	if err = EventCollection.Aggregate(ctx, models.EventPipeline(token).Match(filter).Pipe, result); err != nil {
		return
	}

	return
}

func EventUpdate(ctx context.Context, i *models.EventUpdate, token *vcapool.AccessToken) (result *models.Event, err error) {
	filter := i.PermittedFilter(token)
	if err = EventCollection.UpdateOneAggregate(
		ctx,
		filter,
		vmdb.UpdateSet(i),
		&result,
		models.EventPipeline(token).Match(i.Match()).Pipe,
	); err != nil {
		return
	}
	return
}

func EventDelete(ctx context.Context, i *models.EventParam, token *vcapool.AccessToken) (err error) {
	if err = models.EventDeletePermission(token); err != nil {
		return
	}
	if err = EventCollection.DeleteOne(ctx, i.Match()); err != nil {
		return
	}
	if err = ParticipationCollection.TryDeleteMany(ctx, bson.D{{Key: "event_id", Value: i.ID}}); err != nil {
		return
	}
	return
}

func EventImport(ctx context.Context, event *models.EventImport) (result *models.Event, err error) {
	i := event.EventDatabase()

	taking := models.TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "automatically",
		Modified: vmod.NewModified(),
	}
	i.TakingID = taking.ID

	admin := new(models.UserDatabase)
	adminFilter := bson.D{{Key: "email", Value: "f.wittmann@vivaconagua.org"}}
	if err = UserCollection.FindOne(
		ctx,
		adminFilter,
		admin,
	); err != nil {
		return
	}
	i.CreatorID = admin.ID
	i.InternalASPID = admin.ID
	i.EventState.InternalConfirmation = admin.ID

	if i.CrewID != "" {
		aspRole := new(models.RoleDatabase)
		if err = UserCrewCollection.AggregateOne(ctx, models.EventRolePipeline().Match(bson.D{{Key: "crew_id", Value: i.CrewID}}).Pipe, aspRole); err != nil {
			return
		}
		i.EventASPID = aspRole.UserID
	} else {
		i.EventASPID = admin.ID
	}

	if err = EventCollection.InsertOne(ctx, i); err != nil {
		return
	}
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	eventActivity := models.NewActivityDB(i.CreatorID, "event", i.ID, "Event created", "created")
	if err = ActivityCollection.InsertOne(ctx, eventActivity); err != nil {
		return
	}
	takingActivity := models.NewActivityDB(i.CreatorID, "taking", taking.ID, "Taking automatically created for Event", "auto_created")
	if err = ActivityCollection.InsertOne(ctx, takingActivity); err != nil {
		return
	}

	if err = EventCollection.AggregateOne(ctx, models.EventImportPipeline().Match(i.Match()).Pipe, &result); err != nil {
		return
	}

	// Add participations
	for _, drops_user := range event.Participations {
		participant := new(models.UserDatabase)
		participantFilter := bson.D{{Key: "drops_id", Value: drops_user.DropsID}}
		if err = UserCollection.FindOne(
			ctx,
			participantFilter,
			participant,
		); err != nil {
			return
		}

		participation := drops_user.ParticipationDatabase()
		participation.EventID = result.ID
		participation.UserID = participant.ID
		participation.CrewID = result.CrewID

		if err = ParticipationCollection.InsertOne(ctx, participation); err != nil {
			return
		}
	}

	return
}
