package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AddressCreate struct {
		vcapool.AddressCreate
	}

	AddressParam struct {
		vcapool.AddressParam
	}
	AddressQuery struct {
		vcapool.AddressQuery
	}

	AddressUpdate struct {
		vcapool.AddressUpdate
	}
	Address vcapool.Address
)

var AddressesCollection = Database.Collection("addresses").CreateIndex("user_id", true)

func (i *AddressCreate) Create(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Address, err error) {
	r = i.Address(token)
	err = AddressesCollection.InsertOne(ctx, &r)
	return
}

func (i *AddressParam) Get(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Address, err error) {
	r = new(vcapool.Address)
	if token.Roles.Validate("employee;admin") {
		err = AddressesCollection.FindOne(ctx, bson.M{"_id": i.ID}, r)
	} else {
		err = AddressesCollection.FindOne(ctx, bson.M{"_id": i.ID, "user_id": token.ID}, r)
	}
	return
}

func (i *AddressUpdate) Update(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.Address, err error) {
	if err = AddressesCollection.UpdateOneSet(ctx, bson.M{"_id": i.ID, "user_id": token.ID}, i.AddressUpdate); err != nil {
		return
	}
	r = new(vcapool.Address)
	err = AddressesCollection.FindOne(ctx, bson.M{"_id": i.ID}, r)
	return
}

func (i *AddressParam) Delete(ctx context.Context, token *vcapool.AccessToken) (err error) {
	err = AddressesCollection.DeleteOne(ctx, bson.M{"_id": i.ID, "user_id": token.ID})
	return
}

func (i *AddressQuery) List(ctx context.Context, token *vcapool.AccessToken) (r *vcapool.AddressList, err error) {
	if !token.Roles.Validate("employee;admin") {
		err = vcago.NewPermissionDenied("address")
		return
	}
	pipe := vcago.NewMongoPipe()
	pipe.Match(i.Match())
	r = new(vcapool.AddressList)
	err = AddressesCollection.Aggregate(ctx, pipe.Pipe, &i)
	return
}

//func (i *)
