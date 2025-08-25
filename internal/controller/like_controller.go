package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"actiondelta/internal/model"
	"actiondelta/internal/repository"
)

// ToggleLike 点赞/取消点赞（幂等切换）
// 支持 target_type: record/backstory
func ToggleLike(c *gin.Context) {
	userId := c.GetString("userId")
	var body struct {
		TargetType string `json:"target_type"`
		TargetId   string `json:"target_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.TargetType == "" || body.TargetId == "" {
		respond(c, http.StatusBadRequest, "invalid request", nil)
		return
	}
	tid, err := primitive.ObjectIDFromHex(body.TargetId)
	if err != nil { respond(c, http.StatusBadRequest, "invalid target_id", nil); return }

	filter := bson.M{"userId": userId, "targetType": body.TargetType, "targetId": tid}
	likes := repository.DB().Collection("likes")
	// 是否已点赞
	cnt, _ := likes.CountDocuments(c, filter)
	liked := false
	if cnt > 0 {
		// 取消点赞
		_, _ = likes.DeleteOne(c, filter)
		liked = false
		updateLikeCounter(c, body.TargetType, tid, -1)
	} else {
		// 新增点赞
		_, _ = likes.InsertOne(c, model.Like{UserId: userId, TargetType: body.TargetType, TargetId: tid, CreatedAt: time.Now()})
		liked = true
		updateLikeCounter(c, body.TargetType, tid, +1)
	}
	respond(c, http.StatusOK, "success", gin.H{"liked": liked})
}

// updateLikeCounter 更新目标对象的点赞数（简单累加）
func updateLikeCounter(c *gin.Context, targetType string, targetId primitive.ObjectID, delta int) {
	switch targetType {
	case "record":
		_, _ = repository.DB().Collection("cassettes").UpdateByID(c, targetId, bson.M{"$inc": bson.M{"likeCount": delta}})
	case "backstory":
		_, _ = repository.DB().Collection("backstories").UpdateByID(c, targetId, bson.M{"$inc": bson.M{"likeCount": delta}})
	}
} 