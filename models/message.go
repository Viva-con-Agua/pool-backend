package models

import (
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	MessageCreate struct {
		From           string         `json:"from" bson:"from"`
		Subject        string         `json:"subject" bson:"subject"`
		Message        string         `json:"message" bson:"message"`
		MailboxID      string         `json:"mailbox_id" bson:"mailbox_id"`
		Read           bool           `json:"read" bson:"read"`
		RecipientGroup RecipientGroup `json:"recipient_group" bson:"recipient_group"`
	}
	MessageUpdate struct {
		ID             string         `json:"id" bson:"_id"`
		MessageID      string         `json:"message_id" bson:"message_id"`
		From           string         `json:"from" bson:"from"`
		Subject        string         `json:"subject" bson:"subject"`
		Message        string         `json:"message" bson:"message"`
		UserID         string         `json:"user_id" bson:"user_id"`
		Type           string         `json:"type" bson:"type"`
		MailboxID      string         `json:"mailbox_id" bson:"mailbox_id"`
		Read           bool           `json:"read" bson:"read"`
		RecipientGroup RecipientGroup `json:"recipient_group" bson:"recipient_group"`
		To             []User         `json:"to" bson:"to"`
	}
	MessageParam struct {
		ID string `param:"id"`
	}
	Message struct {
		ID             string         `json:"id" bson:"_id"`
		MessageID      string         `json:"message_id" bson:"message_id"`
		From           string         `json:"from" bson:"from"`
		Subject        string         `json:"subject" bson:"subject"`
		Message        string         `json:"message" bson:"message"`
		UserID         string         `json:"user_id" bson:"user_id"`
		Type           string         `json:"type" bson:"type"`
		MailboxID      string         `json:"mailbox_id" bson:"mailbox_id"`
		Read           bool           `json:"read" bson:"read"`
		RecipientGroup RecipientGroup `json:"recipient_group" bson:"recipient_group"`
		To             []User         `json:"to" bson:"to"`
		Modified       vmod.Modified  `json:"modified" bson:"modified"`
	}
	MessageList []Message

	MessageSubList []MessageSub

	MessageSub struct {
		ID             string         `json:"id" bson:"_id"`
		MessageID      string         `json:"message_id" bson:"message_id"`
		From           string         `json:"from" bson:"from"`
		Subject        string         `json:"subject" bson:"subject"`
		Message        string         `json:"message" bson:"message"`
		UserID         string         `json:"user_id" bson:"user_id"`
		Type           string         `json:"type" bson:"type"`
		MailboxID      string         `json:"mailbox_id" bson:"mailbox_id"`
		Read           bool           `json:"read" bson:"read"`
		RecipientGroup RecipientGroup `json:"recipient_group" bson:"recipient_group"`
		To             []User         `json:"-" bson:"to"`
		Modified       vmod.Modified  `json:"modified" bson:"modified"`
	}

	MessageQuery struct {
		ID          string   `query:"id"`
		From        []string `query:"from"`
		Subject     string   `query:"subject"`
		UserID      string   `query:"user_id"`
		UpdatedTo   string   `query:"updated_to"`
		UpdatedFrom string   `query:"updated_from"`
		CreatedTo   string   `query:"created_to"`
		CreatedFrom string   `query:"created_from"`
	}
)

var MessageCollection = "messages"

func (i *MessageCreate) MessageSub(token *vcapool.AccessToken) *Message {
	return &Message{
		ID:             uuid.NewString(),
		MessageID:      uuid.NewString(),
		From:           i.From,
		Subject:        i.Subject,
		Message:        i.Message,
		MailboxID:      i.MailboxID,
		Read:           i.Read,
		RecipientGroup: i.RecipientGroup,
		Modified:       vmod.NewModified(),
		Type:           "draft",
		UserID:         token.ID,
	}
}

func (i *Message) MessageUpdate() *MessageUpdate {
	return &MessageUpdate{
		ID:             i.ID,
		MessageID:      i.MessageID,
		From:           i.From,
		Subject:        i.Subject,
		Message:        i.Message,
		UserID:         i.UserID,
		Type:           i.Type,
		MailboxID:      i.MailboxID,
		Read:           i.Read,
		RecipientGroup: i.RecipientGroup,
		To:             i.To,
	}
}

func (i *MessageParam) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{
		{Key: "_id", Value: i.ID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "mailbox_id", Value: token.ID}},
			bson.D{{Key: "mailbox_id", Value: token.CrewID}},
		}}}
}

func (i *MessageUpdate) Filter(token *vcapool.AccessToken) bson.D {
	return bson.D{
		{Key: "_id", Value: i.ID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "mailbox_id", Value: token.ID}},
			bson.D{{Key: "mailbox_id", Value: token.CrewID}},
		}}}
}

func (i *Message) Inbox() *[]interface{} {
	inbox := new([]interface{})
	for n := range i.To {
		message := *i
		message.ID = uuid.NewString()
		message.MailboxID = (i.To)[n].ID
		message.Type = "inbox"
		message.Modified = vmod.NewModified()
		*inbox = append(*inbox, message)

	}
	return inbox

}
