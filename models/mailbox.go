package models

import (
	"fmt"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	MailboxDatabase struct {
		ID   string `json:"id" bson:"_id"`
		Type string `json:"type" bson:"type"`
	}
	Mailbox struct {
		ID     string         `json:"id" bson:"_id"`
		Type   string         `json:"type" bson:"type"`
		Inbox  MessageSubList `json:"inbox" bson:"inbox"`
		Outbox MessageList    `json:"outbox" bson:"outbox"`
		Draft  MessageSubList `json:"draft" bson:"draft"`
	}
	MailboxParam struct {
		ID string `param:"id"`
	}
)

var MailboxCollection = "mailbox"

func NewMailboxDatabase(t string) *MailboxDatabase {
	return &MailboxDatabase{
		ID:   uuid.NewString(),
		Type: t,
	}
}

func MailboxPipeline(token *AccessToken) *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(UserCollection, "_id", "mailbox_id", "user")
	pipe.LookupUnwind(CrewCollection, "_id", "mailbox_id", "crew")
	lastSixMonth := fmt.Sprint(time.Now().AddDate(0, -6, 0).Unix())
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPRole)) {
		inboxMatch := vmdb.NewFilter()
		inboxMatch.EqualString("type", "inbox")
		inboxMatch.GteInt64("modified.updated", lastSixMonth)
		inboxMatch.EqualString("user_id", token.ID)
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "inbox", inboxMatch.Bson())
		outboxMatch := vmdb.NewFilter()
		outboxMatch.EqualString("type", "outbox")
		outboxMatch.GteInt64("modified.updated", lastSixMonth)
		outboxMatch.EqualString("user_id", token.ID)
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "outbox", outboxMatch.Bson())
		draftMatch := vmdb.NewFilter()
		draftMatch.EqualString("type", "draft")
		draftMatch.GteInt64("modified.updated", lastSixMonth)
		draftMatch.EqualString("user_id", token.ID)
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "draft", draftMatch.Bson())
	} else {
		inboxMatch := vmdb.NewFilter()
		inboxMatch.GteInt64("modified.updated", lastSixMonth)
		inboxMatch.EqualString("type", "inbox")
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "inbox", inboxMatch.Bson())
		outboxMatch := vmdb.NewFilter()
		outboxMatch.GteInt64("modified.updated", lastSixMonth)
		outboxMatch.EqualString("type", "outbox")
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "outbox", outboxMatch.Bson())
		draftMatch := vmdb.NewFilter()
		draftMatch.GteInt64("modified.updated", lastSixMonth)
		draftMatch.EqualString("type", "draft")
		pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "draft", draftMatch.Bson())
	}
	return pipe
}

func (i *MailboxParam) Permission(token *AccessToken, mailbox *Mailbox) (err error) {
	if token.MailboxID == mailbox.ID {
		return
	} else if token.CrewID == mailbox.ID {
		return
	}
	return vcago.NewPermissionDenied(MailboxCollection)
}

func (i *MailboxParam) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		status := bson.A{}
		status = append(status, bson.D{{Key: "user._id", Value: token.ID}})
		status = append(status, bson.D{{Key: "crew._id", Value: token.CrewID}})
		filter.Append(bson.E{Key: "$or", Value: status})
	}
	return filter.Bson()
}
