package models

import (
	"encoding/base64"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
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

func MessageCrewPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		return vcago.NewPermissionDenied(MessageCollection)
	}
	return
}

func MessageEventPermission(token *AccessToken, event *Event) (err error) {
	if !token.PoolRoles.Validate(ASPRole) && !token.Roles.Validate("admin;employee;pool_employee") && event.EventASPID != token.ID {
		return vcago.NewPermissionDenied(MessageCollection)
	}
	return
}

func (i *MessageCreate) MessageSub(token *AccessToken) *Message {
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

func (i *Message) NotificationMessage(response *vcago.NotificationResponse, user *User) *Message {
	return &Message{
		ID:             uuid.NewString(),
		From:           response.From,
		Subject:        response.Subject,
		Message:        base64.StdEncoding.EncodeToString([]byte(response.Body)),
		MailboxID:      user.MailboxID,
		Read:           true,
		RecipientGroup: i.RecipientGroup,
		Modified:       vmod.NewModified(),
		Type:           "inbox",
		UserID:         response.User.ID,
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

func (i *MessageParam) PermittedFilter(token *AccessToken, crew *Crew) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
		filter.EqualString("user_id", token.ID)
	} else {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
	}
	return filter.Bson()
}

func (i *MessageUpdate) PermittedFilter(token *AccessToken, crew *Crew) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
		filter.EqualString("user_id", token.ID)
	} else {
		filter.EqualStringList("mailbox_id", []string{token.MailboxID, crew.MailboxID})
	}
	return filter.Bson()
}

func (i *RecipientGroup) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !token.Roles.Validate("admin;employee;pool_employee") {
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

func PermittedMessageCreate(token *AccessToken, i *Message, crew *Crew, event *Event) (message *Message, err error) {
	message = i
	if !token.Roles.Validate("admin;employee;pool_employee") {
		if i.RecipientGroup.Type == "event" && token.ID == event.EventASPID {
			// IF IS EVENT ASP -> Force Mailbox and From to CrewMailbox and CrewEmail
			message.MailboxID = crew.MailboxID
			message.From = crew.Email
		} else if token.PoolRoles.Validate(ASPEventRole) {
			// IF IS CREW ASP -> Force Mailbox and From to CrewMailbox and CrewEmail
			message.MailboxID = crew.MailboxID
			message.From = crew.Email
		} else {
			return nil, vcago.NewBadRequest(MessageCollection, "Not allwed to create a message")
			// i.MailboxID = token.MailboxID
			// i.From = token.Email
		}
	}
	return
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
