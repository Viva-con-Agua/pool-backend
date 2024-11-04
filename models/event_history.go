package models

import (
	"fmt"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type EventStateHistoryCreate struct {
	ID       string `bson:"_id" json:"id"`
	UserID   string `json:"user_id" bson:"user_id"`
	CrewID   string `json:"crew_id" bson:"crew_id"`
	EventID  string `json:"event_id" bson:"event_id"`
	Date     int64  `json:"date" bson:"date"`
	OldState string `json:"old_state" bson:"old_state"`
	NewState string `json:"new_state" bson:"new_state"`
}

type EventStateHistoryQuery struct {
	ID            []string `query:"id" qs:"id"`
	UserID        string   `query:"user_id" qs:"user_id"`
	CrewID        string   `query:"crew_id" qs:"crew_id"`
	EventID       string   `query:"event_id" qs:"event_id"`
	DateTo        string   `query:"date_to" qs:"date_to"`
	DateFrom      string   `query:"date_from" qs:"date_from"`
	OldState      string   `query:"old_state" qs:"old_state"`
	NewState      string   `query:"new_state" qs:"new_state"`
	Search        string   `query:"search"`
	SortField     string   `query:"sort"`
	SortDirection string   `query:"sort_dir"`
	Limit         int64    `query:"limit"`
	Skip          int64    `query:"skip"`
	FullCount     string   `query:"full_count"`
}

type EventStateHistory struct {
	ID       string      `bson:"_id" json:"id"`
	UserID   string      `json:"user_id" bson:"user_id"`
	User     UserMinimal `json:"user" bson:"user"`
	CrewID   string      `json:"crew_id" bson:"crew_id"`
	Crew     CrewSimple  `json:"crew" bson:"crew"`
	EventID  string      `json:"event_id" bson:"event_id"`
	Event    EventPublic `json:"event" bson:"event"`
	Date     int64       `json:"date" bson:"date"`
	OldState string      `json:"old_state" bson:"old_state"`
	NewState string      `json:"new_state" bson:"new_state"`
}

type EventStateHistoryNotification struct {
	EventID       string `json:"event_id" bson:"event_id"`
	EventName     string `json:"event_name" bson:"event_name"`
	EventStart    string `json:"event_start" bson:"event_start"`
	EventCrew     string `json:"event_crew" bson:"event_crew"`
	EventArtist   string `json:"event_artist" bson:"event_artist"`
	EventLocation string `json:"event_location" bson:"event_location"`
	PublishedDate string `json:"published_date" bson:"published_date"`
}

var EventStateHistoryCollection = "eventstate_history"

func EventStatePipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(EventCollection, "event_id", "_id", "event")
	pipe.Lookup(ArtistCollection, "event.artist_ids", "_id", "event.artists")
	pipe.LookupUnwind(UserCollection, "user_id", "_id", "user")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	return
}

func EventStateHistoryPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(PoolRoleHistoryCollection)
	}
	return
}

func (i *Event) NewEventStateHistory(old string, new string, token *AccessToken) *EventStateHistoryCreate {
	return &EventStateHistoryCreate{
		ID:       uuid.NewString(),
		UserID:   token.ID,
		CrewID:   i.CrewID,
		EventID:  i.ID,
		OldState: old,
		NewState: new,
		Date:     time.Now().Unix(),
	}
}

func (i *EventStateHistoryQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualString("user_id", i.UserID)
	filter.EqualString("event_id", i.EventID)
	filter.EqualString("crew_id", i.CrewID)
	filter.EqualString("old_state", i.OldState)
	filter.EqualString("new_state", i.NewState)
	filter.GteInt64("date", i.DateFrom)
	filter.LteInt64("date", i.DateTo)
	filter.SearchString([]string{"_id", "name", "crew.name"}, i.Search)

	return filter.Bson()
}

func EventPublishedLastWeek() bson.D {
	filter := vmdb.NewFilter()

	week := 7 * 24
	// if we had fix number of units to subtract, we can use following line instead fo above 2 lines. It does type convertion automatically.
	// then := now.Add(-10 * time.Minute)
	fmt.Printf("7 week ago: %v\n ", time.Now().Add(time.Duration(-week)*time.Hour))

	filter.GteInt64("date", fmt.Sprint(time.Now().Add(time.Duration(-week)*time.Hour).Unix()))
	filter.LteInt64("date", fmt.Sprint(time.Now().Unix()))
	return filter.Bson()
}

func EventHistoryAdminContent(data []EventStateHistoryNotification) *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["Events"] = data
	return content
}
