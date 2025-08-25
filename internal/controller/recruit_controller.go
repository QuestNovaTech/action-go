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

// ListRecruits 招募列表（支持分页与基础筛选）
func ListRecruits(c *gin.Context) {
	page := parseIntDefault(c.DefaultQuery("page", "1"), 1)
	size := parseIntDefault(c.DefaultQuery("size", "20"), 20)
	mode := c.Query("mode")
	status := c.Query("status")
	backstory := c.Query("backstory_id")
	keyword := c.Query("keyword")

	filter := bson.M{}
	if mode != "" {
		filter["mode"] = mode
	}
	if status != "" {
		filter["status"] = status
	}
	if backstory != "" {
		if oid, err := primitive.ObjectIDFromHex(backstory); err == nil {
			filter["backstoryId"] = oid
		}
	}
	if keyword != "" {
		filter["title"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	col := repository.DB().Collection("recruits")
	total, _ := col.CountDocuments(c, filter)
	cur, err := col.Find(c, filter, options.Find().SetSort(bson.M{"createdAt": -1}).SetSkip(int64((page-1)*size)).SetLimit(int64(size)))
	if err != nil {
		respond(c, http.StatusInternalServerError, "查询失败", nil)
		return
	}
	var list []model.Recruit
	_ = cur.All(c, &list)
	respond(c, http.StatusOK, "success", gin.H{"total": total, "list": list})
}

// GetRecruit 招募详情
func GetRecruit(c *gin.Context) {
	idHex := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid id", nil)
		return
	}
	var r model.Recruit
	if err := repository.DB().Collection("recruits").FindOne(c, bson.M{"_id": oid}).Decode(&r); err != nil {
		respond(c, http.StatusNotFound, "not found", nil)
		return
	}
	respond(c, http.StatusOK, "success", r)
}

// CreateRecruit 创建招募
func CreateRecruit(c *gin.Context) {
	userId := c.GetString("userId")
	var body struct {
		BackstoryId      string                  `json:"backstory_id"`
		Mode             string                  `json:"mode"`
		MyCharacters     []string                `json:"myCharacters"`
		TargetCharacters []string                `json:"targetCharacters"`
		Title            string                  `json:"title"`
		CustomContent    string                  `json:"customContent"`
		CustomCharacters []model.CustomCharacter `json:"customCharacters"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.BackstoryId == "" {
		respond(c, http.StatusBadRequest, "invalid request", nil)
		return
	}
	bid, err := primitive.ObjectIDFromHex(body.BackstoryId)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid backstory id", nil)
		return
	}
	now := time.Now()
	rec := model.Recruit{
		Title:            body.Title,
		BackstoryId:      bid,
		CreatorId:        userId,
		Mode:             body.Mode,
		MyCharacters:     body.MyCharacters,
		TargetCharacters: body.TargetCharacters,
		CustomContent:    body.CustomContent,
		CustomCharacters: body.CustomCharacters,
		Status:           "active",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	res, err := repository.DB().Collection("recruits").InsertOne(c, rec)
	if err != nil {
		respond(c, http.StatusInternalServerError, "server error", nil)
		return
	}
	id := res.InsertedID.(primitive.ObjectID)
	respond(c, http.StatusOK, "success", gin.H{"id": id.Hex()})
}

// DeleteRecruit 删除招募（仅发布者）
func DeleteRecruit(c *gin.Context) {
	userId := c.GetString("userId")
	idHex := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid id", nil)
		return
	}
	res := repository.DB().Collection("recruits").FindOneAndDelete(c, bson.M{"_id": oid, "creatorId": userId})
	if res.Err() != nil {
		respond(c, http.StatusForbidden, "forbidden or not found", nil)
		return
	}
	respond(c, http.StatusOK, "success", nil)
}

// AcceptRecruit 接取招募 -> 创建/加入房间
func AcceptRecruit(c *gin.Context) {
	userId := c.GetString("userId")
	idHex := c.Param("id")
	var body struct {
		CharacterId string `json:"character_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond(c, http.StatusBadRequest, "invalid request", nil)
		return
	}
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		respond(c, http.StatusBadRequest, "invalid id", nil)
		return
	}
	// 查找 theater，按 recruitId 复用/创建
	var th model.Theater
	err = repository.DB().Collection("theaters").FindOne(c, bson.M{"recruitId": oid}).Decode(&th)
	if err != nil {
		th = model.Theater{RecruitId: oid, Title: "演绎房间", Mode: "couple", Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		res, _ := repository.DB().Collection("theaters").InsertOne(c, th)
		th.ID = res.InsertedID.(primitive.ObjectID)
	}
	// 加入参与者（去重）
	exists := false
	for _, p := range th.Participants {
		if p.UserId == userId {
			exists = true
			break
		}
	}
	if !exists {
		th.Participants = append(th.Participants, model.TheaterParticipant{UserId: userId, CostumeId: body.CharacterId, JoinTime: time.Now()})
		_, _ = repository.DB().Collection("theaters").UpdateByID(c, th.ID, bson.M{"$set": bson.M{"participants": th.Participants, "updatedAt": time.Now()}})
	}
	respond(c, http.StatusOK, "success", gin.H{"room_id": th.ID.Hex()})
}

func parseIntDefault(s string, def int) int {
	var x int
	_, err := fmtSscan(s, &x)
	if err != nil || x <= 0 {
		return def
	}
	return x
}

// 轻量 fmt.Sscan 等价，避免直接引入 fmt 造成未使用告警
func fmtSscan(s string, p *int) (int, error) { return fmtSscanImpl(s, p) }

// 使用内联实现
func fmtSscanImpl(s string, p *int) (int, error) {
	// 简单转换
	x := 0
	sign := 1
	for i, b := range []byte(s) {
		if i == 0 && b == '-' {
			sign = -1
			continue
		}
		if b < '0' || b > '9' {
			return 0, fmtErr()
		}
		x = x*10 + int(b-'0')
	}
	*p = sign * x
	return 1, nil
}

func fmtErr() error {
	return &scanErr{}
}

type scanErr struct{}

func (e *scanErr) Error() string { return "scan error" }
