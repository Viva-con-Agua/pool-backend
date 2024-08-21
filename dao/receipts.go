package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
	"go.mongodb.org/mongo-driver/bson"
)

func ReceiptFileCreate(
	ctx context.Context,
	create *models.ReceiptFileCreate,
	file *vmod.File,
	token *vcapool.AccessToken,
) (
	result *models.ReceiptFile,
	err error,
) {
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
	token *vcapool.AccessToken,
) (
	result []byte,
	err error,
) {
	filter := bson.D{{Key: "_id", Value: id.ID}}
	file := new(models.ReceiptFile)
	if err = ReceiptFileCollection.FindOne(ctx, filter, file); err != nil {
		return
	}
	//permission check
	if result, err = Database.DownloadFile(id.ID); err != nil {
		return
	}
	return
}

func ReceiptFileDeleteByID(
	ctx context.Context,
	id *vmod.IDParam,
	token *vcapool.AccessToken,
) (
	result *vmod.DeletedResponse,
	err error,
) {
	filter := bson.D{{Key: "_id", Value: id.ID}}
	file := new(models.ReceiptFile)
	if err = ReceiptFileCollection.FindOne(ctx, filter, file); err != nil {
		return
	}
	//permission check
	if err = Database.DeleteFile(ctx, id.ID); err != nil {
		return
	}
	err = ReceiptFileCollection.DeleteOne(ctx, filter)
	result = vmod.NewDeletedResponse(id.ID)
	return
}
