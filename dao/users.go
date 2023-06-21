package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"

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

func UsersGet(ctx context.Context, i *models.UserQuery, token *vcapool.AccessToken) (result *[]models.User, list_size int64, err error) {
	if err = models.UsersPermission(token); err != nil {
		return
	}
	ctx = context.Background()
	filter := i.PermittedFilter(token)
	sort := i.Sort()
	pipeline := models.SortedUserPipeline(sort, false).Match(filter).Sort(sort).Skip(i.Skip, 0).Limit(i.Limit, 100).Pipe
	result = new([]models.User)
	if err = UserCollection.Aggregate(ctx, pipeline, result); err != nil {
		return
	}
	//list_size = 1254
	count := vmod.Count{}
	var cErr error
	if cErr = UserCollection.AggregateOne(ctx, models.UserPipeline(false).Match(filter).Count().Pipe, &count); cErr != nil {
		print(cErr)
		list_size = 1
	} else {
		list_size = int64(count.Total)
	}
	return
}

func UsersGetByCrew(ctx context.Context, i *models.UserQuery, token *vcapool.AccessToken) (result *[]models.UserBasic, err error) {
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

func UsersGetByID(ctx context.Context, i *models.UserParam) (result *models.User, err error) {
	if err = UserCollection.AggregateOne(ctx, models.UserPipelinePublic().Match(i.Match()).Pipe, &result); err != nil {
		return
	}
	return
}

func UsersMinimalGet(ctx context.Context, i *models.UserQuery, token *vcapool.AccessToken) (result *[]models.UserMinimal, err error) {
	filter := i.PermittedFilter(token)
	result = new([]models.UserMinimal)
	if err = UserCollection.Aggregate(ctx, models.UserPipelinePublic().Match(filter).Pipe, result); err != nil {
		return
	}
	return
}

func UsersDeleteUser(ctx context.Context, i *models.UserParam, token *vcapool.AccessToken) (err error) {
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
	if err = ActiveCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = NVMCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
	if err = NVMCollection.TryDeleteMany(ctx, delete); err != nil {
		return
	}
	if err = AvatarCollection.TryDeleteOne(ctx, delete); err != nil {
		return
	}
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
