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

func InitialNats() {
	vcago.Nats.Connect()
	vcago.Nats.Subscribe("system.user.updated", SubscribeUserUpdate)
	vcago.Nats.Subscribe("system.user.import", SubscribeUserImport)
	vcago.Nats.Subscribe("system.user.deleted", SubscribeUserDelete)
	vcago.Nats.Subscribe("system.user.deletion_request", SubscribeUserDeleteRequest)
	vcago.Nats.Subscribe("system.user.deletion_process", SubscribeUserDeleteProcess)
	vcago.Nats.Subscribe("system.notification.publish", SubscribeNotificationPublish)
}

func PublishRoles() {
	result := vmod.WebappAccess{
		Name: "pool",
		Roles: []vmod.AccessRole{{
			Name: "pool_employee",
			Root: []string{"employee", "pool_employee"},
		}, {
			Name: "pool_finance",
			Root: []string{"employee", "pool_employee", "pool_finance"},
		}},
	}
	log.Print(result)
	vcago.Nats.Publish("webapp_role.update", result)
}

func SubscribeUserUpdate(m *models.UserUpdate) {
	result := new(models.User)
	if err := UserCollection.UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: m.ID}},
		vmdb.UpdateSet(m),
		result,
	); err != nil {
		if vmdb.ErrNoDocuments(err) {
			log.Print(err)
		}
	} else {
		vcago.Nats.Publish("pool.user.updated", result)
	}

}

func SubscribeUserImport(m *models.UserDatabase) {
	if m.DropsID != "" {
		if _, err := UserInsert(context.Background(), m); err != nil {
			log.Print(err)
		}
	}
}

func SubscribeUserDelete(m *vmod.DeletedResponse) {
	if err := UserDelete(context.Background(), m.ID); err != nil {
		log.Print(err)
	}
}

func SubscribeNotificationPublish(m *vcago.NotificationResponse) {
	var err error
	message := new(models.Message)
	user := new(models.User)
	if user, err = UsersGetByID(context.Background(), &models.UserParam{ID: m.User.ID}); err != nil {
		log.Print(err)
		return
	}
	result := message.NotificationMessage(m, user)
	if err := MessageCollection.InsertOne(context.Background(), result); err != nil {
		log.Print(err)
		return
	}
}

func SubscribeUserDeleteRequest(m *models.DeletionRequest) {
	data := make(map[string]string)
	var err error

	events := new([]models.EventUserDeletion)
	eventsFilter := bson.D{{Key: "external_asp.email", Value: m.Email}}

	user := new(models.User)
	if m.UserID != "" {
		if err = UserCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: m.UserID}}, user); err != nil {
			log.Print(err)
		} else {
			vmod.WalkStruct(user, "", &data)
			participations := new([]models.ParticipationUserDeletion)
			participationsFilter := bson.D{{Key: "user_id", Value: m.UserID}}
			if err = ParticipationEventViewCollection.Find(context.Background(), participationsFilter, participations); err != nil {
				log.Print(err)
			} else {
				vmod.WalkStruct(participations, "Participations", &data)
			}

			deposits := new([]models.DepositUserDeletion)
			depositsFilter := bson.D{{Key: "creator_id", Value: m.UserID}}
			if err = DepositCollection.Find(context.Background(), depositsFilter, deposits); err != nil {
				log.Print(err)
			} else {
				vmod.WalkStruct(deposits, "Deposits", &data)
			}

			eventsFilter =
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "event_asp_id", Value: m.UserID}},
					bson.D{{Key: "internal_asp_id", Value: m.UserID}},
					bson.D{{Key: "creator_id", Value: m.UserID}},
					bson.D{{Key: "external_asp.email", Value: m.Email}},
				},
				}}
		}

	}

	pipeline := models.EventUserPipeline().Match(eventsFilter).Pipe
	if err = EventCollection.Aggregate(context.Background(), pipeline, events); err != nil {
		log.Print(err)
	}
	matchedEvents := make([]models.EventUserDeletion, 0)
	for _, ev := range *events {
		// match based on EventASPID
		if ev.EventASPID != m.UserID {
			ev.EventASPID = ""
			ev.EventASP = models.User{}
		}

		// match based on InternalASPID
		if ev.InternalASPID != m.UserID {
			ev.InternalASPID = ""
			ev.InternalASP = models.User{}
		}

		// match based on CreatorID
		if ev.CreatorID != m.UserID {
			ev.CreatorID = ""
			ev.Creator = models.User{}
		}

		// match based on ExternalASP email
		if ev.ExternalASP.Email != m.Email {
			ev.ExternalASP = models.UserExternal{}
		}

		matchedEvents = append(matchedEvents, ev)
	}
	if len(matchedEvents) > 0 {
		vmod.WalkStruct(matchedEvents, "Events", &data)
	}

	statusMessage := "Data successfully provided"
	if len(data) == 0 && user.ID == "" {
		statusMessage = "No user data found"
	}

	response := &models.DeletionResponse{
		ID:            m.ID,
		UserID:        user.ID,
		Service:       "pool",
		StatusType:    "data_provided",
		StatusMessage: statusMessage,
		Data:          data,
	}
	vcago.Nats.Publish("system.user.deletion_response", response)
}

func SubscribeUserDeleteProcess(m *models.DeletionRequest) {
	var err error
	user := new(models.User)

	// Clear payment contact records
	if m.UserID != "" {
		if err = UserCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: m.UserID}}, user); err != nil {
			log.Print(err)
		} else {
			ClearUserDataOnDelete(context.Background(), m.UserID)
			// Clear users
			if err = UserCollection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: m.UserID}}); err != nil {
				log.Print(err)
			}
		}
	}

	eventsFilter := bson.D{{Key: "external_asp.email", Value: m.Email}}
	updateEvents := bson.D{{Key: "external_asp", Value: models.UserExternal{}}}
	if err = EventCollection.UpdateMany(context.Background(), eventsFilter, vmdb.UpdateSet(updateEvents)); err != nil {
		log.Print(err)
	}

	// Publish that service deleted user data
	vcago.Nats.Publish("system.user.user_deleted", &models.DeletionResponse{
		ID:            m.ID,
		UserID:        "",
		Service:       "pool",
		StatusType:    "data_deleted",
		StatusMessage: "Data successfully deleted",
	})
}
