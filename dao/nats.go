package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func InitialNats() {
	vcago.Nats.Connect()
	vcago.Nats.Subscribe("system.user.updated", SubscribeUserUpdate)
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
