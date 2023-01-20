package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	DepositCreate struct {
		DepositUnit []DepositUnitCreate `json:"deposit_units"`
		CrewID      string              `json:"crew_id"`
		HasExternal bool                `json:"has_external"`
		External    External            `json:"external"`
	}
	DepositUnitCreate struct {
		TakingID string     `json:"taking_id" bson:"taking_id"`
		Money    vmod.Money `json:"money" bson:"money"`
	}
	DepositUnit struct {
		ID        string         `json:"id" bson:"_id"`
		TakingID  string         `json:"taking_id" bson:"taking_id"`
		Taking    TakingDatabase `json:"taking" bson:"taking"`
		Money     vmod.Money     `json:"money" bson:"money"`
		DepositID string         `json:"deposit_id" bson:"deposit_id"`
		Status    string         `json:"status" bson:"status"`
		Modified  vmod.Modified  `json:"modified" bson:"modified"`
	}
	DepositUnitTaking struct {
		ID        string          `json:"id" bson:"_id"`
		TakingID  string          `json:"taking_id" bson:"taking_id"`
		Taking    TakingDatabase  `json:"taking" bson:"taking"`
		Money     vmod.Money      `json:"money" bson:"money"`
		DepositID string          `json:"deposit_id" bson:"deposit_id"`
		Deposit   DepositDatabase `json:"deposit" bson:"deposit"`
		Status    string          `json:"status" bson:"status"`
		Modified  vmod.Modified   `json:"modified" bson:"modified"`
	}
	DepositUnitUpdate struct {
		ID          string     `json:"id" bson:"_id"`
		TakingID    string     `json:"taking_id" bson:"taking_id"`
		Money       vmod.Money `json:"money" bson:"money"`
		DepositID   string     `json:"deposit_id" bson:"deposit_id"`
		Status      string     `json:"status" bson:"status"`
		UpdateState string     `json:"update_state" bson:"-"`
	}
	DepositUpdate struct {
		ID          string              `json:"id" bson:"_id"`
		Status      string              `json:"status" bson:"status"`
		DepositUnit []DepositUnitUpdate `json:"deposit_units" bson:"-"`
		HasExternal bool                `json:"has_external" bson:"has_external"`
		External    External            `json:"external" bson:"external"`
		Money       vmod.Money          `json:"money" bson:"money"`
	}
	DepositDatabase struct {
		ID               string        `json:"id" bson:"_id"`
		ReasonForPayment string        `json:"reason_for_payment" bson:"reason_for_payment"`
		Status           string        `json:"status" bson:"status"`
		Money            vmod.Money    `json:"money" bson:"money"`
		CrewID           string        `json:"crew_id" bson:"crew_id"`
		CreatorID        string        `json:"creator_id" bson:"creator_id"`
		ConfirmerID      string        `json:"confirmer_id" bson:"confirmer_id"`
		HasExternal      bool          `json:"has_external" bson:"has_external"`
		External         External      `json:"external" bson:"external"`
		Modified         vmod.Modified `json:"modified" bson:"modified"`
	}
	Deposit struct {
		ID               string        `json:"id" bson:"_id" `
		ReasonForPayment string        `json:"reason_for_payment" bson:"reason_for_payment"`
		Status           string        `json:"status" bson:"status"`
		DepositUnit      []DepositUnit `json:"deposit_units" bson:"deposit_units"`
		CrewID           string        `json:"crew_id" bson:"crew_id"`
		Crew             Crew          `json:"crew" bson:"crew"`
		Money            vmod.Money    `json:"money" bson:"money"`
		Creator          User          `json:"creator" bson:"creator"`
		Confirmer        User          `json:"confirmer" bson:"confirmer"`
		HasExternal      bool          `json:"has_external" bson:"has_external"`
		External         External      `json:"external" bson:"external"`
		Modified         vmod.Modified `json:"modified" bson:"modified"`
	}
	DepositQuery struct {
		ID          []string `query:"id"`
		CrewID      string   `query:"crew_id"`
		UpdatedTo   string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom string   `query:"updated_from" qs:"updated_from"`
		CreatedTo   string   `query:"created_to" qs:"created_to"`
		CreatedFrom string   `query:"created_from" qs:"created_from"`
	}
	DepositParam struct {
		ID string `param:"id"`
	}
)

func (i *DepositCreate) DepositDatabase(token *vcapool.AccessToken) (r *DepositDatabase, d []DepositUnit) {
	dIDs := []string{}
	d = []DepositUnit{}
	var amount int64 = 0
	id := uuid.NewString()
	for _, value := range i.DepositUnit {
		depositUnit := &DepositUnit{
			ID:        uuid.NewString(),
			TakingID:  value.TakingID,
			Money:     value.Money,
			DepositID: id,
			Status:    "open",
			Modified:  vmod.NewModified(),
		}
		dIDs = append(dIDs, depositUnit.ID)
		d = append(d, *depositUnit)
		amount += depositUnit.Money.Amount
	}
	currency := "EUR"
	if d != nil {
		currency = d[0].Money.Currency
	}
	r = &DepositDatabase{
		ID:     id,
		Status: "open",
		Money: vmod.Money{
			Amount:   amount,
			Currency: currency,
		},
		CrewID:      i.CrewID,
		CreatorID:   token.ID,
		HasExternal: i.HasExternal,
		External:    i.External,
		Modified:    vmod.NewModified(),
	}
	return
}

func (i *DepositQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.EqualString("crew_id", i.CrewID)
	return filter.Bson()
}
