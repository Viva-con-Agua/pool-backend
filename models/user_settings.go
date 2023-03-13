package models

import "github.com/Viva-con-Agua/vcago/vmod"

type (
	Settings struct {
		ID           string        `json:"id" bson:"_id"`
		Notification Notification  `json:"notification" bson:"notification"`
		Modified     vmod.Modified `json:"modified" bson:"modified"`
	}
	Notification struct {
		FinanceCrew []string `json:"finance_crew_id" bson:"finance_crew_id"`
		FinanceAll  bool     `json:"finance_all" bson:"finance_all"`
	}
)
