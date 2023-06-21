package dao

import (
	"context"
	"fmt"
	"pool-backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// GetNewReasonForPayment create a unique string contains the crew abbreviation,
// the year of the payment and a number.
func GetNewReasonForPayment(ctx context.Context, crewID string) (r string, err error) {
	year := strconv.Itoa(time.Now().Year())
	rfp := new(models.ReasonForPayment)
	if err = ReasonForPaymentCollection.FindOne(
		ctx,
		bson.D{{Key: "crew_id", Value: crewID}, {Key: "year", Value: year}},
		rfp,
	); err != nil {
		if vmdb.ErrNoDocuments(err) {
			err = nil
			crew := new(models.Crew)
			if err = CrewsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: crewID}}, crew); err != nil {
				return
			}
			rfp = &models.ReasonForPayment{
				ID:           uuid.NewString(),
				Abbreviation: crew.Abbreviation,
				CrewID:       crewID,
				Year:         year,
				Count:        1,
			}
			if err = ReasonForPaymentCollection.InsertOne(ctx, rfp); err != nil {
				return
			}
		} else {
			return
		}
	}
	count := fmt.Sprintf("%05d", rfp.Count)
	r = "POOL-" + strings.ToUpper(rfp.Abbreviation) + "-" + rfp.Year + "-" + count
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}}}
	if err = ReasonForPaymentCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: rfp.ID}}, update, nil); err != nil {
		return
	}
	return
}
