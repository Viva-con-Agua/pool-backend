package models

// RecipientGroup represents
type RecipientGroup struct {
	Type             string   `json:"type" bson:"type"`                 // can be crew or event
	CrewID           string   `json:"crew_id" bson:"crew_id"`           // only used for type crew
	EventID          string   `json:"event_id" bson:"event_id"`         // only used for type event
	Active           []string `json:"active_state" bson:"active_state"` //only used for type crew
	NVM              []string `json:"nvm_state" bson:"nvm_state"`       //only used for type crew
	State            []string `json:"state" bson:"state"`               // only used for type event
	IgnoreNewsletter bool     `json:"ignore_newsletter" bson:"ignore_newsletter"`
}

func (i *RecipientGroup) UserQuery() *UserQuery {
	query := new(UserQuery)
	query.CrewID = i.CrewID
	query.ActiveState = i.Active
	query.NVMState = i.NVM
	return query
}
