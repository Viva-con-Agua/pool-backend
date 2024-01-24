package dao

import (
	"context"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func RoleInsert(ctx context.Context, i *models.RoleRequest, token *vcapool.AccessToken) (result *vmod.Role, err error) {
	filter := i.MatchUser()
	user := new(models.User)
	if err = UserCollection.AggregateOne(
		ctx,
		models.UserPipeline(false).Match(filter).Pipe,
		user,
	); err != nil {
		return
	}
	if err = models.RolesPermission(i.Role, user, token); err != nil {
		return
	}
	if result, err = i.NewRole(); err != nil {
		return
	}
	if err = PoolRoleCollection.InsertOne(ctx, result); err != nil {
		return
	}
	if err = PoolRoleHistoryCollection.InsertOne(ctx, models.NewRoleHistory(result, user)); err != nil {
		return
	}
	return
}

func RoleBulkUpdate(ctx context.Context, i *models.RoleBulkRequest, token *vcapool.AccessToken) (result *models.RoleBulkExport, userRolesMap map[string]*models.BulkUserRoles, err error) {
	if err = models.RolesBulkPermission(token); err != nil {
		return
	}

	userCrewRoles := new([]models.User)
	userCrewFilter := i.PermittedFilter(token)
	if err = UserCollection.Aggregate(ctx, models.UserPipelinePublic().Match(userCrewFilter).Pipe, userCrewRoles); err != nil {
		log.Print("Currently no roles set yet")
	}

	userRolesMap = make(map[string]*models.BulkUserRoles)
	result = new(models.RoleBulkExport)
	for _, role := range i.AddedRoles {
		filter := role.MatchUser()
		user := new(models.User)
		if err = UserCollection.AggregateOne(
			ctx,
			models.UserPipeline(false).Match(filter).Pipe,
			user,
		); err != nil {
			return
		}

		if err = models.RolesPermission(role.Role, user, token); err != nil {
			return
		}
		userRole := new(models.RoleDatabase)
		result.Users = append(result.Users, models.ExportRole{UserID: user.ID, Role: role.Role})

		if err = PoolRoleCollection.FindOne(ctx, role.Filter(), userRole); err != nil {

			createdRole := new(vmod.Role)
			if createdRole, err = role.NewRole(); err != nil {
				return
			}
			if err = PoolRoleCollection.InsertOne(ctx, createdRole); err != nil {
				return
			}
			if err = PoolRoleHistoryCollection.InsertOne(ctx, models.NewRoleRequestHistory(&role, user)); err != nil {
				return
			}
			if token.ID != role.UserID {
				if userRolesMap[role.UserID] == nil {
					userRolesMap[role.UserID] = &models.BulkUserRoles{}
				}
				userRolesMap[role.UserID].AddedRoles = append(userRolesMap[role.UserID].AddedRoles, createdRole.Label)
			}
		}

	}
	for _, role := range i.DeletedRoles {
		filter := role.MatchUser()
		user := new(models.User)
		if err = UserCollection.FindOne(
			ctx,
			filter,
			user,
		); err != nil {
			return
		}
		deleteRole := new(vmod.Role)
		if err = PoolRoleCollection.FindOne(
			ctx,
			role.Filter(),
			deleteRole,
		); err != nil {
			return
		}
		if err = models.RolesDeletePermission(deleteRole.Name, token); err != nil {
			return
		}
		if err = PoolRoleCollection.DeleteOne(ctx, role.Filter()); err != nil {
			return
		}
		history := new(models.RoleHistoryUpdate)
		if err = PoolRoleHistoryCollection.FindOne(
			ctx,
			role.FilterHistory(),
			&history,
		); err != nil {
			return
		}
		history.EndDate = time.Now().Unix()
		if err = PoolRoleHistoryCollection.UpdateOne(ctx, role.FilterHistory(), vmdb.UpdateSet(history), history); err != nil {
			return
		}
		if token.ID != role.UserID {
			if userRolesMap[role.UserID] == nil {
				userRolesMap[role.UserID] = &models.BulkUserRoles{}
			}
			userRolesMap[role.UserID].DeletedRoles = append(userRolesMap[role.UserID].DeletedRoles, deleteRole.Label)
		}
	}
	result.CrewID = i.CrewID
	return
}

func RoleBulkConfirm(ctx context.Context, i *[]models.RoleHistory, crew_id string, token *vcapool.AccessToken) (result *models.RoleBulkExport, userRolesMap map[string]*models.AspBulkUserRoles, err error) {
	if err = models.RolesAdminPermission(token); err != nil {
		return
	}

	userRolesMap = make(map[string]*models.AspBulkUserRoles)

	role_filter := bson.D{{Key: "crew.crew_id", Value: crew_id}, {Key: "pool_roles", Value: bson.D{{Key: "$exists", Value: true}, {Key: "$ne", Value: "[]"}}}}
	deleted_roles_users := new([]models.User)
	UserViewCollection.Find(ctx, role_filter, deleted_roles_users)
	deleted_roles := new([]vmod.Role)
	for _, user := range *deleted_roles_users {
		*deleted_roles = append(*deleted_roles, user.PoolRoles...)

		for _, role := range user.PoolRoles {
			if err = PoolRoleCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: role.ID}}); err != nil {
				return
			}
		}
	}

	result = new(models.RoleBulkExport)
	for _, role := range *i {
		filter := role.MatchUser()
		user := new(models.User)
		if err = UserCollection.AggregateOne(
			ctx,
			models.UserPipeline(false).Match(filter).Pipe,
			user,
		); err != nil {
			return
		}

		if err = models.RolesPermission(role.Role, user, token); err != nil {
			return
		}
		userRole := new(models.RoleDatabase)
		result.Users = append(result.Users, models.ExportRole{UserID: user.ID, Role: role.Role})

		if err = PoolRoleCollection.FindOne(ctx, role.FilterRole(), userRole); err != nil {

			createdRole := new(vmod.Role)
			if createdRole, err = role.NewRole(); err != nil {
				return
			}
			if err = PoolRoleCollection.InsertOne(ctx, createdRole); err != nil {
				return
			}

			if token.ID != role.UserID {
				if userRolesMap[role.UserID] == nil {
					userRolesMap[role.UserID] = &models.AspBulkUserRoles{}
				}
				if index := getIndex(createdRole, *deleted_roles); index >= 0 {
					userRolesMap[role.UserID].UnchangedRoles = append(userRolesMap[role.UserID].UnchangedRoles, createdRole.Label)
					*deleted_roles = (*deleted_roles)[:index+copy((*deleted_roles)[index:], (*deleted_roles)[index+1:])]
				} else {
					if userRolesMap[role.UserID] == nil {
						userRolesMap[role.UserID] = &models.AspBulkUserRoles{}
					}
					userRolesMap[role.UserID].AddedRoles = append(userRolesMap[role.UserID].AddedRoles, createdRole.Label)
				}
			}
		}
	}
	for _, role := range *deleted_roles {
		if token.ID != role.UserID {
			if userRolesMap[role.UserID] == nil {
				userRolesMap[role.UserID] = &models.AspBulkUserRoles{}
			}
			userRolesMap[role.UserID].DeletedRoles = append(userRolesMap[role.UserID].DeletedRoles, role.Label)
		}
	}
	result.CrewID = crew_id
	return
}
func getIndex(role *vmod.Role, data []vmod.Role) (index int) {
	for index, search := range data {
		if search.Name == role.Name && search.UserID == role.UserID {
			return index
		}
	}
	return -1
}
func RoleDelete(ctx context.Context, i *models.RoleRequest, token *vcapool.AccessToken) (result *vmod.Role, err error) {
	filter := i.MatchUser()
	user := new(models.User)
	if err = UserCollection.FindOne(
		ctx,
		filter,
		user,
	); err != nil {
		return
	}
	if err = PoolRoleCollection.FindOne(
		ctx,
		i.Filter(),
		&result,
	); err != nil {
		return
	}
	if err = models.RolesDeletePermission(result.Name, token); err != nil {
		return
	}
	if err = PoolRoleCollection.DeleteOne(ctx, i.Filter()); err != nil {
		return
	}
	history := new(models.RoleHistory)
	if err = PoolRoleHistoryCollection.FindOne(
		ctx,
		i.Filter(),
		&history,
	); err != nil {
		return
	}
	history.EndDate = time.Now().Unix()
	if err = PoolRoleHistoryCollection.UpdateOne(ctx, history.Filter(), vmdb.UpdateSet(history), history); err != nil {
		return
	}
	return
}

