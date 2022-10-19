package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
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
	inboxMatch := vmdb.NewFilter()
	inboxMatch.EqualString("type", "inbox")
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "inbox", inboxMatch.Bson())
	outboxMatch := vmdb.NewFilter()
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "outbox", outboxMatch.Bson())
	draftMatch := vmdb.NewFilter()
	pipe.LookupMatch(MessageCollection, "_id", "mailbox_id", "draft", draftMatch.Bson())
	return pipe
}

func (i *MailboxParam) Permission(token *vcapool.AccessToken) (err error) {
	if token.ID == i.ID {
		return
	} else if token.CrewID == i.ID {
		return
	}
	return vcago.NewPermissionDenied("mailbox", i)
}

func (i *MailboxParam) Pipeline() mongo.Pipeline {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	pipe := MailboxPipeline()
	pipe.Match(match.Bson())
	return pipe.Pipe
}
