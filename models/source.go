package models

import (
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	SourceCreate struct {
		Value       string     `json:"value" bson:"value"`
		Description string     `json:"description" bson:"description"`
		HasExternal bool       `json:"has_external" bson:"has_external"`
		TakingID    string     `json:"taking_id" bson:"taking_id"`
		External    External   `json:"external" bson:"external"`
		PaymentType string     `json:"payment_type" bson:"payment_type"`
		Money       vmod.Money `json:"money" bson:"money"`
	}
	Source struct {
		ID          string        `json:"id" bson:"_id"`
		Value       string        `json:"value" bson:"value"`
		Description string        `json:"description" bson:"description"`
		HasExternal bool          `json:"has_external" bson:"has_external"`
		Money       vmod.Money    `json:"money" bson:"money"`
		TakingID    string        `json:"taking_id" bson:"taking_id"`
		External    External      `json:"external" bson:"external"`
		PaymentType string        `json:"payment_type" bson:"payment_type"`
		Modified    vmod.Modified `json:"modified" bson:"modified"`
	}
	SourceUpdate struct {
		ID          string     `json:"id" bson:"_id"`
		Value       string     `json:"value" bson:"value"`
		Description string     `json:"description" bson:"description"`
		HasExternal bool       `json:"has_external" bson:"has_external"`
		Money       vmod.Money `json:"money" bson:"money"`
		TakingID    string     `json:"taking_id" bson:"taking_id"`
		External    External   `json:"external" bson:"external"`
		PaymentType string     `json:"payment_type" bson:"payment_type"`
		UpdateState string     `json:"update_state" bson:"-"`
	}
	SourceQuery struct {
		Value       string `query:"value"`
		UpdatedTo   string `query:"updated_to" qs:"updated_to"`
		UpdatedFrom string `query:"updated_from" qs:"updated_from"`
		CreatedTo   string `query:"created_to" qs:"created_to"`
		CreatedFrom string `query:"created_from" qs:"created_from"`
	}
	SourceParam struct {
		ID string `param:"id"`
	}
	SourceList []Source
	External   struct {
		Organisation     string `json:"organisation" bson:"organisation"`
		ASP              string `json:"asp" bson:"asp"`
		Email            string `json:"email" bson:"email"`
		Address          string `json:"address" bson:"address"`
		Reciept          bool   `json:"reciept" bson:"reciept"`
		Purpose          string `json:"purpose" bson:"purpose"`
		ReasonForPayment string `json:"reason_for_payment" bson:"reason_for_payment"`
	}
)

func (i *SourceList) InsertMany() []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(*i))
	for n, d := range *i {
		interfaceSlice[n] = d
	}
	return interfaceSlice
}

func (i *SourceCreate) Source() *Source {
	return &Source{
		ID:          uuid.NewString(),
		Value:       i.Value,
		Description: i.Description,
		HasExternal: i.HasExternal,
		External:    i.External,
		TakingID:    i.TakingID,
		PaymentType: i.PaymentType,
		Money:       i.Money,
		Modified:    vmod.NewModified(),
	}
}
func (i *SourceUpdate) Source() *Source {
	return &Source{
		ID:          uuid.NewString(),
		Value:       i.Value,
		Description: i.Description,
		HasExternal: i.HasExternal,
		External:    i.External,
		PaymentType: i.PaymentType,
		TakingID:    i.TakingID,
		Money:       i.Money,
		Modified:    vmod.NewModified(),
	}
}

func (i *SourceParam) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *SourceUpdate) Filter() bson.D {
	return bson.D{{Key: "_id", Value: i.ID}}
}

func (i *SourceQuery) Filter() bson.D {
	filter := vmdb.NewFilter()
	filter.LikeString("value", i.Value)
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	return bson.D(*filter)
}
