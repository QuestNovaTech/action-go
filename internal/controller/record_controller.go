package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"actiondelta/internal/model"
	"actiondelta/internal/repository"
)

// CreateRecord 生成戏文：从会话/房间中选择若干消息片段生成
func CreateRecord(c *gin.Context) {
	userId := c.GetString("userId")
	var body struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		BackstoryId  string   `json:"backstory_id"`
		RoomId       string   `json:"room_id"`
		MessageIds   []string `json:"message_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.MessageIds) == 0 {
		respond(c, http.StatusBadRequest, "invalid request", nil)
		return
	}
	// 转换消息ID
	msgOids := make([]primitive.ObjectID, 0, len(body.MessageIds))
	for _, id := range body.MessageIds {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil { msgOids = append(msgOids, oid) }
	}
	if len(msgOids) == 0 { respond(c, http.StatusBadRequest, "invalid message_ids", nil); return }
	// 查询消息，推导参与者
	cur, err := repository.DB().Collection("messages").Find(c, bson.M{"_id": bson.M{"$in": msgOids}})
	if err != nil { respond(c, http.StatusInternalServerError, "server error", nil); return }
	var msgs []model.Message
	_ = cur.All(c, &msgs)
	parts := make(map[string]model.CassetteParticipant)
	for _, m := range msgs {
		cp := model.CassetteParticipant{UserId: m.SenderUserId}
		if m.CharacterInfo != nil { cp.CharacterId = m.CharacterInfo.CharacterId }
		parts[cp.UserId+"/"+cp.CharacterId] = cp
	}
	participants := make([]model.CassetteParticipant, 0, len(parts))
	for _, p := range parts { participants = append(participants, p) }

	var backstoryOID *primitive.ObjectID
	if body.BackstoryId != "" {
		if oid, err := primitive.ObjectIDFromHex(body.BackstoryId); err == nil { backstoryOID = &oid }
	}
	var roomOID *primitive.ObjectID
	if body.RoomId != "" {
		if oid, err := primitive.ObjectIDFromHex(body.RoomId); err == nil { roomOID = &oid }
	}
	rec := model.Cassette{
		Title:        body.Title,
		Description:  body.Description,
		BackstoryId:  backstoryOID,
		RoomId:       roomOID,
		CreatorId:    userId,
		Participants: participants,
		MessageIds:   msgOids,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	res, err := repository.DB().Collection("cassettes").InsertOne(c, rec)
	if err != nil { respond(c, http.StatusInternalServerError, "server error", nil); return }
	id := res.InsertedID.(primitive.ObjectID)
	respond(c, http.StatusOK, "success", gin.H{"id": id.Hex()})
}

// ListRecords 戏文列表（分页/关键字）
func ListRecords(c *gin.Context) {
	page := parseIntDefault(c.DefaultQuery("page", "1"), 1)
	size := parseIntDefault(c.DefaultQuery("size", "20"), 20)
	keyword := c.Query("keyword")
	filter := bson.M{}
	if keyword != "" { filter["title"] = bson.M{"$regex": keyword, "$options": "i"} }
	col := repository.DB().Collection("cassettes")
	total, _ := col.CountDocuments(c, filter)
	cur, err := col.Find(c, filter, options.Find().SetSort(bson.M{"createdAt": -1}).SetSkip(int64((page-1)*size)).SetLimit(int64(size)))
	if err != nil { respond(c, http.StatusInternalServerError, "查询失败", nil); return }
	var list []model.Cassette
	_ = cur.All(c, &list)
	respond(c, http.StatusOK, "success", gin.H{"total": total, "list": list})
}

// GetRecord 戏文详情
func GetRecord(c *gin.Context) {
	idHex := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { respond(c, http.StatusBadRequest, "invalid id", nil); return }
	var r model.Cassette
	if err := repository.DB().Collection("cassettes").FindOne(c, bson.M{"_id": oid}).Decode(&r); err != nil {
		respond(c, http.StatusNotFound, "not found", nil)
		return
	}
	respond(c, http.StatusOK, "success", r)
}

// GetRecordMessages 获取戏文关联的消息列表
func GetRecordMessages(c *gin.Context) {
	idHex := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { respond(c, http.StatusBadRequest, "invalid id", nil); return }
	var r model.Cassette
	if err := repository.DB().Collection("cassettes").FindOne(c, bson.M{"_id": oid}).Decode(&r); err != nil {
		respond(c, http.StatusNotFound, "not found", nil)
		return
	}
	cur, err := repository.DB().Collection("messages").Find(c, bson.M{"_id": bson.M{"$in": r.MessageIds}}, options.Find().SetSort(bson.M{"createdAt": 1}))
	if err != nil { respond(c, http.StatusInternalServerError, "server error", nil); return }
	var list []model.Message
	_ = cur.All(c, &list)
	respond(c, http.StatusOK, "success", gin.H{"messages": list})
} 