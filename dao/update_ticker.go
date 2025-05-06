package dao

import (
	"context"
	"fmt"
	"log"
	"pool-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTicker() {
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{})
	EventStateFinishTicker()
	EventStateClosedTicker()
	UserActiveStateTicker()
	//SendWeeklyNotification()
	SendWeeklyCrewNotification()
	go func() {
		for {
			select {
			case <-ticker.C:
				//SendWeeklyNotification()
				//SendWeeklyCrewNotification()
				EventStateFinishTicker()
				EventStateClosedTicker()
				UserActiveStateTicker()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// EventStateFinishTicker sets all published events that are already over to finished
func EventStateFinishTicker() {
	filter := vmdb.NewFilter()
	filter.EqualString("event_state.state", "published")
	filter.LteInt64("end_at", fmt.Sprint(time.Now().Unix()))
	update := bson.D{{Key: "event_state.state", Value: "finished"}}
	if err := EventCollection.UpdateMany(context.Background(), filter.Bson(), vmdb.UpdateSet(update)); err != nil {
		if !vmdb.ErrNoDocuments(err) {
			log.Print(err)
		}
	}
}

// EventStateClosed
func EventStateClosedTicker() {
	filter := vmdb.NewFilter()
	confirmedFilter := bson.E{Key: "$or", Value: bson.A{
		bson.D{
			{Key: "state.confirmed.amount", Value: bson.D{{Key: "$gte", Value: 1}}},
			{Key: "state.wait.amount", Value: 0},
			{Key: "state.open.amount", Value: 0},
		},
		bson.D{{Key: "state.no_income", Value: true}},
	}}
	filter.Append(confirmedFilter)
	filter.EqualString("event.event_state.state", "finished")
	pipeline := models.TakingPipeline().Match(filter.Bson()).Pipe
	takings := []models.Taking{}
	if err := TakingCollection.Aggregate(context.Background(), pipeline, &takings); err != nil {
		log.Print(err)
	}
	for i := range takings {
		updateFilter := bson.D{{Key: "_id", Value: takings[i].Event.ID}}
		update := bson.D{{Key: "event_state.state", Value: "closed"}}
		e := new(models.Event)
		if err := EventCollection.UpdateOne(context.Background(), updateFilter, vmdb.UpdateSet(update), e); err != nil {
			log.Print(err)
		}
		org := new(models.Organisation)
		if err := OrganisationCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: takings[i].Event.OrganisationID}}, &org); err != nil {
			log.Print(err)
		}
		takings[i].EditorID = org.DefaultAspID
		if err := IDjango.Post(takings[i], "/v1/pool/taking/create/"); err != nil {
			log.Print(err)
		}
		if err := IDjango.Post(e, "/v1/pool/event/update/"); err != nil {
			log.Print(err)
		}
		participations := new([]models.Participation)
		if err := ParticipationCollection.Aggregate(
			context.Background(),
			models.ParticipationPipeline().Match(bson.D{{Key: "event_id", Value: e.ID}}).Pipe,
			participations,
		); err != nil {
			log.Print(err)
		}

		if err := IDjango.Post(participations, "/v1/pool/participations/create/"); err != nil {
			log.Print(err)
		}

	}
}

func UserActiveStateTicker() {
	checkDate := time.Now().Unix() - 15768000
	filter := bson.D{
		{Key: "last_login_date", Value: bson.D{{Key: "$lte", Value: checkDate}}},
		{Key: "active.status", Value: "confirmed"},
	}
	userList := []models.User{}
	pipeline := models.UserPipeline(false).Match(filter).Pipe
	if err := UserCollection.Aggregate(context.Background(), pipeline, &userList); err != nil {
		log.Print(err)
	}
	for _, user := range userList {
		update := bson.D{{Key: "status", Value: "rejected"}}
		userFilter := bson.D{{Key: "_id", Value: user.Active.ID}}
		if err := ActiveCollection.UpdateOne(context.Background(), userFilter, vmdb.UpdateSet(update), nil); err != nil {
			log.Print(err)
		}
	}

}

// SendWeeklyNotification Send a notification mail to festival@ and netzwerk@ for all published events last 7 days
func SendWeeklyNotification() {

	h, _, _ := time.Now().Clock()
	if time.Now().Weekday() != time.Monday {
		return
	} else if h < 9 || h >= 10 {
		return
	}

	filter := models.EventPublishedLastWeek()

	events := []models.EventStateHistory{}
	pipeline := models.EventStatePipeline().Match(filter).Pipe
	if err := EventStateHistoryCollection.Aggregate(context.Background(), pipeline, &events); err != nil {
		log.Print(err)
	}

	eventNotifications := []models.EventStateHistoryNotification{}
	for _, e := range events {
		eventNotifications = append(eventNotifications, models.EventStateHistoryNotification{
			EventID:       e.EventID,
			EventName:     e.Event.Name,
			EventCrew:     e.Crew.Name,
			EventStart:    time.Unix(e.Event.StartAt, 0).Format("02.01.2006 15:04"),
			EventArtist:   models.ToArtistList(e.Event.Artists),
			EventLocation: e.Event.GetLocation(),
			PublishedDate: time.Unix(e.Date, 0).Format("02.01.2006"),
		})
	}

	if err := EventHistoryAdminNotification(context.Background(), eventNotifications); err != nil {
		log.Print(err)

	}
}

// SendWeeklyCrewNotification Send a notification mail to each crew email address for all published events last 7 days
func SendWeeklyCrewNotification() {

	h, _, _ := time.Now().Clock()
	if time.Now().Weekday() != time.Monday {
		return
	} else if h < 9 || h >= 10 {
		return
	}

	filter := models.EventPublishedLastWeek()
	events := []models.EventStateHistory{}
	pipeline := models.EventStatePipeline().Match(filter).Pipe
	if err := EventStateHistoryCollection.Aggregate(context.Background(), pipeline, &events); err != nil {
		log.Print(err)
	}
	m := make(map[string][]models.EventStateHistoryNotification)
	for _, e := range events {
		if e.CrewID == "" {
			continue
		}
		m[e.Crew.Email] = append(m[e.Crew.Email], models.EventStateHistoryNotification{
			EventID:       e.EventID,
			EventName:     e.Event.Name,
			EventCrew:     e.Crew.Name,
			EventStart:    time.Unix(e.Event.StartAt, 0).Format("02.01.2006 15:04"),
			EventArtist:   models.ToArtistList(e.Event.Artists),
			EventLocation: e.Event.GetLocation(),
			PublishedDate: time.Unix(e.Date, 0).Format("02.01.2006"),
		})
	}
	if err := EventHistoryCrewNotification(context.Background(), m); err != nil {
		log.Print(err)
	}
}
