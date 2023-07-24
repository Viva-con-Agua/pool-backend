package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type RoleBulkRequest struct {
	CrewID       string        `json:"crew_id"`
	AddedRoles   []RoleRequest `json:"created"`
	DeletedRoles []RoleRequest `json:"removed"`
}

type RoleBulkExport struct {
	CrewID string       `bson:"crew_id" json:"crew_id"`
	Users  []ExportRole `json:"users"`
}
type ExportRole struct {
	UserID string `json:"uuid"`
	Role   string `json:"role"`
}

type RoleAdminRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RoleImport struct {
	Role    string `json:"role"`
	DropsID string `json:"drops_id"`
}

type RoleDatabase struct {
	ID     string `bson:"_id" json:"id"`
	Name   string `bson:"name" json:"name"`
	UserID string `bson:"user_id" json:"user_id"`
	Label  string `bson:"label" json:"label"`
	Root   string `bson:"root" json:"root"`
}

type BulkUserRoles struct {
	AddedRoles   []string `bson:"created" json:"created"`
	DeletedRoles []string `bson:"deleted" json:"deleted"`
}

var PoolRoleCollection = "pool_roles"

func (i *RoleRequest) NewRole() (r *vmod.Role, err error) {
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

var ASPRole = "asp;finance;operation;education;network;socialmedia;awareness"

func RolesPermission(role string, user *User, token *vcapool.AccessToken) (err error) {
	if user.NVM.Status != "confirmed" {
		return vcago.NewBadRequest(PoolRoleCollection, "nvm required", nil)
	}
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate(role)) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RolesBulkPermission(token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("finance;operation;education;network;socialmedia;awareness;other")) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RolesDeletePermission(role string, token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate(role)) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RoleASP(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "asp",
		Label:  "ASP",
		Root:   "employee;admin",
		UserID: userID,
	}
}

func RoleSupporter(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "supporter",
		Label:  "Supporter",
		Root:   "system",
		UserID: userID,
	}
}

func RoleFinance(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "finance",
		Label:  "Finanzen",
		Root:   "finance;employee;admin",
		UserID: userID,
	}
}

func RoleAction(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "operation",
		Label:  "Aktion",
		Root:   "operation;employee;admin",
		UserID: userID,
	}
}
func RoleEducation(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "education",
		Label:  "Bildung",
		Root:   "education;employee;admin",
		UserID: userID,
	}
}
func RoleNetwork(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "network",
		Label:  "Netzwerk",
		Root:   "network;employee;admin",
		UserID: userID,
	}
}
func RoleSocialMedia(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "socialmedia",
		Label:  "Social Media",
		Root:   "socialmedia;employee;admin",
		UserID: userID,
	}
}
func RoleAwareness(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "awareness",
		Label:  "Awareness",
		Root:   "awareness;employee;admin",
		UserID: userID,
	}
}

func RoleOther(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "other",
		Label:  "Other",
		Root:   "other;employee;admin",
		UserID: userID,
	}
}

func (i *RoleRequest) MatchUser() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.UserID)
	return filter.Bson()
}

func (i *RoleRequest) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("name", i.Role)
	filter.EqualString("user_id", i.UserID)
	return filter.Bson()
}

func (i *RoleBulkRequest) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew.crew_id", token.CrewID)
		filter.ElemMatchList("pool_roles", "name", token.PoolRoles)
	} else {
		filter.EqualString("crew.crew_id", i.CrewID)
		filter.ElemMatchList("pool_roles", "name", []string{"network", "education", "finance", "operation", "awareness", "socialmedia", "other"})
	}
	return filter.Bson()
}
