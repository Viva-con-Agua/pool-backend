package models

import "github.com/Viva-con-Agua/vcago/vmod"

type (
	UserMigrate struct {
		ID        string        `json:"id" bson:"_id"`
		Email     string        `json:"email" bson:"email"`
		FirstName string        `json:"first_name" bson:"first_name"`
		LastName  string        `json:"last_name" bson:"last_name"`
		FullName  string        `json:"full_name" bson:"full_name"`
		DropsID   string        `json:"drops_id" bson:"drops_id"`
		PoolID    string        `json:"pool_id" bson:"pool_id"`
		Crew      CrewMigrate   `json:"crew" bson:"crew"`
		Active    string        `json:"active" bson:"active"`
		NVMDate   int64         `json:"nvm_state" bson:"nvm_state"`
		Address   AddressCreate `json:"address" bson:"address"`
		Profile   ProfileCreate `json:"profile" bson:"profile"`
		Roles     []string      `json:"roles" bson:"roles"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
	CrewMigrate struct {
		ID       string        `json:"id" bson:"id"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
)

/*
func (i *UserMigrate) UserActive() *ActiveUpdate {
	if i.Active == "active" {
		return
	} else if i.Active == "requested" {
		return &vcapool.UserActiveUpdate{
			Status: "requested",
			Since:  time.Now().Unix(),
		}
	}
	return nil
}

func (i *UserMigrate) UserNVM(userID string) *vcapool.UserNVM {
	if i.NVMDate != 0 {
		return &vcapool.UserNVM{
			ID:       uuid.NewString(),
			Status:   "confirmed",
			Since:    i.NVMDate,
			Expired:  0,
			UserID:   userID,
			Modified: vcago.NewModified(),
		}
	}
	return nil
}*/
