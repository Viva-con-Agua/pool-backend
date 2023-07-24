package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func NewsletterCreate(ctx context.Context, i *models.NewsletterCreate, token *vcapool.AccessToken) (result *models.Newsletter, err error) {

	if !token.Roles.Validate("employee;admin") || i.UserID == "" {
		if i.Value == "regional" && token.CrewID == "" {
			return nil, vcago.NewBadRequest(models.NewsletterCollection, "not part of an crew", nil)
		}
		result = i.Newsletter(token)
		if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
			return
		}
	} else {
		if i.Value == "regional" {
			crew := new(models.UserCrew)
			if err = UserCrewCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: i.UserID}}, crew); err != nil {
				return nil, vcago.NewBadRequest(models.NewsletterCollection, "not part of an crew", nil)
			}
		}
		result = i.NewsletterAdmin()
		if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
			return
		}
	}

	return
}

func NewsletterDelete(ctx context.Context, i *models.NewsletterParam, token *vcapool.AccessToken) (result *models.Newsletter, err error) {
	filter := i.Match()
	if err = NewsletterCollection.FindOne(ctx, filter, &result); err != nil {
		return
	}
	if err = result.DeletePermission(token); err != nil {
		return
	}
	if err = NewsletterCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}

func NewsletterImport(ctx context.Context, i *models.NewsletterImport) (result *models.Newsletter, err error) {
	user := new(models.UserDatabase)
	userFilter := bson.D{{Key: "drops_id", Value: i.DropsID}}
	if err = UserCollection.FindOne(ctx, userFilter, user); err != nil {
		return
	}
	if i.Value == "regional" {
		crew := new(models.UserCrew)
		if err = UserCrewCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: user.ID}}, crew); err != nil {
			return nil, vcago.NewBadRequest(models.NewsletterCollection, "not part of an crew", nil)
		}
	}
	result = i.ToNewsletter(user.ID)
	if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}
