package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type AddressCreate struct {
	vcapool.AddressCreate
}

type AddressParam struct {
	vcapool.AddressParam
}

type AddressUpdate struct {
	vcapool.AddressUpdate
}

type Address vcapool.Address

var AddressesCollection = Database.Collection("addresses").CreateIndex("user_id", true)

func (i *AddressCreate) Create(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Address, err error) {
	r = i.Address(token)
	err = AddressesCollection.InsertOne(ctx, &r)
	return
}

func (i *AddressParam) Get(ctx context.Context) (r *vcapool.Address, err error) {
	r = new(vcapool.Address)
	err = AddressesCollection.FindOne(ctx, i.Filter(), r)
	return
}

func (i *AddressUpdate) Update(ctx context.Context) (r *vcapool.Address, err error) {
	if err = AddressesCollection.UpdateOneSet(ctx, bson.M{"_id": i.ID}, i); err != nil {
		return
	}
	r = new(vcapool.Address)
	err = AddressesCollection.FindOne(ctx, bson.M{"_id": i.ID}, r)
	return
}

func (i *AddressParam) Delete(ctx context.Context) (err error) {
	err = AddressesCollection.DeleteOne(ctx, i.Filter())
	return
}

type AddressQuery struct {
	ID string `query:"id"`
}

func (i *AddressQuery) Filter() bson.M {
	f := vcago.NewMongoFilter()
	f.Equal("_id", i.ID)
	return f.Filter
}

type AddressList []Address

func (i *AddressList) Get(ctx context.Context, filter bson.M) (err error) {
	err = AddressesCollection.Find(ctx, filter, &i)
	return
}

//func (i *)
