package dao

import (
	"context"
	"log"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"go.mongodb.org/mongo-driver/bson"
)

func NewsletterCreate(ctx context.Context, i *models.NewsletterCreate, token *models.AccessToken) (result *models.Newsletter, err error) {

	if !token.Roles.Validate("admin;employee;pool_employee") || i.UserID == "" {
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
			filter := bson.D{{Key: "_id", Value: i.UserID}, {Key: "crew._id", Value: bson.D{{Key: "$ne", Value: ""}}}}
			if err = UserCollection.FindOne(ctx, filter, crew); err != nil {
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

func NewsletterDelete(ctx context.Context, i *models.NewsletterParam, token *models.AccessToken) (result *models.Newsletter, err error) {
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
		user := new(models.User)
		filter := bson.D{{Key: "_id", Value: user.ID}, {Key: "crew._id", Value: bson.D{{Key: "$ne", Value: ""}}}}
		if err = UserCollection.FindOne(ctx, filter, user); err != nil {
			return nil, vcago.NewBadRequest(models.NewsletterCollection, "not part of an crew", nil)
		}
	}
	result = i.ToNewsletter(user.ID)
	if err = NewsletterCollection.InsertOne(ctx, result); err != nil {
		return
	}
	return
}

func NewsletterSync(ctx context.Context, i *models.User, token *models.AccessToken) (result *[]models.Newsletter, err error) {
	export := &models.NewsletterExport{
		UserID:     i.ID,
		Newsletter: i.Newsletter,
	}
	go func() {
		if err = IDjango.Post(export, "/v1/pool/newsletters/"); err != nil {
			log.Print(err)
		}
	}()
	result = &i.Newsletter
	return
}
