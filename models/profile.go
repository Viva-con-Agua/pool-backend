package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ProfileCreate struct {
		Gender            string `bson:"gender" json:"gender"`
		Phone             string `bson:"phone" json:"phone"`
		Mattermost        string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate         int64  `bson:"birthdate" json:"birthdate"`
		BirthdateDatetime string `bson:"birthdate_datetime" json:"birthdate_datetime"`
	}
	ProfileUpdate struct {
		ID                string `bson:"_id" json:"id"`
		Gender            string `bson:"gender" json:"gender"`
		Phone             string `bson:"phone" json:"phone"`
		Mattermost        string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate         int64  `bson:"birthdate" json:"birthdate"`
		BirthdateDatetime string `bson:"birthdate_datetime" json:"birthdate_datetime"`
	}
	Profile struct {
		ID                string        `bson:"_id" json:"id"`
		Gender            string        `bson:"gender" json:"gender"`
		Phone             string        `bson:"phone" json:"phone"`
		Mattermost        string        `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate         int64         `bson:"birthdate" json:"birthdate"`
		BirthdateDatetime string        `bson:"birthdate_datetime" json:"birthdate_datetime"`
		UserID            string        `bson:"user_id" json:"user_id"`
		Modified          vmod.Modified `bson:"modified" json:"modified"`
	}
	ProfileParam struct {
		ID string `param:"id"`
	}
	ProfileMinimal struct {
		Mattermost        string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate         int64  `bson:"birthdate" json:"birthdate"`
		BirthdateDatetime string `bson:"birthdate_datetime" json:"birthdate_datetime"`
		UserID            string `bson:"user_id" json:"user_id"`
	}
	ProfileImport struct {
		Gender     string `bson:"gender" json:"gender"`
		Phone      string `bson:"phone" json:"phone"`
		Mattermost string `bson:"mattermost_username" json:"mattermost_username"`
		Birthdate  int64  `bson:"birthdate" json:"birthdate"`
		DropsID    string `bson:"drops_id" json:"drops_id"`
	}
)

var ProfileCollection = "profiles"

func (i *ProfileParam) ProfileSyncPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(ProfileCollection)
	}
	return
}

func (i *ProfileCreate) Profile(userID string) *Profile {
	birthdate := time.Unix(i.Birthdate, 0)
	birthdateDatetime := ""
	if i.Birthdate != 0 {
		birthdateDatetime = birthdate.Format("2006-01-02")
	}

	return &Profile{
		ID:                uuid.NewString(),
		Gender:            i.Gender,
		Phone:             i.Phone,
		Mattermost:        i.Mattermost,
		Birthdate:         i.Birthdate,
		BirthdateDatetime: birthdateDatetime,
		UserID:            userID,
		Modified:          vmod.NewModified(),
	}
}

func (i *ProfileImport) Profile(userID string) *Profile {
	return &Profile{
		ID:         uuid.NewString(),
		Gender:     i.Gender,
		Phone:      i.Phone,
		Mattermost: i.Mattermost,
		Birthdate:  i.Birthdate,
		UserID:     userID,
		Modified:   vmod.NewModified(),
	}
}

func (i *ProfileUpdate) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}

func (i *ProfileUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ProfileParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}