func RoleNotification(ctx context.Context, i map[string]*models.BulkUserRoles) (err error) {
	for index, role := range i {
		user := new(models.User)
		if err = UserCollection.FindOne(
			ctx,
			bson.D{{Key: "_id", Value: index}},
			user,
		); err != nil {
			return
		}
		mail := vcago.NewMailData(user.Email, "pool-backend", "role_update", "pool", user.Country)
		mail.AddUser(user.User())
		mail.AddContent(user.RoleContent(role))
		vcago.Nats.Publish("system.mail.job", mail)
	}
	return
}

func AspRoleNotification(ctx context.Context, i map[string]*models.AspBulkUserRoles) (err error) {
	for index, role := range i {
		user := new(models.User)
		if err = UserCollection.FindOne(
			ctx,
			bson.D{{Key: "_id", Value: index}},
			user,
		); err != nil {
			return
		}
		mail := vcago.NewMailData(user.Email, "pool-backend", "asp_role_update", "pool", user.Country)
		mail.AddUser(user.User())
		mail.AddContent(user.AspRoleContent(role))
		vcago.Nats.Publish("system.mail.job", mail)
	}
	return
}

func RoleAdminNotification(ctx context.Context, crewID *models.CrewParam) (err error) {
	crew := new(models.Crew)
	mail := vcago.NewMailData("netzwerk@vivaconagua.org", "pool-backend", "role_network", "pool", "de")
	mail.AddContent(models.RoleAdminContent(crew))
	vcago.Nats.Publish("system.mail.job", mail)
	return
}
