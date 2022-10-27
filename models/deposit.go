package models

import "github.com/Viva-con-Agua/vcago/vmod"

type (
	DepositUnit struct {
		ID        string        `json:"id" bson:"_id" bson:"sources_ids"`
		TakingID  string        `json:"taking_id" bson:"taking_id"`
		Money     vmod.Money    `json:"money" bson:"money"`
		DepositID string        `json:"deposit_id" bson:"deposit_id"`
		Status    string        `json:"status" bson:"status"`
		Modified  vmod.Modified `json:"modified" bson:"modified"`
	}
	DepositDatabase struct {
		ID               string        `json:"id" bson:"_id" bson:"sources_ids"`
		ReasonForPayment string        `json:"reason_for_payment" bson:"sources_ids"`
		Status           string        `json:"status" bson:"status"`
		DepositUnitID    []string      `json:"deposit_unit_ids" bson:"deposit_unit_ids"`
		Money            vmod.Money    `json:"money" bson:"money"`
		Modified         vmod.Modified `json:"modified" bson:"modified"`
	}
	Deposit struct {
		ID               string        `json:"id" bson:"_id" bson:"sources_ids"`
		ReasonForPayment string        `json:"reason_for_payment" bson:"sources_ids"`
		Status           string        `json:"status" bson:"status"`
		DepositUnitID    []string      `json:"deposit_unit_ids" bson:"deposit_unit_ids"`
		DepositUnit      []DepositUnit `json:"deposit_units" bson:"deposit_units"`
		Money            vmod.Money    `json:"money" bson:"money"`
		Modified         vmod.Modified `json:"modified" bson:"modified"`
	}
)
