package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ParticipationCreate struct {
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
	}

	ParticipationUpdate struct {
		ID      string `json:"id" bson:"_id"`
		Status  string `json:"status" bson:"status"`
		Comment string `json:"comment" bson:"comment"`
		//Confirmer UserInternal `json:"confirmer" bson:"confirmer"`
	}
	ParticipationDatabase struct {
		ID      string `json:"id" bson:"_id"`
		UserID  string `json:"user_id" bson:"user_id"`
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
		Status  string `json:"status" bson:"status"`
		CrewID  string `json:"crew_id" bson:"crew_id"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	Participation struct {
		ID      string `json:"id" bson:"_id"`
		UserID  string `json:"user_id" bson:"user_id"`
		User    User   `json:"user" bson:"user"`
		EventID string `json:"event_id" bson:"event_id"`
		Comment string `json:"comment" bson:"comment"`
		Status  string `json:"status" bson:"status"`
		Event   Event  `json:"event" bson:"event"`
		CrewID  string `json:"crew_id" bson:"crew_id"`
		Crew    Crew   `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	UserParticipation struct {
		ID      string      `json:"id" bson:"_id"`
		EventID string      `json:"event_id" bson:"event_id"`
		Comment string      `json:"comment" bson:"comment"`
		Status  string      `json:"status" bson:"status"`
		Event   EventPublic `json:"event" bson:"event"`
		CrewID  string      `json:"crew_id" bson:"crew_id"`
		Crew    CrewName    `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	EventParticipation struct {
		ID      string    `json:"id" bson:"_id"`
		UserID  string    `json:"user_id" bson:"user_id"`
		User    User      `json:"user" bson:"user"`
		EventID string    `json:"event_id" bson:"event_id"`
		Comment string    `json:"comment" bson:"comment"`
		Status  string    `json:"status" bson:"status"`
		Event   ListEvent `json:"event" bson:"event"`
		CrewID  string    `json:"crew_id" bson:"crew_id"`
		Crew    Crew      `json:"crew" bson:"crew"`
		//Confirmer UserInternal   `json:"confirmer" bson:"confirmer"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	ParticipationMinimal struct {
		ID      string `json:"id" bson:"_id"`
		EventID string `json:"event_id" bson:"event_id"`
		Status  string `json:"status" bson:"status"`
	}

	ParticipationParam struct {
		ID string `param:"id"`
	}

	ParticipationQuery struct {
		ID       []string `query:"id" qs:"id"`
		EventID  []string `query:"event_id" qs:"event_id"`
		Comment  []string `query:"comment" bson:"comment"`
		Status   []string `query:"status" bson:"status"`
		UserId   []string `query:"user_id" bson:"user_id"`
		CrewName []string `query:"crew_name" bson:"crew_name"`
		CrewId   []string `query:"crew_id" qs:"crew_id"`
	}
	ParticipationStateRequest struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	ParticipationImport struct {
		DropsID string `json:"drops_id"`
		Comment string `json:"comment"`
	}
)

func (i *ParticipationCreate) ParticipationDatabase(token *vcapool.AccessToken) *ParticipationDatabase {
	return &ParticipationDatabase{
		ID:       uuid.NewString(),
		UserID:   token.ID,
		EventID:  i.EventID,
		Comment:  i.Comment,
		Status:   "requested",
		CrewID:   token.CrewID,
		Modified: vmod.NewModified(),
	}
}

func (i *ParticipationImport) ParticipationDatabase() *ParticipationDatabase {
	return &ParticipationDatabase{
		ID:       uuid.NewString(),
		Comment:  i.Comment,
		Status:   "requested",
		Modified: vmod.NewModified(),
	}
}

func ParticipationPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind("users", "user_id", "_id", "user")
	pipe.LookupUnwind("profiles", "user_id", "user_id", "user.profile")
	pipe.LookupUnwind("events", "event_id", "_id", "event")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	pipe.LookupUnwind("users", "event.event_asp_id", "_id", "event.event_asp")
	pipe.LookupUnwind("profiles", "event.event_asp_id", "user_id", "event.event_asp.profile")
	pipe.LookupUnwind("users", "event.internal_asp_id", "_id", "event.internal_asp")
	pipe.LookupUnwind("profiles", "event.internal_asp_id", "user_id", "event.internal_asp.profile")
	pipe.LookupUnwind("users", "event.creator_id", "_id", "event.creator")
	pipe.LookupUnwind("profiles", "event.creator_id", "user_id", "event.creator.profile")
	pipe.Lookup("artists", "event.artist_ids", "_id", "event.artists")
	pipe.LookupUnwind("organizers", "event.organizer_id", "_id", "event.organizer")
	pipe.LookupUnwind("crews", "event.crew_id", "_id", "event.crew")
	return
}

func ParticipationAspPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind("events", "event_id", "_id", "event")
	pipe.LookupUnwind("users", "event.event_asp_id", "_id", "event.event_asp")
	pipe.LookupUnwind("profiles", "event.event_asp_id", "user_id", "event.event_asp.profile")
	return
}

func (i *ParticipationQuery) Filter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualStringList("event_id", i.EventID)
	filter.EqualStringList("status", i.Status)
	filter.EqualStringList("comment", i.Comment)
	filter.EqualStringList("user_id", i.UserId)
	filter.EqualStringList("crew.name", i.CrewName)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	} else {
		filter.EqualStringList("crew_id", i.CrewId)
	}
	return filter.Bson()
}

func (i *ParticipationQuery) FilterUser(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualStringList("event_id", i.EventID)
	filter.EqualStringList("status", i.Status)
	filter.EqualStringList("comment", i.Comment)
	filter.EqualStringList("crew.name", i.CrewName)
	filter.EqualStringList("crew_id", i.CrewId)
	filter.EqualString("user_id", token.ID)
	return filter
}

func (i *ParticipationQuery) FilterAspInformation(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualStringList("event_id", i.EventID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("status", "confirmed")
		filter.EqualString("user_id", token.ID)
	}
	return filter
}

func (i *EventParam) FilterEvent(token *vcapool.AccessToken) *vmdb.Filter {
	filter := vmdb.NewFilter()
	filter.EqualString("event_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		filter.EqualString("event.event_asp_id", token.ID)
	} else if !token.Roles.Validate("employee;admin") {
		filter.EqualString("event.crew_id", token.CrewID)
	}
	return filter
}

func (i *ParticipationDatabase) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *ParticipationUpdate) Match() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	return match.Bson()
}

func (i *ParticipationUpdate) Filter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualString("user_id", token.ID)
	} else if !token.Roles.Validate("employee;admin") {
		match.EqualString("crew_id", token.CrewID)
	}
	return match.Bson()
}

func (i *ParticipationParam) Filter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualString("user_id", token.ID)
	} else if !token.Roles.Validate("employee;admin") {
		match.EqualString("crew_id", token.CrewID)
	}
	return match.Bson()
}

func (i *ParticipationStateRequest) PermittedFilter(token *vcapool.AccessToken) bson.D {
	if i.Status == "withdrawn" {
		return bson.D{{Key: "_id", Value: i.ID}, {Key: "user", Value: token.ID}}
	} else if i.Status == "confirmed" || i.Status == "rejected" {
		if token.Roles.Validate("employee;admin") {
			return bson.D{{Key: "_id", Value: i.ID}}
		} else if token.Roles.Validate("operation") {
			return bson.D{{Key: "_id", Value: i.ID}, {Key: "crew.id", Value: token.CrewID}}
		}
	}
	return bson.D{{Key: "_id", Value: "not_defined"}}
}

func (i *ParticipationStateRequest) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}
