package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var UserCrewCollection = "user_crew"

type (
	UserCrewCreate struct {
		CrewID string `json:"crew_id"`
	}
	UserCrew struct {
		ID        string        `bson:"_id" json:"id"`
		UserID    string        `bson:"user_id" json:"user_id"`
		Name      string        `bson:"name" json:"name"`
		Email     string        `bson:"email" json:"email"`
		Roles     []vmod.Role   `bson:"roles" json:"roles"`
		CrewID    string        `bson:"crew_id" json:"crew_id"`
		MailboxID string        `bson:"mailbox_id" json:"mailbox_id"`
		Modified  vmod.Modified `bson:"modified" json:"modified"`
	}
	UserCrewMinimal struct {
		ID     string `bson:"_id" json:"id"`
		UserID string `bson:"user_id" json:"user_id"`
		Name   string `bson:"name" json:"name"`
		Email  string `bson:"email" json:"email"`
		CrewID string `bson:"crew_id" json:"crew_id"`
	}
	UserCrewUpdate struct {
		ID     string `bson:"_id" json:"id"`
		UserID string `bson:"user_id" json:"user_id"`
		Name   string `bson:"name" json:"name"`
		Email  string `bson:"email" json:"email"`
		CrewID string `bson:"crew_id" json:"crew_id"`
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

func NewUserCrew(userID string, crewID string, name string, email string, mailboxID string) *UserCrew {
	return &UserCrew{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      name,
		Email:     email,
		CrewID:    crewID,
		MailboxID: mailboxID,
		Modified:  vmod.NewModified(),
	}
}

func (i *UserCrewUpdate) UserCrewUpdatePermission(token *vcapool.AccessToken) (err error) {
	if token.ID != i.UserID {
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

func (i *UserCrewUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualString("user_id", token.ID)
	return filter.Bson()
}
