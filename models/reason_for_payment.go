package models

type ReasonForPayment struct {
	ID           string `bson:"_id"`
	CrewID       string `bson:"crew_id"`
	Abbreviation string `bson:"abbreviation"`
	Count        int    `bson:"count"`
	Year         string `bson:"year"`
}
