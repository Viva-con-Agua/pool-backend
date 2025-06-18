package models

import (
	"slices"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
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

type AspBulkUserRoles struct {
	AddedRoles     []string `bson:"created" json:"created"`
	DeletedRoles   []string `bson:"deleted" json:"deleted"`
	UnchangedRoles []string `bson:"unchanged" json:"unchanged"`
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
func (i *RoleRequest) NewRoleHistory(user *User) (r *RoleHistoryDatabase) {
	return &RoleHistoryDatabase{
		ID:        uuid.NewString(),
		Role:      i.Role,
		UserID:    user.ID,
		CrewID:    user.Crew.CrewID,
		StartDate: time.Now().Unix(),
		Confirmed: false,
		Modified:  vmod.NewModified(),
	}
}

func NewRoleHistory(i *vmod.Role, user *User) (r *RoleHistoryDatabase) {
	return &RoleHistoryDatabase{
		ID:        uuid.NewString(),
		Role:      i.Name,
		UserID:    user.ID,
		CrewID:    user.Crew.CrewID,
		StartDate: time.Now().Unix(),
		Confirmed: true,
		Modified:  vmod.NewModified(),
	}
}
func NewRoleRequestHistory(i *RoleRequest, user *User) (r *RoleHistoryDatabase) {
	return &RoleHistoryDatabase{
		ID:        uuid.NewString(),
		Role:      i.Role,
		UserID:    user.ID,
		CrewID:    user.Crew.CrewID,
		StartDate: time.Now().Unix(),
		Confirmed: true,
		Modified:  vmod.NewModified(),
	}
}

var ASPRole = "other;asp;finance;operation;education;network;socialmedia;awareness"
var ASPEventRole = "network;operation;education"

func RolesPermission(role string, user *User, token *AccessToken, options []string) (err error) {
	if user.NVM.Status != "confirmed" && slices.Contains(options, "nvm") {
		return vcago.NewBadRequest(PoolRoleCollection, "nvm required", nil)
	}
	if user.Active.Status != "confirmed" {
		return vcago.NewBadRequest(PoolRoleCollection, "active required", nil)
	}
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(role)) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RolesBulkPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RolesDeletePermission(role string, token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(role)) {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RolesAdminPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(PoolRoleCollection)
	}
	return
}

func RoleASP(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "asp",
		Root:   "admin;employee;pool_employee",
		UserID: userID,
	}
}

func RoleSupporter(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "supporter",
		Root:   "system",
		UserID: userID,
	}
}

func RoleFinance(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "finance",
		Root:   "finance;admin;employee;pool_employee",
		UserID: userID,
	}
}

func RoleAction(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "operation",
		Root:   "operation;admin;employee;pool_employee",
		UserID: userID,
	}
}
func RoleEducation(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "education",
		Root:   "education;admin;employee;pool_employee",
		UserID: userID,
	}
}
func RoleNetwork(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "network",
		Root:   "network;admin;employee;pool_employee",
		UserID: userID,
	}
}
func RoleSocialMedia(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "socialmedia",
		Root:   "socialmedia;admin;employee;pool_employee",
		UserID: userID,
	}
}
func RoleAwareness(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "awareness",
		Root:   "awareness;admin;employee;pool_employee",
		UserID: userID,
	}
}

func RoleOther(userID string) *vmod.Role {
	return &vmod.Role{
		ID:     uuid.NewString(),
		Name:   "other",
		Root:   "other;admin;employee;pool_employee",
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
func (i *RoleRequest) FilterHistory() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("role", i.Role)
	filter.EqualString("user_id", i.UserID)
	return filter.Bson()
}

func (i *RoleBulkRequest) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew.crew_id", token.CrewID)
		filter.ElemMatchList("pool_roles", "name", token.PoolRoles)
	} else {
		filter.EqualString("crew.crew_id", i.CrewID)
		filter.ElemMatchList("pool_roles", "name", []string{"network", "education", "finance", "operation", "awareness", "socialmedia", "other", "asp"})
	}
	return filter.Bson()
}
