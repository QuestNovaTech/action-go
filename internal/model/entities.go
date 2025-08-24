package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户基本资料。
type User struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Phone      string             `bson:"phone" json:"phone"`
    Nickname   string             `bson:"nickname" json:"nickname"`
    Avatar     string             `bson:"avatar" json:"avatar"`
    UserId     string             `bson:"userId" json:"user_id"`
    UserOpenId string             `bson:"userOpenId" json:"user_open_id"`
    WordCount  int                `bson:"wordCount" json:"word_count"`
    CoinCount  int                `bson:"coinCount" json:"coin_count"`
    GemCount   int                `bson:"gemCount" json:"gem_count"`
    Gender     string             `bson:"gender" json:"gender"`
    Bio        string             `bson:"bio" json:"bio"`
    Online     bool               `bson:"online" json:"online"`
    LastSeenAt time.Time          `bson:"lastSeenAt" json:"last_seen_at"`
    CreatedAt  time.Time          `bson:"createdAt" json:"created_at"`
    UpdatedAt  time.Time          `bson:"updatedAt" json:"updated_at"`
    DeletedAt  *time.Time         `bson:"deletedAt" json:"deleted_at"`
}

type AuthCode struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Phone     string             `bson:"phone" json:"phone"`
    Code      string             `bson:"code" json:"code"`
    Scene     string             `bson:"scene" json:"scene"`
    CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
    ExpireAt  time.Time          `bson:"expireAt" json:"expire_at"`
}

type FriendRequest struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    RequesterId string             `bson:"requesterId" json:"requester_id"`
    RecipientId string             `bson:"recipientId" json:"recipient_id"`
    Greeting    string             `bson:"greeting" json:"greeting"`
    Status      string             `bson:"status" json:"status"` // pending 待处理 / accepted 已同意 / rejected 已拒绝
    CreatedAt   time.Time          `bson:"createdAt" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updatedAt" json:"updated_at"`
}

type FriendEdge struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserA     string             `bson:"userA" json:"user_a"`
    UserB     string             `bson:"userB" json:"user_b"`
    CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
}

type BlockEdge struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserId        string             `bson:"userId" json:"user_id"`
    BlockedUserId string             `bson:"blockedUserId" json:"blocked_user_id"`
    CreatedAt     time.Time          `bson:"createdAt" json:"created_at"`
}

type Group struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Avatar    string             `bson:"avatar" json:"avatar"`
    OwnerId   string             `bson:"ownerId" json:"owner_id"`
    CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
    UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`
}

type GroupMember struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    GroupId  primitive.ObjectID `bson:"groupId" json:"group_id"`
    UserId   string             `bson:"userId" json:"user_id"`
    Role     string             `bson:"role" json:"role"` // owner 群主 / admin 管理员 / member 普通成员
    JoinedAt time.Time          `bson:"joinedAt" json:"joined_at"`
}

type Conversation struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ConversationId   string             `bson:"conversationId" json:"conversation_id"`
    ConversationType string             `bson:"conversationType" json:"conversation_type"` // dm 私聊 / group 群聊 / room 演绎房间
    Participants     []string           `bson:"participants" json:"participants"`
    LastSeq          int64              `bson:"lastSeq" json:"last_seq"`
    LastMessage      string             `bson:"lastMessage" json:"last_message"`
    UpdatedAt        time.Time          `bson:"updatedAt" json:"updated_at"`
}

type Message struct {
    ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
    ConversationId   string              `bson:"conversationId" json:"conversation_id"`
    ConversationType string              `bson:"conversationType" json:"conversation_type"`
    Seq              int64               `bson:"seq" json:"seq"`
    RoomId           *primitive.ObjectID `bson:"roomId,omitempty" json:"room_id,omitempty"`
    GroupId          *primitive.ObjectID `bson:"groupId,omitempty" json:"group_id,omitempty"`
    SenderUserId     string              `bson:"senderUserId" json:"sender_user_id"`
    MessageType      string              `bson:"messageType" json:"message_type"`
    Element          MessageElement      `bson:"element" json:"element"`
    CharacterInfo    *CharacterInfo      `bson:"characterInfo,omitempty" json:"character_info,omitempty"`
    CreatedAt        time.Time           `bson:"createdAt" json:"created_at"`
    UpdatedAt        time.Time           `bson:"updatedAt" json:"updated_at"`
    DeletedAt        *time.Time          `bson:"deletedAt" json:"deleted_at"`
}

type MessageElement struct {
    Type string                 `bson:"type" json:"type"`
    Data map[string]interface{} `bson:"data" json:"data"`
}

type CharacterInfo struct {
    CharacterId string `bson:"characterId" json:"character_id"`
    Name        string `bson:"name" json:"name"`
    Avatar      string `bson:"avatar" json:"avatar"`
}

type Theater struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    RecruitId       primitive.ObjectID   `bson:"recruitId" json:"recruit_id"`
    BackstoryId     primitive.ObjectID   `bson:"backstoryId" json:"backstory_id"`
    Title           string               `bson:"title" json:"title"`
    Subtitle        string               `bson:"subtitle" json:"subtitle"`
    Mode            string               `bson:"mode" json:"mode"`
    BackgroundStory string               `bson:"backgroundStory" json:"background_story"`
    Participants    []TheaterParticipant `bson:"participants" json:"participants"`
    Status          string               `bson:"status" json:"status"`
    CreatedAt       time.Time            `bson:"createdAt" json:"created_at"`
    UpdatedAt       time.Time            `bson:"updatedAt" json:"updated_at"`
}

type TheaterParticipant struct {
    UserId       string    `bson:"userId" json:"user_id"`
    CostumeId    string    `bson:"costumeId" json:"costume_id"`
    CostumeName  string    `bson:"costumeName" json:"costume_name"`
    Avatar       string    `bson:"avatar" json:"avatar"`
    JoinTime     time.Time `bson:"joinTime" json:"join_time"`
    MessageCount int       `bson:"messageCount" json:"message_count"`
}

// FollowEdge 关注关系
type FollowEdge struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FollowerId  string             `bson:"followerId" json:"follower_id"`
    FollowingId string             `bson:"followingId" json:"following_id"`
    CreatedAt   time.Time          `bson:"createdAt" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updatedAt" json:"updated_at"`
}

