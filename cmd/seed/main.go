package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"actiondelta/internal/config"
	"actiondelta/internal/indexer"
	"actiondelta/internal/model"
	"actiondelta/internal/repository"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	fmt.Println("Seeding MongoDB with sample data...")
	if err := config.Load(); err != nil {
		panic(err)
	}
	if err := repository.InitMongo(context.Background()); err != nil {
		panic(err)
	}
	defer repository.CloseMongo(context.Background())
	if err := indexer.EnsureAllIndexes(context.Background()); err != nil {
		panic(err)
	}

	db := repository.DB()
	now := time.Now()

	// Users
	u1 := model.User{Phone: "13800000000", UserId: "u_13800000000", UserOpenId: "13800000000", Nickname: "Alice", Avatar: "", CreatedAt: now, UpdatedAt: now}
	u2 := model.User{Phone: "13900000000", UserId: "u_13900000000", UserOpenId: "13900000000", Nickname: "Bob", Avatar: "", CreatedAt: now, UpdatedAt: now}
	upsert := options.Update().SetUpsert(true)
	_, _ = db.Collection("users").UpdateOne(context.Background(), bson.M{"userId": u1.UserId}, bson.M{"$set": u1}, upsert)
	_, _ = db.Collection("users").UpdateOne(context.Background(), bson.M{"userId": u2.UserId}, bson.M{"$set": u2}, upsert)

	// User stats
	_, _ = db.Collection("user_stats").UpdateOne(context.Background(), bson.M{"userId": u1.UserId}, bson.M{"$setOnInsert": bson.M{"followersCount": 0, "followingCount": 0, "postsCount": 0, "updatedAt": now}}, upsert)
	_, _ = db.Collection("user_stats").UpdateOne(context.Background(), bson.M{"userId": u2.UserId}, bson.M{"$setOnInsert": bson.M{"followersCount": 0, "followingCount": 0, "postsCount": 0, "updatedAt": now}}, upsert)

	// Group owned by u1, and both members
	gID := primitive.NewObjectID()
	g := model.Group{ID: gID, Name: "Test Group", Avatar: "", OwnerId: u1.UserId, CreatedAt: now, UpdatedAt: now}
	_, _ = db.Collection("groups").UpdateOne(context.Background(), bson.M{"_id": gID}, bson.M{"$set": g}, upsert)
	_, _ = db.Collection("group_members").UpdateOne(context.Background(), bson.M{"groupId": gID, "userId": u1.UserId}, bson.M{"$set": model.GroupMember{GroupId: gID, UserId: u1.UserId, Role: "owner", JoinedAt: now}}, upsert)
	_, _ = db.Collection("group_members").UpdateOne(context.Background(), bson.M{"groupId": gID, "userId": u2.UserId}, bson.M{"$set": model.GroupMember{GroupId: gID, UserId: u2.UserId, Role: "member", JoinedAt: now}}, upsert)

	// DM conversation between u1 and u2
	dmId := "dm_" + u1.UserId + "_" + u2.UserId
	_, _ = db.Collection("conversations").UpdateOne(context.Background(), bson.M{"conversationId": dmId}, bson.M{"$setOnInsert": model.Conversation{ConversationId: dmId, ConversationType: "dm", Participants: []string{u1.UserId, u2.UserId}, LastSeq: 0, LastMessage: "", UpdatedAt: now}}, upsert)

	// Backstory (简单模拟一条)
	bsID := primitive.NewObjectID()
	_, _ = db.Collection("backstories").UpdateOne(context.Background(), bson.M{"_id": bsID}, bson.M{"$setOnInsert": bson.M{
		"title":        "示例剧本",
		"subtitle":     "副标题",
		"authorName":   "System",
		"authorUserId": u1.UserId,
		"tags":         []string{"恋爱", "校园"},
		"createdAt":    now,
		"updatedAt":    now,
	}}, upsert)

	// Recruit (关联该剧本)
	reID := primitive.NewObjectID()
	rec := model.Recruit{ID: reID, Title: "对戏招募-示例", BackstoryId: bsID, CreatorId: u1.UserId, Mode: "couple", MyCharacters: []string{"A"}, TargetCharacters: []string{"B"}, Status: "active", CreatedAt: now, UpdatedAt: now}
	_, _ = db.Collection("recruits").UpdateOne(context.Background(), bson.M{"_id": reID}, bson.M{"$set": rec}, upsert)

	// Theater for recruit
	thID := primitive.NewObjectID()
	th := model.Theater{ID: thID, RecruitId: reID, Title: "演绎房间", Mode: "couple", Status: "active", Participants: []model.TheaterParticipant{{UserId: u1.UserId, CostumeId: "A", JoinTime: now}}, CreatedAt: now, UpdatedAt: now}
	_, _ = db.Collection("theaters").UpdateOne(context.Background(), bson.M{"_id": thID}, bson.M{"$set": th}, upsert)

	// Two messages in DM
	m1 := model.Message{ID: primitive.NewObjectID(), ConversationId: dmId, ConversationType: "dm", Seq: 1, SenderUserId: u1.UserId, MessageType: "user", Element: model.MessageElement{Type: "text", Data: bson.M{"text": "你好"}}, CreatedAt: now, UpdatedAt: now}
	m2 := model.Message{ID: primitive.NewObjectID(), ConversationId: dmId, ConversationType: "dm", Seq: 2, SenderUserId: u2.UserId, MessageType: "user", Element: model.MessageElement{Type: "text", Data: bson.M{"text": "你好呀"}}, CreatedAt: now.Add(time.Second), UpdatedAt: now.Add(time.Second)}
	_, _ = db.Collection("messages").UpdateOne(context.Background(), bson.M{"_id": m1.ID}, bson.M{"$set": m1}, upsert)
	_, _ = db.Collection("messages").UpdateOne(context.Background(), bson.M{"_id": m2.ID}, bson.M{"$set": m2}, upsert)

	// Cassette built from the two messages
	csID := primitive.NewObjectID()
	cass := model.Cassette{ID: csID, Title: "示例戏文", Description: "由两条消息生成", CreatorId: u1.UserId, Participants: []model.CassetteParticipant{{UserId: u1.UserId}, {UserId: u2.UserId}}, MessageIds: []primitive.ObjectID{m1.ID, m2.ID}, CreatedAt: now, UpdatedAt: now}
	_, _ = db.Collection("cassettes").UpdateOne(context.Background(), bson.M{"_id": csID}, bson.M{"$set": cass}, upsert)

	fmt.Println("Seed complete. Data:")
	fmt.Printf("  Users: %s (phone=%s), %s (phone=%s)\n", u1.UserId, u1.Phone, u2.UserId, u2.Phone)
	fmt.Printf("  Group ID: %s (use as conversation_id for group messaging)\n", gID.Hex())
	fmt.Printf("  DM Conversation ID: %s\n", dmId)
	fmt.Printf("  Backstory ID: %s\n", bsID.Hex())
	fmt.Printf("  Recruit ID: %s\n", reID.Hex())
	fmt.Printf("  Theater(Room) ID: %s\n", thID.Hex())
	fmt.Printf("  Cassette ID: %s\n", csID.Hex())
}
