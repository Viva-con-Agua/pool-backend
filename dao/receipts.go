package dao

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"pool-backend/models"
	"strconv"

	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

func ReceiptFileCreate(
	ctx context.Context,
	create *models.ReceiptFileCreate,
	file *vmod.File,
	token *models.AccessToken,
) (
	result *models.ReceiptFile,
	err error,
) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	deposit := new(models.Deposit)
	filter := bson.D{{Key: "_id", Value: create.DepositID}}
	if err = DepositCollection.FindOne(ctx, filter, deposit); err != nil {
		return
	}
	//permission check
	result = create.ReceiptFile()
	if err = ReceiptFileCollection.InsertOne(ctx, result); err != nil {
		return
	}
	if err = Database.UploadFile(file, result.ID); err != nil {
		return
	}
	return
}

func ReceiptFileGetByID(
	ctx context.Context,
	id *vmod.IDParam,
	token *models.AccessToken,
) (
	result []byte,
	err error,
) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	file := new(models.ReceiptFile)
	if err = ReceiptFileCollection.FindOne(ctx, id.Filter(), file); err != nil {
		return
	}
	//permission check
	if result, err = Database.DownloadFile(id.ID); err != nil {
		return
	}
	return
}

func ReceiptFileZipGetByID(
	ctx context.Context,
	id *vmod.IDParam,
	token *models.AccessToken,
) (
	result []byte,
	err error,
) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	deposit := new(models.Deposit)
	if err = DepositCollection.AggregateOne(ctx, models.DepositPipeline().Match(id.Filter()).Pipe, deposit); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for index, file := range deposit.Receipts {
		var bbuffer []byte
		if bbuffer, err = Database.DownloadFile(file.ID); err != nil {
			return
		}
		var f io.Writer
		if f, err = w.Create(deposit.ReasonForPayment + "_" + strconv.Itoa(index+1) + ".png"); err != nil {
			return
		}
		if _, err = f.Write(bbuffer); err != nil {
			return
		}
	}
	err = w.Close()
	result = buf.Bytes()
	return
}

func ReceiptFileDeleteByID(
	ctx context.Context,
	id *vmod.IDParam,
	token *models.AccessToken,
) (
	result *vmod.DeletedResponse,
	err error,
) {
	if err = models.DepositPermission(token); err != nil {
		return
	}
	file := new(models.ReceiptFile)
	if err = ReceiptFileCollection.FindOne(ctx, id.Filter(), file); err != nil {
		return
	}
	//permission check
	if err = Database.DeleteFile(ctx, id.ID); err != nil {
		return
	}
	err = ReceiptFileCollection.DeleteOne(ctx, id.Filter())
	result = vmod.NewDeletedResponse(id.ID)
	return
}
