package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	DepositCreate struct {
		DepositUnit []DepositUnitCreate `json:"deposit_units"`
		CrewID      string              `json:"crew_id" validate:"required"`
		HasExternal bool                `json:"has_external"`
		External    External            `json:"external"`
	}
	DepositUnitCreate struct {
		TakingID string     `json:"taking_id" bson:"taking_id"`
		Money    vmod.Money `json:"money" bson:"money"`
	}
	DepositUnit struct {
		ID        string        `json:"id" bson:"_id"`
		TakingID  string        `json:"taking_id" bson:"taking_id"`
		Taking    Taking        `json:"taking" bson:"taking"`
		Money     vmod.Money    `json:"money" bson:"money"`
		DepositID string        `json:"deposit_id" bson:"deposit_id"`
		Status    string        `json:"status" bson:"status"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
	DepositUnitTaking struct {
		ID        string          `json:"id" bson:"_id"`
		TakingID  string          `json:"taking_id" bson:"taking_id"`
		Taking    Taking          `json:"taking" bson:"taking"`
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
		ID            string              `json:"id" bson:"_id"`
		CrewID        string              `json:"crew_id" bson:"crew_id"`
		Status        string              `json:"status" bson:"status"`
		DepositUnit   []DepositUnitUpdate `json:"deposit_units" bson:"-"`
		HasExternal   bool                `json:"has_external" bson:"has_external"`
		External      External            `json:"external" bson:"external"`
		UpdateState   string              `json:"update_state" bson:"-"`
		DateOfDeposit int64               `json:"date_of_deposit" bson:"date_of_deposit"`
		Money         vmod.Money          `json:"money" bson:"money"`
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
		DateOfDeposit    int64         `json:"date_of_deposit" bson:"date_of_deposit"`
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
		Receipts         []ReceiptFile `json:"receipts" bson:"receipts"`
		DateOfDeposit    int64         `json:"date_of_deposit" bson:"date_of_deposit"`
		External         External      `json:"external" bson:"external"`
		Modified         vmod.Modified `json:"modified" bson:"modified"`
	}
	DepositQuery struct {
		ID               []string `query:"id"`
		Name             string   `query:"deposit_unit_name"`
		ReasonForPayment string   `query:"reason_for_payment"`
		CrewID           []string `query:"crew_id"`
		Status           []string `query:"deposit_status"`
		Creator          []string `query:"deposit_creator"`
		Confirmer        []string `query:"deposit_confirmer"`
		HasExternal      string   `query:"deposit_has_external"`
		UpdatedTo        string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom      string   `query:"updated_from" qs:"updated_from"`
		CreatedTo        string   `query:"created_to" qs:"created_to"`
		CreatedFrom      string   `query:"created_from" qs:"created_from"`
		vmdb.Query
	}
	DepositParam struct {
		ID     string `param:"id"`
		CrewID string `param:"crew_id"`
	}
)

var DepositCollection = "deposits"
var DepositUnitCollection = "deposit_units"
var DepositUnitTakingView = "deposit_unit_taking"

func DepositPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate("finance")) {
		return vcago.NewPermissionDenied(DepositCollection)
	}
	return
}

func (i *DepositParam) DepositSyncPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(DepositCollection)
	}
	return
}

func DepositPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(DepositUnitCollection, "_id", "deposit_id", "deposit_units")
	pipe.LookupUnwind(TakingDepositView, "deposit_units.taking_id", "_id", "deposit_units.taking")
	pipe.Append(bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"}, {Key: "deposit_units", Value: bson.D{
				{Key: "$push", Value: "$deposit_units"},
			}},
		}},
	})
	pipe.LookupUnwind(DepositCollection, "_id", "_id", "deposits")
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "deposits.deposit_units", Value: "$deposit_units"}}}})
	pipe.Append(bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$deposits"}}}})
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	pipe.Lookup(ReceiptFileCollection, "_id", "deposit_id", "receipts")
	return pipe
}

func DepositPipelineList() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(DepositUnitCollection, "_id", "deposit_id", "deposit_units")
	pipe.LookupUnwind(TakingDepositView, "deposit_units.taking_id", "_id", "deposit_units.taking")
	pipe.Append(bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"}, {Key: "deposit_units", Value: bson.D{
				{Key: "$push", Value: "$deposit_units"},
			}},
		}},
	})
	pipe.LookupUnwind(DepositCollection, "_id", "_id", "deposits")
	pipe.Append(bson.D{{Key: "$addFields", Value: bson.D{{Key: "deposits.deposit_units", Value: "$deposit_units"}}}})
	pipe.Append(bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$deposits"}}}})
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	return pipe
}
func DepositPipelineCount(filter bson.D) *vmdb.Pipeline {
	pipe := DepositPipelineList()
	pipe.Match(filter)
	pipe.Limit(300, 300)
	pipe.Append(bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil}, {Key: "list_size", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
		}},
	})
	pipe.Append(bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}}}})
	return pipe
}

func Match(id string) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", id)
	return filter.Bson()
}

/*
Should this be used somewhere?
func UpdateWaitTaking(amount int64) bson.D {
	return bson.D{{Key: "$inc", Value: bson.D{{Key: "state.open.amount", Value: -amount}, {Key: "state.wait.amount", Value: amount}}}}
}
*/

func (i *DepositCreate) DepositDatabase(token *AccessToken) (r *DepositDatabase, d []DepositUnit) {
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

func (i *DepositUpdate) DepositDatabase(current *Deposit) (r *DepositUpdate, create []DepositUnit, update []DepositUnitUpdate, delete []DepositUnit) {
	create = []DepositUnit{}
	update = []DepositUnitUpdate{}
	delete = []DepositUnit{}
	var amount int64 = 0
	for _, value_current := range current.DepositUnit {
		contains := false
		for _, value_update := range i.DepositUnit {
			if value_update.ID == value_current.ID {
				contains = true
			}
		}
		if !contains {
			delete = append(delete, value_current)
		}
	}
	for _, value := range i.DepositUnit {
		if value.ID == "" {
			create = append(create, DepositUnit{
				ID:        uuid.NewString(),
				TakingID:  value.TakingID,
				Money:     value.Money,
				DepositID: i.ID,
				Status:    "open",
				Modified:  vmod.NewModified(),
			})
			amount += value.Money.Amount
		} else {
			update = append(update, value)
			amount += value.Money.Amount
		}
	}
	currency := "EUR"
	if len(i.DepositUnit) != 0 {
		currency = i.DepositUnit[0].Money.Currency
	}
	r = i
	r.Money.Amount = amount
	r.Money.Currency = currency
	if current.Status != "wait" && i.Status == "wait" {
		r.DateOfDeposit = time.Now().Unix()
	}
	return
}

func (i *DepositQuery) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.SearchString([]string{"deposit_units.taking.name", "reason_for_payment"}, i.Search)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew_id", token.CrewID)
	} else {
		filter.EqualStringList("crew_id", i.CrewID)
	}
	filter.EqualStringList("status", i.Status)
	filter.EqualBool("has_external", i.HasExternal)
	filter.LikeString("deposit_units.taking.name", i.Name)
	filter.LikeString("reason_for_payment", i.ReasonForPayment)
	return filter.Bson()
}

func (i *DepositParam) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew_id", token.CrewID)
	}
	return filter.Bson()
}

func (i *DepositQuery) Sort() bson.D {
	sort := vmdb.NewSort()
	sort.Add(i.SortField, i.SortDirection)
	return sort.Bson()
}
