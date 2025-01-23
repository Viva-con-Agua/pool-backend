package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var UserCrewCollection = "user_crew"

type (
	UserCrewCreate struct {
		CrewID string `json:"crew_id"`
	}
	UsersCrewCreate struct {
		CrewID         string `json:"crew_id"`
		UserID         string `json:"user_id"`
		OrganisationID string `bson:"organisation_id" json:"organisation_id"`
	}
	UserCrew struct {
		ID             string        `bson:"_id" json:"id"`
		UserID         string        `bson:"user_id" json:"user_id"`
		Name           string        `bson:"name" json:"name"`
		Email          string        `bson:"email" json:"email"`
		Roles          []vmod.Role   `bson:"roles" json:"roles"`
		CrewID         string        `bson:"crew_id" json:"crew_id"`
		OrganisationID string        `bson:"organisation_id" json:"organisation_id"`
		Organisation   string        `bson:"organisation" json:"organisation"`
		MailboxID      string        `bson:"mailbox_id" json:"mailbox_id"`
		Modified       vmod.Modified `bson:"modified" json:"modified"`
	}
	UserCrewMinimal struct {
		ID             string `bson:"_id" json:"id"`
		UserID         string `bson:"user_id" json:"user_id"`
		Name           string `bson:"name" json:"name"`
		Email          string `bson:"email" json:"email"`
		CrewID         string `bson:"crew_id" json:"crew_id"`
		OrganisationID string `bson:"organisation_id" json:"organisation_id"`
		Organisation   string `bson:"organisation" json:"organisation"`
	}
	UserCrewUpdate struct {
		ID             string `bson:"_id" json:"id"`
		UserID         string `bson:"user_id" json:"user_id"`
		Name           string `bson:"name" json:"name"`
		Email          string `bson:"email" json:"email"`
		CrewID         string `bson:"crew_id" json:"crew_id"`
		OrganisationID string `bson:"organisation_id" json:"organisation_id"`
	}
	UserCrewParam struct {
		ID string `param:"id"`
	}
	UserCrewImport struct {
		DropsID string   `json:"drops_id"`
		CrewID  string   `json:"crew_id"`
		NVMDate int64    `json:"nvm_date"`
		Created int64    `json:"created"`
		Active  string   `json:"active"`
		Roles   []string `json:"roles"`
	}
)

func NewUserCrew(userID string, crew *Crew) *UserCrew {
	return &UserCrew{
		ID:             uuid.NewString(),
		UserID:         userID,
		Name:           crew.Name,
		Email:          crew.Email,
		CrewID:         crew.ID,
		OrganisationID: crew.OrganisationID,
		MailboxID:      crew.MailboxID,
		Modified:       vmod.NewModified(),
	}
}

func (i *UserCrewUpdate) UserCrewUpdatePermission(token *AccessToken) (err error) {
	if token.ID != i.UserID {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}

func (i *UserCrewUpdate) UsersCrewUpdatePermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}

func (i *UsersCrewCreate) UsersCrewCreatePermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(CrewCollection)
	}
	return
}

func (i *UserCrewImport) ToActive(userID string) (result *Active) {
	result = NewActive(userID, i.CrewID)
	if i.Active != "" {
		if i.Active == "active" {
			result.Status = "confirmed"
		}
		if i.Active == "requested" {
			result.Status = "requested"
		}
		result.Since = i.Created
		result.Modified.Created = i.Created
	}
	return
}

func (i *UserCrewImport) ToNVM(userID string) (result *NVM) {
	result = NewNVM(userID)
	if i.NVMDate != 0 {
		result.Since = i.NVMDate
		result.Modified.Created = i.NVMDate
		result.Status = "confirmed"
	}
	return
}

func (i *UserCrewImport) ToRoles(userID string) (result []vmod.Role) {
	result = []vmod.Role{}
	for _, role := range i.Roles {
		switch role {
		case "asp":
			result = append(result, *RoleASP(userID))
		case "finance":
			result = append(result, *RoleFinance(userID))
		case "operation":
			result = append(result, *RoleAction(userID))
		case "education":
			result = append(result, *RoleEducation(userID))
		case "network":
			result = append(result, *RoleNetwork(userID))
		case "socialmedia":
			result = append(result, *RoleSocialMedia(userID))
		case "awareness":
			result = append(result, *RoleAwareness(userID))
		case "other":
			result = append(result, *RoleOther(userID))
		}
	}
	return
}

func (i *UserCrewCreate) CrewFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.CrewID)
	return filter.Bson()
}

func (i *UsersCrewCreate) CrewFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.CrewID)
	return filter.Bson()
}

func (i *UserCrewUpdate) CrewFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.CrewID)
	return filter.Bson()
}

func (i *UserCrewUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *UserCrewUpdate) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}
