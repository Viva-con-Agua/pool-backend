package models

import (
	"strconv"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleHistoryCreate struct {
	UserID    string        `json:"user_id"`
	Role      string        `json:"role"`
	CrewID    string        `json:"crew_id"`
	Confirmed bool          `json:"confirmed" bson:"confirmed"`
	StartDate int64         `json:"start_date" bson:"start_date"`
	EndDate   int64         `json:"end_date" bson:"end_date"`
	Modified  vmod.Modified `json:"modified" bson:"modified"`
}

type RoleHistoryUpdate struct {
	ID        string `bson:"_id" json:"id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Role      string `json:"role" bson:"role"`
	CrewID    string `json:"crew_id" bson:"crew_id"`
	Confirmed bool   `json:"confirmed" bson:"confirmed"`
	StartDate int64  `json:"start_date" bson:"start_date"`
	EndDate   int64  `json:"end_date" bson:"end_date"`
}
type RoleHistoryRequest struct {
	UserID    string        `json:"user_id"`
	Role      string        `json:"role"`
	CrewID    string        `json:"crew_id"`
	Confirmed bool          `json:"confirmed" bson:"confirmed"`
	StartDate int64         `json:"start_date" bson:"start_date"`
	EndDate   int64         `json:"end_date" bson:"end_date"`
	Modified  vmod.Modified `json:"modified" bson:"modified"`
}

type RoleHistoryBulkRequest struct {
	CrewID     string        `json:"crew_id"`
	AddedRoles []RoleRequest `json:"created"`
}

type RoleHistory struct {
	ID        string          `bson:"_id" json:"id"`
	UserID    string          `json:"user_id" bson:"user_id"`
	Role      string          `json:"role" bson:"role"`
	CrewID    string          `json:"crew_id" bson:"crew_id"`
	Crew      UserCrewMinimal `json:"crew" bson:"crew"`
	Profile   ProfileMinimal  `json:"profile" bson:"profile"`
	User      UserMinimal     `json:"user" bson:"user"`
	Confirmed bool            `json:"confirmed" bson:"confirmed"`
	StartDate int64           `json:"start_date" bson:"start_date"`
	EndDate   int64           `json:"end_date" bson:"end_date"`
	Modified  vmod.Modified   `json:"modified" bson:"modified"`
}

type RoleHistoryDatabase struct {
	ID        string        `bson:"_id" json:"id"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Role      string        `json:"role" bson:"role"`
	CrewID    string        `json:"crew_id" bson:"crew_id"`
	Confirmed bool          `json:"confirmed" bson:"confirmed"`
	StartDate int64         `json:"start_date" bson:"start_date"`
	EndDate   int64         `json:"end_date" bson:"end_date"`
	Modified  vmod.Modified `json:"modified" bson:"modified"`
}

var PoolRoleHistoryCollection = "pool_roles_history"

func RolesHistoryPermittedPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(UserCrewCollection, "user_id", "user_id", "crew")
	pipe.LookupUnwind(ProfileCollection, "user_id", "user_id", "profile")
	pipe.LookupUnwind(UserCollection, "user_id", "_id", "user")
	return
}

func RolesHistoryPermission(user *User, token *vcapool.AccessToken) (err error) {
	if user.NVM.Status != "confirmed" {
		return vcago.NewBadRequest(PoolRoleHistoryCollection, "nvm required", nil)
	}
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(PoolRoleHistoryCollection)
	}
	return
}

func RolesHistoryAdminPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(PoolRoleHistoryCollection)
	}
	return
}

func (i *RoleHistory) NewRole() (r *vmod.Role, err error) {
	switch i.Role {
	case "asp":
		return RoleASP(i.UserID), err
	case "finance":
		return RoleFinance(i.UserID), err
	case "operation":
		return RoleAction(i.UserID), err
	case "education":
		return RoleEducation(i.UserID), err
	case "network":
		return RoleNetwork(i.UserID), err
	case "socialmedia":
		return RoleSocialMedia(i.UserID), err
	case "awareness":
		return RoleAwareness(i.UserID), err
	case "other":
		return RoleOther(i.UserID), err
	default:
		return nil, vcago.NewValidationError("role not supported: " + i.Role)
	}
}
func (i *RoleHistoryCreate) NewRoleHistory() *RoleHistory {
	return &RoleHistory{
		ID:        uuid.NewString(),
		Role:      i.Role,
		UserID:    i.UserID,
		CrewID:    i.CrewID,
		Confirmed: i.Confirmed,
		StartDate: i.StartDate,
		EndDate:   i.EndDate,
		Modified:  vmod.NewModified(),
	}
}

func (i *RoleHistory) MatchUser() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.UserID)
	return filter.Bson()
}

func (i *RoleHistory) FilterRole() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("name", i.Role)
	filter.EqualString("user_id", i.UserID)
	return filter.Bson()
}

func (i *RoleHistory) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("role", i.Role)
	filter.EqualString("user_id", i.UserID)
	return filter.Bson()
}

func (i *RoleHistoryRequest) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("role", i.Role)
	filter.EqualBool("confirmed", strconv.FormatBool(i.Confirmed))
	filter.EqualString("crew_id", i.CrewID)
	filter.EqualString("user_id", i.UserID)
	return filter.Bson()
}

func (i *RoleHistoryDatabase) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("role", i.Role)
	filter.EqualString("user_id", i.UserID)
	filter.EqualInt("end_date", "0")
	return filter.Bson()
}

func (i *RoleHistoryBulkRequest) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew.crew_id", token.CrewID)
	} else {
		filter.EqualString("crew.crew_id", i.CrewID)
	}
	filter.ElemMatchList("pool_roles", "name", []string{"network", "education", "finance", "operation", "awareness", "socialmedia", "other"})
	return filter.Bson()
}

func (i *RoleHistoryRequest) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew_id", token.CrewID)
	} else {
		filter.EqualString("crew_id", i.CrewID)
	}
	filter.EqualBool("confirmed", strconv.FormatBool(i.Confirmed))
	return filter.Bson()
}
