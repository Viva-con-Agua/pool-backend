package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

/*
type Address struct {
	ID          string         `json:"_id" bson:"id"`
	Street      string         `json:"street" bson:"street"`
	Number      string         `json:"number" bson:"number"`
	Zip         string         `json:"zip" bson:"zip"`
	City        string         `json:"city" bson:"city"`
	Country     string         `json:"country" bson:"country"`
	Additionals string         `json:"additionals" bson:"additionals"`
	UserID      string         `json:"user_id" bson:"user_id"`
	Modified    vcago.Modified `json:"modified" bson:"modified"`
}*/

type Address vcapool.Address

/*
func (i *Address) ToVca() *vcapool.Address {
	return &vcapool.Address{
		ID:          i.ID,
		Street:      i.Street,
		Number:      i.Number,
		Zip:         i.Zip,
		City:        i.City,
		Country:     i.Country,
		Additionals: i.Additionals,
		UserID:      i.UserID,
		Modified:    i.Modified,
	}
}*/

var AddressesCollection = Database.Collection("addresses").CreateIndex("user_id", true)

func (i *Address) Create(ctx context.Context) (err error) {
	i.ID = uuid.NewString()
	i.Modified = vcago.NewModified()
	err = AddressesCollection.InsertOne(ctx, &i)
	return
}

func (i *Address) Get(ctx context.Context, id string) (err error) {
	err = AddressesCollection.FindOne(ctx, bson.M{"_id": id}, &id)
	return
}

func (i *Address) Update(ctx context.Context) (err error) {
	i.Modified.Update()
	err = AddressesCollection.UpdateOne(ctx, bson.M{"_id": i.ID}, &i)
	return
}

func (i *Address) Delete(ctx context.Context) (err error) {
	err = AddressesCollection.DeleteOne(ctx, bson.M{"_id": i.ID})
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
