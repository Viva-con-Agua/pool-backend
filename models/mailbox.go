package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func MailboxPipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwind(UserCollection, "_id", "mailbox_id", "user")
	pipe.LookupUnwind(CrewCollection, "_id", "mailbox_id", "crew")
	inboxMatch := vmdb.NewFilter()
	inboxMatch.EqualString("type", "inbox")
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "inbox", inboxMatch.Bson())
	outboxMatch := vmdb.NewFilter()
	outboxMatch.EqualString("type", "outbox")
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "outbox", outboxMatch.Bson())
	draftMatch := vmdb.NewFilter()
	draftMatch.EqualString("type", "draft")
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "draft", draftMatch.Bson())
	return pipe
}

func (i *MailboxParam) Permission(token *vcapool.AccessToken, mailbox *Mailbox) (err error) {
	if token.MailboxID == mailbox.ID {
		return
	} else if token.CrewID == mailbox.ID {
		return
	}
	return vcago.NewPermissionDenied(MailboxCollection)
}

func (i *MailboxParam) Pipeline() mongo.Pipeline {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	pipe := MailboxPipeline()
	pipe.Match(filter.Bson())
	return pipe.Pipe
}

func (i *MailboxParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("operation;network;finance;education;socialmedia;awareness;asp;")) {
		filter.EqualString("user._id", token.ID)
		filter.EqualString("type", "user")
	} else if !token.Roles.Validate("employee;admin") {
		status := bson.A{}
		status = append(status, bson.D{{Key: "user._id", Value: token.ID}})
		status = append(status, bson.D{{Key: "crew._id", Value: token.CrewID}})
		filter.Append(bson.E{Key: "$or", Value: status})
	}
	return filter.Bson()
}
