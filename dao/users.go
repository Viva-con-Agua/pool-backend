package dao

import (
	"context"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"

	"go.mongodb.org/mongo-driver/bson"
)

func UserInsert(ctx context.Context, i *models.UserDatabase) (result *models.User, err error) {
	//create mailbox
	mailbox := models.NewMailboxDatabase("user")
	if err = MailboxCollection.InsertOne(ctx, mailbox); err != nil {
		return
	}
	// refer the mailbox.ID
	i.MailboxID = mailbox.ID
	// insert user
	if err = UserCollection.InsertOne(ctx, i); err != nil {
		return
	}
	//select user from database
	if err = UserCollection.AggregateOne(
		ctx,
		models.UserPipeline(false).Match(models.UserMatch(i.ID)).Pipe,
		&result,
	); err != nil {
		return
	}
	return
}

func UsersGet(i *models.UserQuery, token *models.AccessToken) (result *[]models.ListUser, listSize int64, err error) {
	if err = models.UsersPermission(token); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filter := i.PermittedFilter(token)
	sort := i.Sort()
	pipeline := models.SortedUserPermittedPipeline(token).SortFields(sort).Match(filter).Sort(sort).Skip(i.Skip, 0).Limit(i.Limit, 100).Pipe
	result = new([]models.ListUser)
	if err = UserCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	count := vmod.Count{}
	var cErr error
	ctxCount, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if cErr = UserCollection.AggregateOne(ctxCount, models.SortedUserPermittedPipeline(token).Match(filter).Limit(500, 500).Count().Pipe, &count); cErr != nil {
		log.Print(cErr)
		listSize = 1
	} else {
		listSize = int64(count.Total)
	}
	return
}

func UsersGetByCrew(ctx context.Context, i *models.UserQuery, token *models.AccessToken) (result *[]models.UserBasic, err error) {
	if err = i.CrewUsersPermission(token); err != nil {
		return
	}
	filter := i.PermittedUserFilter(token)
	result = new([]models.UserBasic)
	if err = UserCollection.Aggregate(ctx, models.UserPipelinePublic().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func UsersUserGetByID(ctx context.Context, i *models.UserParam, token *models.AccessToken) (result *models.User, err error) {
	if err = models.UsersPermission(token); err != nil {
		return
	}
	if err = UserCollection.AggregateOne(ctx, models.UserPermittedPipeline(token).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func UsersUserGetByIDAPIKey(ctx context.Context, i *models.UserParam) (result *models.User, err error) {
	if err = UserCollection.AggregateOne(ctx, models.UserPipeline(true).Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func UsersGetByID(ctx context.Context, i *models.UserParam) (result *models.User, err error) {
	if err = UserCollection.AggregateOne(ctx, models.UserPipelinePublic().Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func UsersMinimalGet(ctx context.Context, i *models.UserQuery, token *models.AccessToken) (result *[]models.UserMinimal, err error) {
	filter := i.PermittedFilter(token)
	result = new([]models.UserMinimal)
	if err = UserCollection.Aggregate(ctx, models.UserPipelinePublic().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func UsersDeleteUser(ctx context.Context, i *models.UserParam, token *models.AccessToken) (err error) {
	if err = i.UsersDeletePermission(token); err != nil {
		return
	}
	if err = UserDelete(ctx, i.ID); err != nil {
		return
	}
	return
}

func UserDelete(ctx context.Context, id string) (err error) {
	user := new(models.User)
	if err = UserCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}, user); err != nil {
		return
	}
	delete := bson.D{{Key: "user_id", Value: id}}
	if err = AddressesCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = ProfileCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = UserCrewCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	/*
		if err = ActiveCollection.TryDeleteOne(ctx, delete); err != nil {
			return
		}*/
	/*
		if err = NVMCollection.TryDeleteOne(ctx, delete); err != nil {
			return
		}
		if err = NVMCollection.TryDeleteMany(ctx, delete); err != nil {
			return
		}*/
	/*
		if err = AvatarCollection.TryDeleteOne(ctx, delete); err != nil {
			return
		}*/
	if err = MailboxCollection.TryDeleteOne(ctx, bson.D{{Key: "_id", Value: user.MailboxID}}); err != nil {
		return
	}
	if err = MessageCollection.TryDeleteMany(ctx, bson.D{{Key: "mailbox_id", Value: user.MailboxID}}); err != nil {
		return
	}
	if err = UserCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		return
	}
	return
}

func UserSync(ctx context.Context, i *models.ProfileParam, token *models.AccessToken) (result *models.User, err error) {
	profile := new(models.Profile)
	if err = ProfileCollection.FindOne(ctx, i.Match(), profile); err != nil {
		return
	}
	if result, err = ProfileGetByID(ctx, &models.UserParam{ID: profile.UserID}, token); err != nil {
		return
	}
	if err = IDjango.Post(result, "/v1/pool/user/"); err != nil {
		log.Print(err)
		err = nil
	}
	return
}

func UserOrganisationUpdate(ctx context.Context, i *models.UserOrganisationUpdate, token *models.AccessToken) (result *models.User, err error) {
	if err = token.AccessPermission(); err != nil {
		return
	}
	if err = UserCollection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(i),
		&result,
	); err != nil {
		return
	}
	return
}
