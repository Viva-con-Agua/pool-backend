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
		Currency  string         `json:"currency"`
		NewSource []SourceCreate `json:"new_sources"`
		Comment   string         `json:"comment"`
	}
	TakingUpdate struct {
		ID       string         `json:"id" bson:"_id"`
		Name     string         `json:"name" bson:"name"`
		CrewID   string         `json:"crew_id" bson:"crew_id"`
		Currency string         `json:"currency" bson:"currency"`
		Sources  []SourceUpdate `json:"sources" bson:"-"`
		State    *TakingState   `json:"-;omitempty" bson:"state"`
		Comment  string         `json:"comment"`
	}

	TakingDatabase struct {
		ID       string        `json:"id" bson:"_id"`
		Name     string        `json:"name" bson:"name"`
		CrewID   string        `json:"crew_id" bson:"crew_id"`
		Type     string        `json:"type" bson:"type"`
		Comment  string        `json:"comment" bson:"comment"`
		Currency string        `json:"currency" bson:"currency"`
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
		Currency     string              `json:"currency" bson:"currency"`
		DepositUnits []DepositUnitTaking `json:"deposit_units" bson:"deposit_units"`
		Activity     []Activity          `json:"activity" bson:"activity"`
		Money        vmod.Money          `json:"money" bson:"money"`
	}
	TakingState struct {
		Open      vmod.Money `json:"open" bson:"open"`
		Confirmed vmod.Money `json:"confirmed" bson:"confirmed"`
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
