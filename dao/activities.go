package dao

import (
	"context"
	"errors"
	"pool-backend/models"
)

func SystemActivity(ctx context.Context, activity string, modelID string, userID string) (err error) {
	var model, message, status string
	switch activity {
	case "taking_updated":
		model, message, status = "taking", "Successfully updated", "updated"
	case "taking_created":
		model, message, status = "taking", "Successfully created", "created"
	default:
		return errors.New("ACTIVITY_NOT_FOUND")
	}
	err = ActivityInsert(ctx, userID, model, modelID, message, status)
	return
}

func ActivityInsert(ctx context.Context, userID string, model string, modelID string, message, status string) (err error) {
	activity := models.NewActivityDB(userID, model, modelID, message, status)
	err = ActivityCollection.InsertOne(ctx, activity)
	return
}
