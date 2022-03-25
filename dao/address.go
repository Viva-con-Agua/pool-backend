package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Address vcapool.Address

var AddressesCollection = Database.Collection("addresses").CreateIndex("user_id", true)

func (i *Address) Create(ctx context.Context) (err error) {
	i.ID = uuid.NewString()
	i.Modified = vcago.NewModified()
	err = AddressesCollection.InsertOne(ctx, &i)
	return
}

func (i *Address) Get(ctx context.Context, filter bson.M) (err error) {
	err = AddressesCollection.FindOne(ctx, filter, i)
	return
}

func (i *Address) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	update := bson.M{"$set": &i}
	err = AddressesCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, update)
	return
}

func (i *Address) Delete(ctx context.Context, filter bson.M) (err error) {
	err = AddressesCollection.DeleteOne(ctx, filter)
	return
}

type AddressQuery struct {
	ID string `query:"id"`
}

func (i *AddressQuery) Filter() bson.M {
	f := vcago.NewMongoFilterM()
	f.Equal("_id", i.ID)
	return f.Filter
}

type AddressList []Address

func (i *AddressList) Get(ctx context.Context, filter bson.M) (err error) {
	err = AddressesCollection.Find(ctx, filter, &i)
	return
}
