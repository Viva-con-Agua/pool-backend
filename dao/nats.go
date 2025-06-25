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
