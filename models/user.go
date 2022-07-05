package models

import (
	"pool-user/dao"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserEmail struct {
		Email string `json:"email"`
	}
	UserDatabase struct {
		ID            string         `json:"id,omitempty" bson:"_id"`
		Email         string         `json:"email" bson:"email" validate:"required,email"`
		FirstName     string         `bson:"first_name" json:"first_name" validate:"required"`
		LastName      string         `bson:"last_name" json:"last_name" validate:"required"`
		FullName      string         `bson:"full_name" json:"full_name"`
		DisplayName   string         `bson:"display_name" json:"display_name"`
		Roles         vcago.RoleList `json:"system_roles" bson:"system_roles"`
		Country       string         `bson:"country" json:"country"`
		PrivacyPolicy bool           `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed     bool           `bson:"confirmed" json:"confirmed"`
		DropsID       string         `bson:"drops_id" json:"drops_id"`
		LastUpdate    string         `bson:"last_update" json:"last_update"`
		Modified      vcago.Modified `json:"modified" bson:"modified"`
	}
	UserUpdate struct {
		ID            string         `json:"id,omitempty" bson:"_id"`
		Email         string         `json:"email" bson:"email" validate:"required,email"`
		FirstName     string         `bson:"first_name" json:"first_name" validate:"required"`
		LastName      string         `bson:"last_name" json:"last_name" validate:"required"`
		FullName      string         `bson:"full_name" json:"full_name"`
		DisplayName   string         `bson:"display_name" json:"display_name"`
		Roles         vcago.RoleList `json:"system_roles" bson:"system_roles"`
		Country       string         `bson:"country" json:"country"`
		PrivacyPolicy bool           `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed     bool           `bson:"confirmed" json:"confirmed"`
		DropsID       string         `bson:"drops_id" json:"drops_id"`
		LastUpdate    string         `bson:"last_update" json:"last_update"`
	}
	User struct {
		ID            string         `json:"id,omitempty" bson:"_id"`
		Email         string         `json:"email" bson:"email" validate:"required,email"`
		FirstName     string         `bson:"first_name" json:"first_name" validate:"required"`
		LastName      string         `bson:"last_name" json:"last_name" validate:"required"`
		FullName      string         `bson:"full_name" json:"full_name"`
		DisplayName   string         `bson:"display_name" json:"display_name"`
		Roles         vcago.RoleList `json:"system_roles" bson:"system_roles"`
		Country       string         `bson:"country" json:"country"`
		PrivacyPolicy bool           `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed     bool           `bson:"confirmed" json:"confirmed"`
		LastUpdate    string         `bson:"last_update" json:"last_update"`
		//extends the vcago.User
		DropsID   string         `bson:"drops_id" json:"drops_id"`
		Profile   Profile        `json:"profile" bson:"profile,truncate"`
		Crew      UserCrew       `json:"crew" bson:"crew,omitempty"`
		Avatar    Avatar         `bson:"avatar,omitempty" json:"avatar"`
		Address   Address        `json:"address" bson:"address,omitempty"`
		PoolRoles vcago.RoleList `json:"pool_roles" bson:"pool_roles,omitempty"`
		Active    Active         `json:"active" bson:"active,omitempty"`
		NVM       NVM            `json:"nvm" bson:"nvm,omitempty"`
		Modified  vcago.Modified `json:"modified" bson:"modified"`
	}
	UserParam struct {
		ID string `param:"id"`
	}
	UserQuery struct {
		ID            []string `query:"id" qs:"id"`
		FirstName     string   `query:"first_name" qs:"first_name"`
		LastName      string   `query:"last_name" qs:"last_name"`
		FullName      string   `query:"full_name" qs:"full_name"`
		DisplayName   string   `query:"display_name" qs:"display_name"`
		ActiveState   []string `query:"active_state" qs:"active_state"`
		SystemRoles   []string `query:"system_roles" qs:"system_roles"`
		PoolRoles     []string `query:"pool_roles" qs:"pool_roles"`
		PrivacyPolicy string   `query:"privacy_policy" qs:"privacy_policy"`
		NVMState      []string `query:"nvm_state" qs:"nvm_state"`
		CrewID        string   `query:"crew_id" qs:"crew_id"`
		Country       string   `query:"country" qs:"country"`
		Confirmed     string   `query:"confirmed" qs:"confirmed"`
		UpdatedTo     string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom   string   `query:"updated_from" qs:"updated_from"`
		CreatedTo     string   `query:"created_to" qs:"created_to"`
		CreatedFrom   string   `query:"created_from" qs:"created_from"`
	}
)

func NewUserDatabase(user *vcago.User) *UserDatabase {
	return &UserDatabase{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FullName,
		DisplayName:   user.DisplayName,
		Roles:         user.Roles,
		Country:       user.Country,
		PrivacyPolicy: user.PrivacyPolicy,
		Confirmed:     user.Confirmd,
		LastUpdate:    user.LastUpdate,
		Modified:      vcago.NewModified(),
	}
}

func NewUserUpdate(user *vcago.User) *UserUpdate {
	return &UserUpdate{
		ID:            user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FullName,
		DisplayName:   user.DisplayName,
		Roles:         user.Roles,
		Country:       user.Country,
		PrivacyPolicy: user.PrivacyPolicy,
		Confirmed:     user.Confirmd,
		LastUpdate:    user.LastUpdate,
	}
}

func UserPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(dao.AddressesCollection.Name, "_id", "user_id", "address")
	pipe.LookupUnwind(dao.ProfilesCollection.Name, "_id", "user_id", "profile")
	pipe.LookupUnwind(dao.UserCrewCollection.Name, "_id", "user_id", "crew")
	pipe.LookupUnwind(dao.ActiveCollection.Name, "_id", "user_id", "active")
	pipe.LookupUnwind(dao.NVMCollection.Name, "_id", "user_id", "nvm")
	pipe.Lookup(dao.PoolRoleCollection.Name, "_id", "user_id", "pool_roles")
	pipe.LookupUnwind(dao.AvatarCollection.Name, "_id", "user_id", "avatar")
	return
}

func UserMatch(userID string) (match *vmdb.Match) {
	match = vmdb.NewMatch()
	match.EqualString("_id", userID)
	return
}

func UserMatchEmail(email string) (match *vmdb.Match) {
	match = vmdb.NewMatch()
	match.EqualString("email", email)
	return
}

func (i *UserUpdate) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *User) AuthToken() (r *vcapool.AuthToken, err error) {
	accessToken := &vcapool.AccessToken{
		ID:            i.ID,
		Email:         i.Email,
		FirstName:     i.FirstName,
		LastName:      i.LastName,
		FullName:      i.FullName,
		DisplayName:   i.DisplayName,
		Roles:         *i.Roles.Cookie(),
		Country:       i.Country,
		PrivacyPolicy: i.PrivacyPolicy,
		Confirmd:      i.Confirmed,
		LastUpdate:    i.LastUpdate,
		Phone:         i.Profile.Phone,
		Gender:        i.Profile.Gender,
		Birthdate:     i.Profile.Birthdate,
		CrewName:      i.Crew.Name,
		CrewID:        i.Crew.CrewID,
		CrewEmail:     i.Crew.Email,
		AddressID:     i.Address.ID,
		PoolRoles:     *i.PoolRoles.Cookie(),
		ActiveState:   i.Active.Status,
		NVMState:      i.NVM.Status,
		AvatarID:      i.Avatar.ID,
		Modified:      i.Modified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	refreshToken := &vcapool.RefreshToken{
		UserID: i.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	r = new(vcapool.AuthToken)
	if r.AccessToken, err = accessToken.SignedString(vcapool.JWTSecret); err != nil {
		return
	}
	r.RefreshToken, err = refreshToken.SignedString(vcapool.JWTSecret)
	return
}

func (i UserParam) Pipeline() mongo.Pipeline {
	match := vmdb.NewMatch()
	match.EqualString("_id", i.ID)
	return UserPipeline().Match(match).Pipe
}

func (i *UserParam) Filter(token *vcapool.AccessToken) bson.D {
	if token.Roles.Validate("employee;admin") {
		return bson.D{{Key: "_id", Value: i.ID}}
	} else {
		return bson.D{{Key: "_id", Value: token.ID}}
	}
}

func (i *UserQuery) Match() *vmdb.Match {
	match := vmdb.NewMatch()
	match.LikeString("first_name", i.FirstName)
	match.LikeString("last_name", i.LastName)
	match.LikeString("full_name", i.FullName)
	match.LikeString("display_name", i.DisplayName)
	match.EqualString("crew.crew_id", i.CrewID)
	match.ElemMatchList("system_roles", "name", i.SystemRoles)
	match.ElemMatchList("pool_roles", "name", i.PoolRoles)
	match.EqualBool("privacy_policy", i.PrivacyPolicy)
	match.EqualStringList("active.status", i.ActiveState)
	match.EqualStringList("nvm.status", i.NVMState)
	match.EqualString("crew.crew_id", i.CrewID)
	match.EqualString("country", i.Country)
	match.EqualBool("confirmed", i.Confirmed)
	match.GteInt64("modified.updated", i.UpdatedFrom)
	match.GteInt64("modified.created", i.CreatedFrom)
	match.LteInt64("modified.updated", i.UpdatedTo)
	match.LteInt64("modified.created", i.CreatedTo)
	return match
}

func (i *UserQuery) Pipeline() mongo.Pipeline {
	match := i.Match()
	return UserPipeline().Match(match).Pipe
}