package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func RoleHistoryInsert(ctx context.Context, i *models.RoleHistoryCreate, token *vcapool.AccessToken) (result *models.RoleHistory, err error) {
	if err = models.RolesHistoryAdminPermission(token); err != nil {
		return
	}
	if result = i.NewRoleHistory(); err != nil {
		return
	}
	if err = PoolRoleHistoryCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func RoleHistoryBulkInsert(ctx context.Context, i *models.RoleHistoryBulkRequest, token *vcapool.AccessToken) (result *models.RoleBulkExport, err error) {
	if err = models.RolesBulkPermission(token); err != nil {
		return
	}

	if token.Roles.Validate("admin;employee") {
		if _, err = RoleHistoryDelete(ctx, &models.RoleHistoryRequest{CrewID: i.CrewID, Confirmed: false}, token); err != nil {
			return
		}
	}
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

		if err = models.RolesHistoryPermission(user, token); err != nil {
			return
		}
		userRoleHistory := new(models.RoleHistoryDatabase)
		result.Users = append(result.Users, models.ExportRole{UserID: user.ID, Role: role.Role})

		if err = PoolRoleHistoryCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: user.ID}, {Key: "role", Value: role.Role}, {Key: "end_date", Value: int64(0)}, {Key: "crew_id", Value: i.CrewID}}, userRoleHistory); err != nil {
			if err = PoolRoleHistoryCollection.InsertOne(ctx, role.NewRoleHistory(user)); err != nil {
				return
			}
		}

	}
	result.CrewID = i.CrewID
	return
}

func RoleHistoryGet(ctx context.Context, i *models.RoleHistoryRequest, token *vcapool.AccessToken) (result *[]models.RoleHistory, list_size int64, err error) {
	result = new([]models.RoleHistory)
	pipeline := models.RolesHistoryPermittedPipeline()
	if err = PoolRoleHistoryCollection.Aggregate(
		ctx,
		pipeline.Match(i.PermittedFilter(token)).Pipe,
		result,
	); err != nil {
		return
	}
	list_size = int64(len(*result))
	return
}

func RoleHistoryConfirm(ctx context.Context, i *models.RoleHistoryRequest, token *vcapool.AccessToken) (result *[]models.RoleHistory, err error) {
	if err = models.RolesHistoryAdminPermission(token); err != nil {
		return
	}
	i.Confirmed = false
	result = new([]models.RoleHistory)
	if err = PoolRoleHistoryCollection.Find(ctx, i.PermittedFilter(token), result); err != nil {
		return
	}
	if err = PoolRoleHistoryCollection.UpdateMany(ctx, i.PermittedFilter(token), vmdb.UpdateSet(bson.D{{Key: "confirmed", Value: true}})); err != nil {
		return
	}
	return
}

func RoleHistoryDelete(ctx context.Context, i *models.RoleHistoryRequest, token *vcapool.AccessToken) (result *models.RoleHistory, err error) {
	if err = models.RolesHistoryAdminPermission(token); err != nil {
		return
	}
	if err = PoolRoleHistoryCollection.FindOne(
		ctx,
		i.Filter(),
		&result,
	); err != nil {
		return
	}
	if err = PoolRoleHistoryCollection.DeleteMany(ctx, i.Filter()); err != nil {
		return
	}

	return
}

func RoleHistoryAdminNotification(ctx context.Context, crewID *models.CrewParam) (err error) {
	crew := new(models.Crew)
	if err = CrewsCollection.FindOne(ctx, crewID.Match(), crew); err != nil {
		log.Print("No crew found")
	}
	mail := vcago.NewMailData("netzwerk@vivaconagua.org", "pool-backend", "asp_selection_network", "pool", "de")
	mail.AddContent(models.RoleAdminContent(crew))
	vcago.Nats.Publish("system.mail.job", mail)
	return
}
