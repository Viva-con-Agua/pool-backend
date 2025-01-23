package models

import (
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
)

type (
	ReceiptFileCreate struct {
		DepositID string `bson:"deposit_id" json:"deposit_id"`
	}

	ReceiptFile struct {
		ID        string        `bson:"_id" json:"id"`
		DepositID string        `bson:"deposit_id" json:"deposit_id"`
		Modified  vmod.Modified `bson:"modified" json:"modified"`
	}
)

var (
	ReceiptFileCollection = "receipts_files"
)

func (i *ReceiptFileCreate) ReceiptFile() *ReceiptFile {
	return &ReceiptFile{
		ID:        uuid.NewString(),
		DepositID: i.DepositID,
		Modified:  vmod.NewModified(),
	}
}
