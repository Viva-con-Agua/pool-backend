package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	TakingCreate struct {
		Name      string         `json:"name" bson:"name"`
		CrewID    string         `json:"crew_id" bson:"crew_id"`
		External  External       `json:"external" bson:"external"`
		NewSource []SourceCreate `json:"new_sources"`
		Comment   string         `json:"comment"`
	}
	TakingUpdate struct {
		ID       string         `json:"id" bson:"_id"`
		Name     string         `json:"name" bson:"name"`
		CrewID   string         `json:"crew_id" bson:"crew_id"`
		External External       `json:"external" bson:"external"`
		Sources  []SourceUpdate `json:"sources" bson:"-"`
		State    *TakingState   `json:"-;omitempty" bson:"state"`
		Comment  string         `json:"comment"`
	}
	External struct {
		Organisation string `json:"organisation" bson:"organisation"`
		ASP          string `json:"asp" bson:"asp"`
		Email        string `json:"email" bson:"email"`
		Address      string `json:"address" bson:"address"`
		Reciept      bool   `json:"reciept" bson:"reciept"`
		Purpose      string `json:"purpose" bson:"purpose"`
	}
	TakingDatabase struct {
		ID       string        `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		CrewID   string        `json:"crew_id" bson:"crew_id"`
		Type     string        `json:"type" bson:"type"`
		External External      `json:"external" bson:"external"`
		Comment  string        `json:"comment" bson:"comment"`
		Status   string        `json:"status" bson:"status"`
		State    TakingState   `json:"state" bson:"state"`
		Modified vmod.Modified `json:"modified" bson:"modified"`
	}
	Taking struct {
		ID       string      `json:"id" bson:"_id"`
		Name     string      `json:"name" bson:"name"`
		Type     string      `json:"type" bson:"type"`
		CrewID   string      `json:"crew_id" bson:"crew_id"`
		Crew     Crew        `json:"crew" bson:"crew"`
		Event    Event       `json:"event" bson:"event"`
		External External    `json:"external" bson:"external"`
		Source   []Source    `json:"sources" bson:"sources"`
		Status   string      `json:"status" bson:"status"`
		State    TakingState `json:"state" bson:"state"`
		Comment  string      `json:"comment" bson:"comment"`
		Activity []Activity  `json:"activity" bson:"activity"`
	}
	TakingState struct {
		Confirmed vmod.Money `json:"confirmed" bson:"confirmed"`
		Open      vmod.Money `json:"open" bson:"open"`
		Wait      vmod.Money `json:"wait" bson:"wait"`
	}
	TakingParam struct {
		ID string `param:"id"`
	}
	TakingQuery struct {
		ID []string `query:"id"`
	}
)

func (i *TakingCreate) TakingDatabase() *TakingDatabase {
	takingState := new(TakingState)
	takingState.Open.Amount = 0
	for n := range i.NewSource {
		takingState.Open.Amount += i.NewSource[n].Money.Amount
	}
	takingState.Wait.Amount = 0
	takingState.Confirmed.Amount = 0
	return &TakingDatabase{
		ID:       uuid.NewString(),
		Name:     i.Name,
		CrewID:   i.CrewID,
		Type:     "manually",
		External: i.External,
		Comment:  i.Comment,
		Status:   "open",
		State:    *takingState,
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

func NewTakingsPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.Lookup("sources", "_id", "taking_id", "sources")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	pipe.LookupUnwind("events", "_id", "taking_id", "event")
	return pipe
}

func (i *TakingQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	return filter.Bson()
}

func (i *TakingUpdate) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *TakingUpdate) Update() bson.D {
	return vmdb.UpdateSet(i)
}
func (i *TakingParam) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *TakingQuery) Pipeline() *vmdb.Pipeline {
	return NewTakingsPipeline().Match(i.Filter())
}

func (i *TakingParam) Pipeline() *vmdb.Pipeline {
	return NewTakingsPipeline().Match(i.Filter())
}
