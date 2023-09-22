package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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

func RoleAdminNotification(ctx context.Context, crewID *models.CrewParam) (err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, crewID.Match(), crew); err != nil {
	}
	mail := vcago.NewMailData("netzwerk@vivaconagua.org", "pool-backend", "role_network", "pool", "de")
	mail.AddContent(models.RoleAdminContent(crew))
	vcago.Nats.Publish("system.mail.job", mail)
	return
}
