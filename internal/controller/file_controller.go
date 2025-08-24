package controller

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"actiondelta/internal/repository"
)

const (
	maxAvatarSizeBytes = 5 * 1024 * 1024
)

// UploadAvatar 头像上传接口：校验、裁剪、压缩并存入 Mongo GridFS。
func UploadAvatar(c *gin.Context) {
	userId := c.GetString("userId")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respond(c, http.StatusBadRequest, "文件缺失", nil)
		return
	}
	if header.Size <= 0 || header.Size > maxAvatarSizeBytes {
		respond(c, http.StatusBadRequest, "文件大小不合法", nil)
		return
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		respond(c, http.StatusBadRequest, "读取文件失败", nil)
		return
	}
	data := buf.Bytes()
	if err := validateImageFile(data); err != nil {
		respond(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 解码
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		respond(c, http.StatusBadRequest, "图片解码失败", nil)
		return
	}
	// 裁剪为正方形，限制最大尺寸
	cropped := cropSquare(img)
	resized := imaging.Fit(cropped, 1024, 1024, imaging.Lanczos)
	thumb := imaging.Resize(cropped, 200, 200, imaging.Lanczos)

	// 编码JPEG
	var outFull, outThumb bytes.Buffer
	if err := jpeg.Encode(&outFull, resized, &jpeg.Options{Quality: 85}); err != nil {
		respond(c, http.StatusInternalServerError, "图片编码失败", nil)
		return
	}
	if err := jpeg.Encode(&outThumb, thumb, &jpeg.Options{Quality: 80}); err != nil {
		respond(c, http.StatusInternalServerError, "缩略图编码失败", nil)
		return
	}

	// 保存到 GridFS
	fullName := "avatar-" + uuid.NewString() + ".jpg"
	thumbName := "avatar-thumb-" + uuid.NewString() + ".jpg"
	fullID, err := saveToGridFS(c, outFull.Bytes(), fullName)
	if err != nil {
		respond(c, http.StatusInternalServerError, "保存失败", nil)
		return
	}
	thumbID, err := saveToGridFS(c, outThumb.Bytes(), thumbName)
	if err != nil {
		respond(c, http.StatusInternalServerError, "保存失败", nil)
		return
	}

	// 更新用户头像（保存可直接访问的 API URL）
	avatarURL := "/api/file/" + fullID.Hex()
	thumbURL := "/api/file/" + thumbID.Hex()
	if err := updateUserAvatar(c, userId, avatarURL, thumbURL); err != nil {
		// 不阻断返回
	}

	respond(c, http.StatusOK, "上传成功", gin.H{
		"file_id":       fullID.Hex(),
		"thumb_file_id": thumbID.Hex(),
		"avatar_url":    avatarURL,
		"thumbnail_url": thumbURL,
		"uploaded_at":   time.Now().UTC(),
	})
}

// GetFile 按文件ID从 GridFS 流式输出（当前固定 image/jpeg）。
func GetFile(c *gin.Context) {
	idHex := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid id", nil)
		return
	}
	bucket, err := repository.GridFS()
	if err != nil {
		respond(c, http.StatusInternalServerError, "storage error", nil)
		return
	}
	stream, err := bucket.OpenDownloadStream(oid)
	if err != nil {
		respond(c, http.StatusNotFound, "file not found", nil)
		return
	}
	defer stream.Close()
	c.Header("Content-Type", "image/jpeg")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, stream); err != nil {
		// 传输中断无需额外处理
	}
}

func cropSquare(img image.Image) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	size := w
	if h < w {
		size = h
	}
	return imaging.CropCenter(img, size, size)
}

func validateImageFile(data []byte) error {
	if len(data) < 8 {
		return errors.New("文件过小或损坏")
	}
	// 简单魔数校验（JPEG/PNG/GIF）
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8}) { // JPEG
		return nil
	}
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) { // PNG
		return nil
	}
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46}) { // GIF
		return nil
	}
	return errors.New("不支持的文件格式")
}

// saveToGridFS 将数据以给定文件名写入 GridFS，返回文件ID。
func saveToGridFS(c *gin.Context, data []byte, filename string) (primitive.ObjectID, error) {
	bucket, err := repository.GridFS()
	if err != nil {
		return primitive.NilObjectID, err
	}
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer uploadStream.Close()
	if _, err := uploadStream.Write(data); err != nil {
		return primitive.NilObjectID, err
	}
	if oid, ok := uploadStream.FileID.(primitive.ObjectID); ok {
		return oid, nil
	}
	return primitive.NilObjectID, errors.New("invalid file id")
}

// updateUserAvatar 更新用户头像信息
func updateUserAvatar(c *gin.Context, userId, avatarURL, thumbnailURL string) error {
	update := bson.M{
		"$set": bson.M{
			"avatar":    avatarURL,
			"thumbnail": thumbnailURL,
			"updatedAt": time.Now(),
		},
	}
	_, err := repository.DB().Collection("users").UpdateOne(
		c,
		bson.M{"userId": userId},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}
