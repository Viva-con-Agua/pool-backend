package models

import (
	"strings"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	TakingCreate struct {
		Name      string         `json:"name" bson:"name"`
		CrewID    string         `json:"crew_id" bson:"crew_id"`
		NewSource []SourceCreate `json:"new_sources"`
		Comment   string         `json:"comment"`
	}
	TakingUpdate struct {
		ID      string         `json:"id" bson:"_id"`
		Name    string         `json:"name" bson:"name"`
		CrewID  string         `json:"crew_id" bson:"crew_id"`
		Sources []SourceUpdate `json:"sources" bson:"-"`
		State   *TakingState   `json:"-;omitempty" bson:"state"`
		Comment string         `json:"comment"`
	}

	TakingDatabase struct {
		ID       string        `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		CrewID   string        `json:"crew_id" bson:"crew_id"`
		Type     string        `json:"type" bson:"type"`
		Comment  string        `json:"comment" bson:"comment"`
		State    TakingState   `json:"state" bson:"state"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	Taking struct {
		ID           string              `json:"id" bson:"_id"`
		Name         string              `json:"name" bson:"name"`
		Type         string              `json:"type" bson:"type"`
		CrewID       string              `json:"crew_id" bson:"crew_id"`
		Crew         Crew                `json:"crew" bson:"crew"`
		Event        Event               `json:"event" bson:"event"`
		Source       []Source            `json:"sources" bson:"sources"`
		State        TakingState         `json:"state" bson:"state"`
		Comment      string              `json:"comment" bson:"comment"`
		EditorID     string              `json:"editor_id" bson:"-"`
		DepositUnits []DepositUnitTaking `json:"deposit_units" bson:"deposit_units"`
		Activities   []Activity          `json:"activities" bson:"activities"`
		Money        vmod.Money          `json:"money" bson:"money"`
		Creator      UserDatabase        `json:"creator" bson:"creator"`
	}
	TakingState struct {
		Open      vmod.Money `json:"open" bson:"open"`
		Confirmed vmod.Money `json:"confirmed" bson:"confirmed"`
		Wait      vmod.Money `json:"wait" bson:"wait"`
	}
	TakingParam struct {
		ID     string `param:"id"`
		CrewID string `param:"crew_id"`
	}
	TakingQuery struct {
		ID              []string `query:"id"`
		Name            string   `query:"name"`
		Search          string   `query:"search"`
		CrewID          []string `query:"crew_id"`
		EventName       string   `query:"event_name"`
		EventState      []string `query:"event_state"`
		EventEndFrom    string   `query:"event_end_from"`
		EventEndTo      string   `query:"event_end_to"`
		Status          []string `query:"status"`
		StatusOpen      bool     `query:"status_open"`
		StatusConfirmed bool     `query:"status_confirmed"`
		StatusNone      bool     `query:"status_none"`
		StatusWait      bool     `query:"status_wait"`
	}
)

var TakingCollection = "takings"

func TakingPermission(token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate("finance")) {
		return vcago.NewPermissionDenied(DepositCollection)
	}
	return
}

func TakingPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.Lookup(DepositUnitTakingView, "_id", "taking_id", "deposit_units")
	pipe.LookupMatch(DepositUnitTakingView, "_id", "taking_id", "wait", bson.D{{Key: "deposit.status", Value: bson.D{{Key: "$in", Value: bson.A{"wait", "open"}}}}})
	pipe.LookupMatch(DepositUnitTakingView, "_id", "taking_id", "confirmed", bson.D{{Key: "deposit.status", Value: "confirmed"}})
	pipe.Lookup(SourceCollection, "_id", "taking_id", "sources")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	pipe.LookupUnwind(EventCollection, "_id", "taking_id", "event")
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.wait.amount", Value: bson.D{{Key: "$sum", Value: "$wait.money.amount"}}},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.confirmed.amount", Value: bson.D{{Key: "$sum", Value: "$confirmed.money.amount"}}},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "money.amount", Value: bson.D{{Key: "$sum", Value: "$sources.money.amount"}}}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "state.open.amount", Value: bson.D{
		{Key: "$subtract", Value: bson.A{"$money.amount", bson.D{{Key: "$add", Value: bson.A{"$state.wait.amount", "$state.confirmed.amount"}}}}},
	}}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.wait.currency", Value: "$currency"},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "state.confirmed.currency", Value: "$currency"},
	}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "money.currency", Value: "$currency"}}}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "state.open.currency", Value: "$currency"}}}})
	pipe.Lookup(ActivityUserView, "_id", "model_id", "activities")
	pipe.LookupUnwindMatch(ActivityUserView, "_id", "model_id", "dummy", bson.D{{Key: "status", Value: "created"}})
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "creator", Value: "$dummy.user"}}}})
	return pipe
}

func (i *TakingCreate) TakingDatabase() *TakingDatabase {
	return &TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "manually",
		Comment:  i.Comment,
		Modified: vmod.NewModified(),
	}
}

func (i *TakingCreate) SourceList(id string) *SourceList {
	r := new(SourceList)
	for n := range i.NewSource {
		source := i.NewSource[n].Source()
		source.TakingID = id
		*r = append(*r, *source)
	}
	return r
}
func (i *TakingUpdate) SourceList(id string) *SourceList {
	r := new(SourceList)
	for _, v := range i.Sources {
		if v.ID == "" {
			source := v.Source()
			source.TakingID = id
			*r = append(*r, *source)
		}
	}
	return r
}

/*
Should this be used somewhere?
func (i *TakingUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}
*/

func (i *TakingQuery) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualStringList("crew_id", []string{token.CrewID})
	} else {
		filter.EqualStringList("crew_id", i.CrewID)
	}
	filter.LikeString("name", i.Name)
	filter.EqualStringList("event.event_state.state", i.EventState)
	filter.LikeString("event.name", i.EventName)
	filter.GteInt64("event.end_at", i.EventEndFrom)
	filter.LteInt64("event.end_at", i.EventEndTo)
	status := bson.A{}
	if i.StatusOpen || i.StatusConfirmed || i.StatusWait || i.StatusNone {
		if i.StatusOpen {
			status = append(status, bson.D{{Key: "state.open.amount", Value: bson.D{{Key: "$gte", Value: 1}}}})
		}
		if i.StatusConfirmed {
			status = append(status, bson.D{{Key: "state.confirmed.amount", Value: bson.D{{Key: "$gte", Value: 1}}}})
		}
		if i.StatusWait {
			status = append(status, bson.D{{Key: "state.wait.amount", Value: bson.D{{Key: "$gte", Value: 1}}}})
		}
		if i.StatusNone {
			status = append(status, bson.D{
				{Key: "state.wait.amount", Value: 0},
				{Key: "state.confirmed.amount", Value: 0},
				{Key: "state.open.amount", Value: 0},
			})
		}
		filter.Append(bson.E{Key: "$or", Value: status})
	}
	return filter.Bson()
}

func (i *TakingUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *TakingParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !token.Roles.Validate("employee;admin") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *Taking) UpdatePermission(token *vcapool.AccessToken) error {
	if i.Event.ID != "" {
		if !token.Roles.Validate("employee;admin") {
			if !token.PoolRoles.Validate("finance") {
				return vcago.NewPermissionDenied("taking")
			}
			if !strings.Contains("published finished", i.Event.EventState.State) {
				return vcago.NewBadRequest("taking", "event_failure")
			}
		}
	}
	return nil
}
