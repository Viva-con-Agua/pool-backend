package dao

type Additionals struct {
	ID     string `json:"id" bson:"_id"`
	Role   string `json:"role" bson:"role"`
	Active string `json:"active" bson:"active"`
	NVM    string `json:"nvm" bson:"nvm"`
	Crew   Crew   `json:"crew" bson:"-"`
	UserID string `json:"user_id" bson:"user_id"`
}
