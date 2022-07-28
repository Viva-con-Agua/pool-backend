package models

import (
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

type RoleAdminRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (i *RoleRequest) New() (r *vmod.Role, err error) {
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
	default:
		return nil, vcago.NewValidationError("role not supported: " + i.Role)
	}
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

func (i *RoleRequest) MatchUser() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.UserID)
	return match.Bson()
}

func (i *RoleRequest) Filter() bson.D {
	return bson.D{{Key: "name", Value: i.Role}, {Key: "user_id", Value: i.UserID}}
}
