package token

import (
	"bytes"
	"io"
	"net/http"
	"pool-backend/dao"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type AvatarHandler struct {
	vcago.Handler
}

var Avatar = &AvatarHandler{*vcago.NewHandler("avatar")}

func (i *AvatarHandler) Routes(group *echo.Group) {
	group.Use(i.Context)
	group.DELETE("/:id", i.Delete, accessCookie)
	group.POST("/upload", i.Upload, accessCookie)
	group.GET("/img/:id", i.GetByID, accessCookie)
	group.DELETE("/img/:id", i.Delete, accessCookie)
}

// Upload
// @Security CookieAuth
// @Summary Uploads a Avatar
// @Description creates an  Avatar object.
// @Tags Avatar
// @Accept json
// @Produce json
// @Param form body models.AvatarFile true "Avatar File"
// @Model: vcago.Response
// @Success 201 {object} vcago.Response{payload=models.Avatar}
// @Router /users/avatar [post]
func (i *AvatarHandler) Upload(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	result := models.NewAvatar(token)
	filter := bson.D{{Key: "_id", Value: token.ID}}
	update := bson.D{{Key: "avatar", Value: result}}
	if err = dao.UserCollection.UpdateOne(c.Ctx(), filter, vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	file := new(models.AvatarFile)
	if file.File, file.Header, err = c.Request().FormFile("file"); err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file.File); err != nil {
		return
	}
	bucket := new(gridfs.Bucket)
	if bucket, err = gridfs.NewBucket(dao.Database.Database); err != nil {
		return
	}
	uploadStream := new(gridfs.UploadStream)
	if uploadStream, err = bucket.OpenUploadStreamWithID(result.ID, file.Header.Filename); err != nil {
		return
	}
	defer uploadStream.Close()
	if _, err = uploadStream.Write(buf.Bytes()); err != nil {
		return
	}

	return c.Created(result)
}

// GetByID
// @Security CookieAuth
// @Summary Get a  Avatar by ID
// @Tags Avatar
// @Accept json
// @Produce image/png
// @Param id path string true "Avatar ID"
// @Success 200 {file} string "Download Image"
// @Router /users/avatar/{id} [get]
func (i *AvatarHandler) GetByID(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AvatarParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	var buf bytes.Buffer
	bucket, _ := gridfs.NewBucket(dao.Database.Database)
	if _, err = bucket.DownloadToStream(body.ID, &buf); err != nil {
		return
	}
	return c.Stream(http.StatusOK, "image/png", bytes.NewReader(buf.Bytes()))
}

// DeleteByID
// @Security CookieAuth
// @Summary Get a  Avatar by ID
// @Tags Avatar
// @Accept json
// @Produce json
// @Param id path string true "Avatar ID"
// @Success 200 {object} vmod.DeletedResponse
// @Router /users/avatar/{id} [delete]
func (i *AvatarHandler) Delete(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(models.AvatarParam)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	token := new(models.AccessToken)
	if err = c.AccessToken(token); err != nil {
		return
	}
	update := bson.D{{Key: "avatar", Value: models.NewAvatarClean()}}
	if err = dao.UserCollection.UpdateOne(c.Ctx(), body.PermittedFilter(token), vmdb.UpdateSet(update), nil); err != nil {
		return
	}
	if err = dao.FSChunkCollection.DeleteOne(c.Ctx(), body.MatchChunk()); err != nil {
		return
	}
	if err = dao.FSFilesCollection.DeleteOne(c.Ctx(), body.Match()); err != nil {
		return
	}
	return c.Deleted(body.ID)
}
