package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
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
		ToID           []string       `json:"to_id" bson:"to_id"`
		To             []TOData       `json:"to" bson:"-"`
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
		To             []TOData       `json:"-" bson:"to"`
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
	TOData struct {
		UserID    string `bson:"user_id" json:"user_id"`
		MailboxID string `bson:"mailbox_id" json:"mailbox_id"`
		Email     string `bson:"email" json:"email"`
	}
)

var MessageCollection = "messages"

func MessageCrewPermission(token *vcapool.AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(MessageCollection)
	}
	return
}

func MessageEventPermission(token *vcapool.AccessToken, event *Event) (err error) {
	if !token.PoolRoles.Validate(ASPRole) && !token.Roles.Validate("admin;employee") && event.EventASPID != token.ID {
		return vcago.NewPermissionDenied(MessageCollection)
	}
	return
}

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
	}
}

func (i *MessageParam) PermittedFilter(token *vcapool.AccessToken, crew *Crew) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		filter.EqualString("mailbox_id", token.MailboxID)
	} else {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
	}
	return filter.Bson()
}

func (i *MessageUpdate) PermittedFilter(token *vcapool.AccessToken, crew *Crew) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		filter.EqualString("mailbox_id", token.MailboxID)
	} else {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
	}
	return filter.Bson()
}

func (i *RecipientGroup) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee") {
		filter.EqualString("crew.crew_id", token.CrewID)
	} else {
		filter.EqualString("crew.crew_id", i.CrewID)
	}
	filter.EqualStringList("active.status", i.Active)
	filter.EqualStringList("nvm.status", i.NVM)
	if !i.IgnoreNewsletter {
		filter.EqualString("newsletter.value", "regional")
	}
	return filter.Bson()
}

func (i *RecipientGroup) FilterEvent() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.EventID)
	return filter.Bson()
}

func (i *Message) PermittedCreate(token *vcapool.AccessToken, crew *Crew, event *Event) *Message {
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		if !(i.RecipientGroup.Type == "event" || token.ID == event.EventASPID) { // USER
			i.MailboxID = token.MailboxID
			i.From = token.Email
		} else { // EVENT ASP
			i.MailboxID = crew.MailboxID
			if !(i.From == token.Email || i.From == crew.Email) {
				i.From = token.Email
			}
		}
	} else if !(token.Roles.Validate("employee;admin")) { // ASP
		if i.MailboxID == crew.MailboxID {
			i.RecipientGroup.CrewID = crew.ID
			if !(i.From == token.Email || i.From == crew.Email) {
				i.From = token.Email
			}
		} else {
			i.MailboxID = token.MailboxID
			i.From = token.Email
		}
	}
	// ADMIN
	return i
}

func (i *Message) Inbox() *[]interface{} {
	inbox := new([]interface{})
	for n := range i.To {
		message := *i
		message.ID = uuid.NewString()
		message.MailboxID = (i.To)[n].MailboxID
		message.Type = "inbox"
		message.Modified = vmod.NewModified()
		*inbox = append(*inbox, message)
	}
	return inbox
}

func (i *Message) ToEmails() (result []string) {
	result = []string{}
	for _, value := range i.To {
		result = append(result, value.Email)
	}
	return
}