// UserStats 用户统计信息
type UserStats struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserId         string             `bson:"userId" json:"user_id"`
    FollowersCount int                `bson:"followersCount" json:"followers_count"`
    FollowingCount int                `bson:"followingCount" json:"following_count"`
    PostsCount     int                `bson:"postsCount" json:"posts_count"`
    UpdatedAt      time.Time          `bson:"updatedAt" json:"updated_at"`
}

// UserActivity 用户活动（用于主页最近动态）
type UserActivity struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserId       string             `bson:"userId" json:"user_id"`
    ActivityType string             `bson:"activityType" json:"activity_type"`
    TargetType   string             `bson:"targetType" json:"target_type"`
    TargetId     string             `bson:"targetId" json:"target_id"`
    Title        string             `bson:"title" json:"title"`
    Content      string             `bson:"content" json:"content"`
    CreatedAt    time.Time          `bson:"createdAt" json:"created_at"`
}

// Recruit 招募实体，用于发起演绎
type Recruit struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title            string             `bson:"title" json:"title"`
    BackstoryId      primitive.ObjectID `bson:"backstoryId" json:"backstory_id"`
    CreatorId        string             `bson:"creatorId" json:"creator_id"`
    Mode             string             `bson:"mode" json:"mode"` // couple/crowd/drama
    MyCharacters     []string           `bson:"myCharacters" json:"my_characters"`
    TargetCharacters []string           `bson:"targetCharacters" json:"target_characters"`
    CustomContent    string             `bson:"customContent" json:"custom_content"`
    CustomCharacters []CustomCharacter  `bson:"customCharacters" json:"custom_characters"`
    Status           string             `bson:"status" json:"status"` // active/completed/cancelled
    CreatedAt        time.Time          `bson:"createdAt" json:"created_at"`
    UpdatedAt        time.Time          `bson:"updatedAt" json:"updated_at"`
    DeletedAt        *time.Time         `bson:"deletedAt" json:"deleted_at"`
}

// CustomCharacter 自定义角色定义
type CustomCharacter struct {
    CharacterId  string `bson:"characterId" json:"character_id"`
    Name         string `bson:"name" json:"name"`
    Avatar       string `bson:"avatar" json:"avatar"`
    Illustration string `bson:"illustration" json:"illustration"`
    Background   string `bson:"background" json:"background"`
}

// Cassette 戏文
type Cassette struct {
    ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    Title        string               `bson:"title" json:"title"`
    Description  string               `bson:"description" json:"description"`
    BackstoryId  *primitive.ObjectID  `bson:"backstoryId,omitempty" json:"backstory_id,omitempty"`
    RoomId       *primitive.ObjectID  `bson:"roomId,omitempty" json:"room_id,omitempty"`
    CreatorId    string               `bson:"creatorId" json:"creator_id"`
    Participants []CassetteParticipant `bson:"participants" json:"participants"`
    MessageIds   []primitive.ObjectID `bson:"messageIds" json:"message_ids"`
    LikeCount    int                  `bson:"likeCount" json:"like_count"`
    ViewCount    int                  `bson:"viewCount" json:"view_count"`
    CreatedAt    time.Time            `bson:"createdAt" json:"created_at"`
    UpdatedAt    time.Time            `bson:"updatedAt" json:"updated_at"`
    DeletedAt    *time.Time           `bson:"deletedAt" json:"deleted_at"`
}

// CassetteParticipant 戏文参与者
type CassetteParticipant struct {
    UserId        string `bson:"userId" json:"user_id"`
    CharacterId   string `bson:"characterId" json:"character_id"`
    CharacterName string `bson:"characterName" json:"character_name"`
}

// Like 点赞
type Like struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserId     string             `bson:"userId" json:"user_id"`
    TargetType string             `bson:"targetType" json:"target_type"` // backstory/record
    TargetId   primitive.ObjectID `bson:"targetId" json:"target_id"`
    CreatedAt  time.Time          `bson:"createdAt" json:"created_at"`
}

