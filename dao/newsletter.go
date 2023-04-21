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
			return nil, vcago.NewBadRequest("newsletter", "not part of an crew", nil)
		}
		result = i.Newsletter(token)
		if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
			return
		}
	} else {
		if i.Value == "regional" {
			crew := new(models.UserCrew)
			if err = UserCrewCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: i.UserID}}, crew); err != nil {
				return nil, vcago.NewBadRequest("newsletter", "not part of an crew", nil)
			}
		}
		result = i.NewsletterAdmin()
		if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
			return
		}
	}

	return
}

func NewsletterDelete(ctx context.Context, i *models.NewsletterParam, token *vcapool.AccessToken) (newletter *models.Newsletter, err error) {
	newletter = new(models.Newsletter)
	filter := bson.D{{Key: "_id", Value: i.ID}}
	if err = NewsletterCollection.FindOne(ctx, filter, newletter); err != nil {
		return
	}
	if !token.Roles.Validate("employee;admin") {
		if token.ID != newletter.UserID {
			return nil, vcago.NewPermissionDenied("newsletter", i.ID)
		}
	}
	if err = NewsletterCollection.DeleteOne(ctx, filter); err != nil {
		return
	}
	return
}
