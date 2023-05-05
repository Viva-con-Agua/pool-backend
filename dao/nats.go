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
