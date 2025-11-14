package models

import (
	"strings"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserEmail struct {
		Email string `json:"email"`
	}
	UserDatabase struct {
		ID             string        `json:"id,omitempty" bson:"_id"`
		Email          string        `json:"email" bson:"email" validate:"required,email"`
		FirstName      string        `bson:"first_name" json:"first_name" validate:"required"`
		LastName       string        `bson:"last_name" json:"last_name" validate:"required"`
		FullName       string        `bson:"full_name" json:"full_name"`
		DisplayName    string        `bson:"display_name" json:"display_name"`
		Roles          vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Country        string        `bson:"country" json:"country"`
		PrivacyPolicy  bool          `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed      bool          `bson:"confirmed" json:"confirmed"`
		DropsID        string        `bson:"drops_id" json:"drops_id"`
		LastUpdate     string        `bson:"last_update" json:"last_update"`
		OrganisationID string        `bson:"organisation_id" json:"organisation_id"`
		MailboxID      string        `bson:"mailbox_id" json:"mailbox_id"`
		LastLoginDate  int64         `bson:"last_login_date" json:"last_login_date"`
		NVM            NVM           `json:"nvm" bson:"nvm,omitempty"`
		Active         Active        `json:"active" bson:"active,omitempty"`
		Modified       vmod.Modified `json:"modified" bson:"modified"`
	}
	UserUpdate struct {
		ID            string        `json:"id,omitempty" bson:"_id"`
		Email         string        `json:"email" bson:"email" validate:"required,email"`
		FirstName     string        `bson:"first_name" json:"first_name" validate:"required"`
		LastName      string        `bson:"last_name" json:"last_name" validate:"required"`
		FullName      string        `bson:"full_name" json:"full_name"`
		DisplayName   string        `bson:"display_name" json:"display_name"`
		Roles         vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Country       string        `bson:"country" json:"country"`
		PrivacyPolicy bool          `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed     bool          `bson:"confirmed" json:"confirmed"`
		DropsID       string        `bson:"drops_id" json:"drops_id"`
		LastLoginDate int64         `bson:"last_login_date" json:"last_login_date"`
		LastUpdate    string        `bson:"last_update" json:"last_update"`
	}
	UserOrganisationUpdate struct {
		ID             string `json:"id,omitempty" bson:"_id"`
		OrganisationID string `bson:"organisation_id" json:"organisation_id"`
	}
	User struct {
		ID            string        `json:"id,omitempty" bson:"_id"`
		Email         string        `json:"email" bson:"email" `
		FirstName     string        `bson:"first_name" json:"first_name" `
		LastName      string        `bson:"last_name" json:"last_name" `
		FullName      string        `bson:"full_name" json:"full_name"`
		DisplayName   string        `bson:"display_name" json:"display_name"`
		Roles         vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Country       string        `bson:"country" json:"country"`
		PrivacyPolicy bool          `bson:"privacy_policy" json:"privacy_policy"`
		Confirmed     bool          `bson:"confirmed" json:"confirmed"`
		LastUpdate    string        `bson:"last_update" json:"last_update"`
		MailboxID     string        `bson:"mailbox_id" json:"mailbox_id"`
		LastLoginDate int64         `bson:"last_login_date" json:"last_login_date"`
		//extends the vcago.User
		DropsID        string        `bson:"drops_id" json:"drops_id"`
		Profile        Profile       `json:"profile" bson:"profile,truncate"`
		Crew           UserCrew      `json:"crew" bson:"crew,omitempty"`
		OrganisationID string        `bson:"organisation_id" json:"organisation_id"`
		Organisation   string        `bson:"organisation" json:"organisation,omitempty"`
		Avatar         Avatar        `bson:"avatar,omitempty" json:"avatar"`
		Address        Address       `json:"address" bson:"address,omitempty"`
		AddressID      string        `json:"address_id" bson:"address_id"`
		PoolRoles      vmod.RoleList `json:"pool_roles" bson:"pool_roles,omitempty"`
		Active         Active        `json:"active" bson:"active,omitempty"`
		NVM            NVM           `json:"nvm" bson:"nvm,omitempty"`
		Newsletter     []Newsletter  `json:"newsletter" bson:"newsletter"`
		Modified       vmod.Modified `json:"modified" bson:"modified"`
	}
	ListUser struct {
		ID          string        `json:"id,omitempty" bson:"_id"`
		Email       string        `json:"email" bson:"email" `
		FirstName   string        `bson:"first_name" json:"first_name" `
		LastName    string        `bson:"last_name" json:"last_name" `
		FullName    string        `bson:"full_name" json:"full_name"`
		DisplayName string        `bson:"display_name" json:"display_name"`
		Roles       vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Country     string        `bson:"country" json:"country"`
		Confirmed   bool          `bson:"confirmed" json:"confirmed"`
		//extends the vcago.User
		Profile    ProfileMinimal  `json:"profile" bson:"profile,truncate"`
		Crew       UserCrewMinimal `json:"crew" bson:"crew,omitempty"`
		Avatar     Avatar          `bson:"avatar,omitempty" json:"avatar"`
		PoolRoles  vmod.RoleList   `json:"pool_roles" bson:"pool_roles,omitempty"`
		Active     Active          `json:"active" bson:"active,omitempty"`
		NVM        NVM             `json:"nvm" bson:"nvm,omitempty"`
		Newsletter []Newsletter    `json:"newsletter" bson:"newsletter"`
		Modified   vmod.Modified   `json:"modified" bson:"modified"`
	}
	UserParticipant struct {
		ID          string   `json:"id,omitempty" bson:"_id"`
		Email       string   `json:"email" bson:"email" `
		FirstName   string   `bson:"first_name" json:"first_name" `
		LastName    string   `bson:"last_name" json:"last_name" `
		FullName    string   `bson:"full_name" json:"full_name"`
		DisplayName string   `bson:"display_name" json:"display_name"`
		Country     string   `bson:"country" json:"country"`
		Profile     Profile  `json:"profile" bson:"profile,truncate"`
		Crew        CrewName `json:"crew" bson:"crew,omitempty"`
		Avatar      Avatar   `bson:"avatar,omitempty" json:"avatar"`
		Active      Active   `json:"active" bson:"active,omitempty"`
	}
	UserPublic struct {
		ID             string        `json:"id,omitempty" bson:"_id"`
		FirstName      string        `bson:"first_name" json:"first_name" `
		LastName       string        `bson:"last_name" json:"last_name" `
		FullName       string        `bson:"full_name" json:"full_name"`
		DisplayName    string        `bson:"display_name" json:"display_name"`
		Roles          vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Country        string        `bson:"country" json:"country"`
		Confirmed      bool          `bson:"confirmed" json:"confirmed"`
		LastUpdate     string        `bson:"last_update" json:"last_update"`
		OrganisationID string        `bson:"organisation_id" json:"organisation_id"`
		Organisation   string        `bson:"organisation" json:"organisation"`
		//extends the vcago.User
		DropsID   string        `bson:"drops_id" json:"drops_id"`
		Profile   Profile       `json:"profile" bson:"profile,truncate"`
		Crew      UserCrew      `json:"crew" bson:"crew,omitempty"`
		Avatar    Avatar        `bson:"avatar,omitempty" json:"avatar"`
		PoolRoles vmod.RoleList `json:"pool_roles" bson:"pool_roles,omitempty"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
	UserMinimal struct {
		ID          string `json:"id,omitempty" bson:"_id"`
		FirstName   string `bson:"first_name" json:"first_name" `
		LastName    string `bson:"last_name" json:"last_name" `
		FullName    string `bson:"full_name" json:"full_name"`
		DisplayName string `bson:"display_name" json:"display_name"`
	}
	UserBasic struct {
		ID        string `json:"id,omitempty" bson:"_id"`
		FirstName string `bson:"first_name" json:"first_name" `
		LastName  string `bson:"last_name" json:"last_name" `
		FullName  string `bson:"full_name" json:"full_name"`
		//Profile     ProfileMinimal `bson:"profile" json:"profile"`
		DisplayName string        `bson:"display_name" json:"display_name"`
		Roles       vmod.RoleList `json:"system_roles" bson:"system_roles"`
		Avatar      Avatar        `bson:"avatar,omitempty" json:"avatar"`
		PoolRoles   vmod.RoleList `json:"pool_roles" bson:"pool_roles,omitempty"`
		NVM         NVM           `json:"nvm" bson:"nvm,omitempty"`
	}
	UserContact struct {
		ID        string         `json:"id,omitempty" bson:"_id"`
		Email     string         `json:"email" bson:"email" `
		FirstName string         `bson:"first_name" json:"first_name" `
		LastName  string         `bson:"last_name" json:"last_name" `
		FullName  string         `bson:"full_name" json:"full_name"`
		Profile   ProfileMinimal `bson:"profile" json:"profile"`
	}
	UserParam struct {
		ID string `param:"id"`
	}
	UserQuery struct {
		ID             []string `query:"id"`
		Email          string   `query:"email"`
		FirstName      string   `query:"first_name"`
		LastName       string   `query:"last_name"`
		FullName       string   `query:"full_name"`
		DisplayName    string   `query:"display_name" qs:"display_name"`
		OrganisationID []string `query:"organisation_id" qs:"organisation_id"`
		Abbreviation   []string `query:"organisation_abbreviation" qs:"organisation_abbreviation"`
		ActiveState    []string `query:"active_state" qs:"active_state"`
		NVMState       []string `query:"nvm_state" qs:"nvm_state"`
		SystemRoles    []string `query:"system_roles"`
		PoolRoles      []string `query:"pool_roles"`
		CrewID         string   `query:"crew_id" qs:"crew_id"`
		FullCount      string   `query:"full_count"`
		vmdb.Query
	}
)

var UserCollection = "users"
var UserView = "users_view"

func (i *User) User() *vmod.User {
	return &vmod.User{
		ID:            i.ID,
		Email:         i.Email,
		FirstName:     i.FirstName,
		LastName:      i.LastName,
		FullName:      i.FirstName + " " + i.LastName,
		RealName:      i.FirstName + " " + i.LastName,
		DisplayName:   i.DisplayName,
		Roles:         i.Roles,
		Country:       i.Country,
		PrivacyPolicy: i.PrivacyPolicy,
		Confirmd:      i.Confirmed,
		LastUpdate:    i.LastUpdate,
	}
}

func NewUserDatabase(user *vmod.User) *UserDatabase {
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
		NVM:           *NVMClean(),
		Active:        *ActiveClean(),
		Modified:      vmod.NewModified(),
	}
}

func NewUserUpdate(user *vmod.User) *UserUpdate {
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

func UserPipeline(user bool) (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	if user {
		pipe.LookupUnwind(AddressesCollection, "_id", "user_id", "address")
	} else {
		pipe.LookupUnwind(AddressesCollection, "_id", "user_id", "address_data")
		pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "address_id", Value: "$address_data._id"}}}})
	}
	//pipe.LookupUnwind(ProfileCollection, "_id", "user_id", "profile")
	//pipe.LookupUnwind(UserCrewCollection, "_id", "user_id", "crew")
	//pipe.LookupUnwind(ActiveCollection, "_id", "user_id", "active")
	//pipe.LookupUnwind(NVMCollection, "_id", "user_id", "nvm")
	pipe.Lookup(PoolRoleCollection, "_id", "user_id", "pool_roles")
	pipe.Lookup(NewsletterCollection, "_id", "user_id", "newsletter")
	//pipe.LookupUnwind(AvatarCollection, "_id", "user_id", "avatar")

	return
}

func SortedUserPermittedPipeline(token *AccessToken) (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	//pipe.LookupUnwind(ProfileCollection, "_id", "user_id", "profile")
	//pipe.LookupUnwind(UserCrewCollection, "_id", "user_id", "crew")
	//pipe.LookupUnwind(ActiveCollection, "_id", "user_id", "active")
	//	pipe.LookupUnwind(NVMCollection, "_id", "user_id", "nvm")
	pipe.Lookup(PoolRoleCollection, "_id", "user_id", "pool_roles")
	return
}

func UserPermittedPipeline(token *AccessToken) (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	if token.Roles.Validate("admin") {
		pipe.LookupUnwind(AddressesCollection, "_id", "user_id", "address")
		pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "address_id", Value: "$address._id"}}}})
	} else {
		pipe.LookupUnwind(AddressesCollection, "_id", "user_id", "address_data")
		pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "address_id", Value: "$address_data._id"}}}})
	}
	//pipe.LookupUnwind(ProfileCollection, "_id", "user_id", "profile")
	//pipe.LookupUnwind(UserCrewCollection, "_id", "user_id", "crew")
	//pipe.LookupUnwind(ActiveCollection, "_id", "user_id", "active")
	//	pipe.LookupUnwind(NVMCollection, "_id", "user_id", "nvm")
	pipe.Lookup(PoolRoleCollection, "_id", "user_id", "pool_roles")
	if token.Roles.Validate("admin;employee;pool_employee") {
		pipe.Lookup(NewsletterCollection, "_id", "user_id", "newsletter")
	}
	//pipe.LookupUnwind(AvatarCollection, "_id", "user_id", "avatar")

	return
}

func UserPipelinePublic() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	//pipe.LookupUnwind(UserCrewCollection, "_id", "user_id", "crew")
	pipe.Lookup(PoolRoleCollection, "_id", "user_id", "pool_roles")
	//	pipe.LookupUnwind(NVMCollection, "_id", "user_id", "nvm")
	//pipe.LookupUnwind(AvatarCollection, "_id", "user_id", "avatar")
	return
}

func UserMatch(userID string) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", userID)
	return filter.Bson()
}

func UserMatchEmail(email string) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("email", email)
	return filter.Bson()
}

func UserCountPipeline(filter bson.D) *vmdb.Pipeline {
	pipe := UserPipeline(false)
	pipe.Match(filter)
	pipe.Append(bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil}, {Key: "list_size", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
		}},
	})
	pipe.Append(bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}}}})
	return pipe
}
func (i *User) AuthToken() (r *vcago.AuthToken, err error) {
	accessToken := &AccessToken{
		ID:             i.ID,
		Email:          i.Email,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
		FullName:       i.FullName,
		DisplayName:    i.DisplayName,
		Roles:          *i.Roles.Cookie(),
		Country:        i.Country,
		PrivacyPolicy:  i.PrivacyPolicy,
		Confirmd:       i.Confirmed,
		LastUpdate:     i.LastUpdate,
		Phone:          i.Profile.Phone,
		Gender:         i.Profile.Gender,
		Birthdate:      i.Profile.Birthdate,
		CrewName:       i.Crew.Name,
		CrewID:         i.Crew.CrewID,
		CrewEmail:      i.Crew.Email,
		OrganisationID: i.OrganisationID,
		AddressID:      i.Address.ID,
		PoolRoles:      *i.PoolRoles.Cookie(),
		ActiveState:    i.Active.Status,
		NVMState:       i.NVM.Status,
		AvatarID:       i.Avatar.ID,
		MailboxID:      i.MailboxID,
		Modified:       i.Modified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	refreshToken := &vcago.RefreshToken{
		UserID: i.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	return vcago.NewAuthToken(accessToken, refreshToken)
}

func (i *ProfileUpdate) ToUserUpdate() *ProfileUpdate {
	return &ProfileUpdate{
		ID:         i.ID,
		Mattermost: i.Mattermost,
		Phone:      i.Phone,
		Gender:     i.Gender,
		Birthdate:  i.Birthdate,
	}
}

func UsersPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(UserCollection)
	}
	return
}

func (i *UserParam) UsersDeletePermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || i.ID == token.ID) {
		return vcago.NewPermissionDenied(UserCollection)
	}
	return
}

func UsersDetailsPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin;employee;pool_employee") {
		return vcago.NewPermissionDenied(UserCollection)
	}
	return
}

func UsersEditPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(UserCollection)
	}
	return
}

func (i *UserQuery) CrewUsersPermission(token *AccessToken) (err error) {
	if i.CrewID != "" && i.CrewID != token.CrewID {
		return vcago.NewPermissionDenied(UserCollection)
	}
	return
}

func (i *UserParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i UserQuery) Sort() bson.D {
	sort := vmdb.NewSort()
	sort.Add(i.SortField, i.SortDirection)
	return sort.Bson()
}

func (i *UserQuery) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualBool("confirmed", "true")
	filter.LikeString("email", i.Email)
	filter.LikeString("first_name", i.FirstName)
	filter.LikeString("last_name", i.LastName)
	filter.LikeString("full_name", i.FullName)
	filter.LikeString("display_name", i.DisplayName)
	filter.ElemMatchList("system_roles", "name", i.SystemRoles)
	filter.ElemMatchList("pool_roles", "name", i.PoolRoles)
	filter.EqualStringList("active.status", i.ActiveState)
	filter.EqualStringList("nvm.status", i.NVMState)
	if token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualStringList("crew.organisation_id", i.OrganisationID)
		filter.EqualString("crew.crew_id", i.CrewID)
	} else {
		filter.EqualString("crew.crew_id", token.CrewID)
	}
	filter.SearchString([]string{"_id", "first_name", "last_name", "email"}, i.Search)
	return filter.Bson()
}

func (i *UserQuery) PermittedUserFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.ElemMatchList("pool_roles", "name", []string{"network", "education", "finance", "operation", "awareness", "socialmedia", "other"})
	filter.EqualBool("confirmed", "true")
	filter.LikeString("first_name", i.FirstName)
	filter.LikeString("last_name", i.LastName)
	filter.LikeString("full_name", i.FullName)
	filter.LikeString("display_name", i.DisplayName)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew.crew_id", token.CrewID)
	} else {
		filter.EqualString("crew.crew_id", i.CrewID)
	}
	filter.ElemMatchList("system_roles", "name", i.SystemRoles)
	filter.EqualStringList("active.status", i.ActiveState)
	filter.EqualStringList("nvm.status", i.NVMState)
	filter.SearchString([]string{"_id", "first_name", "last_name", "full_name", "email"}, i.Search)
	return filter.Bson()
}

func (i *UserQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualBool("confirmed", "true")
	filter.LikeString("first_name", i.FirstName)
	filter.LikeString("last_name", i.LastName)
	filter.LikeString("full_name", i.FullName)
	filter.LikeString("display_name", i.DisplayName)
	filter.ElemMatchList("system_roles", "name", i.SystemRoles)
	filter.ElemMatchList("pool_roles", "name", i.PoolRoles)
	filter.EqualStringList("active.status", i.ActiveState)
	filter.EqualStringList("nvm.status", i.NVMState)
	filter.EqualString("crew.crew_id", i.CrewID)
	filter.SearchString([]string{"_id", "first_name", "last_name", "email"}, i.Search)
	return filter.Bson()
}

func (i *User) RoleContent(roles *BulkUserRoles) *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["AddedRoles"] = strings.Join(roles.AddedRoles, ", ")
	content.Fields["DeletedRoles"] = strings.Join(roles.DeletedRoles, ", ")
	return content
}

func (i *User) AspRoleContent(roles *AspBulkUserRoles) *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["AddedRoles"] = strings.Join(roles.AddedRoles, ", ")
	content.Fields["DeletedRoles"] = strings.Join(roles.DeletedRoles, ", ")
	content.Fields["UnchangedRoles"] = strings.Join(roles.UnchangedRoles, ", ")
	return content
}

func RoleAdminContent(crew *Crew) *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["Crew"] = crew
	return content
}
