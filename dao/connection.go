package dao

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

var MailSend = vcago.NewMailSend()

var Database = vmdb.NewDatabase("pool-user").Connect()

// UserCollection represents the database collection of the User model.
var UserCollection = Database.Collection("users").CreateIndex("email", true)

// UserCrewCollection represents the database collection of the UserCrew model.
var UserCrewCollection = Database.Collection("user_crew").CreateIndex("user_id", true)

// ActiveCollection represents the database collection of the Active model.
var ActiveCollection = Database.Collection("active").CreateIndex("user_id", true)

// NVMCollection represents the database collection of the NVM model.
var NVMCollection = Database.Collection("nvm").CreateIndex("user_id", true)

// AdressesCollection represents the database colltection of the Address model.
var AddressesCollection = Database.Collection("addresses").CreateIndex("user_id", true)

// CrewsCollection represents the database collection of the Crew model.
var CrewsCollection = Database.Collection("crews").CreateIndex("name", true)

// ProfilesCollection represents the database collection of the Profile model.
var ProfilesCollection = Database.Collection("profiles").CreateIndex("user_id", true)

// AvatarCollection represents the database collection of the Avatar model.
var AvatarCollection = Database.Collection("avatar").CreateIndex("user_id", true)

// PoolRoleCollection represents the database collection of the PoolRole Collection.
var PoolRoleCollection = Database.Collection("pool_roles").CreateMultiIndex(bson.D{{Key: "name", Value: 1}, {Key: "user_id", Value: 1}}, true)
