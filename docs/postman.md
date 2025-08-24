# Postman
## 环境变量
- 必填环境变量（运行前手动填）
  - baseUrl: http://localhost:8080
  - phone: 13800000000
  - code: 123456
  - backstoryId: 从 seed 脚本输出复制（例如 65f…abc）
  - avatarFile: 本地头像绝对路径（留空则跳过头像上传相关用例）
  - otherUserId: 可选，用于关注/好友等双人用例（示例 u_13900000000）

- 运行中自动写入（初始置空）
  - accessToken: ""
  - refreshToken: ""
  - userId: ""
  - avatarFileId: ""
  - avatarUrl: ""
  - recruitId: ""
  - roomId: ""
  - conversationId: ""
  - msgId1: ""
  - msgId2: ""
  - recordId: ""

- 可能需要手动临时设定（仅当你单独跑对应请求）
  - requestId: 好友申请 ID（若未按集合顺序运行、或未在测试脚本中自动提取时，手动设置）

注意：
跑脚本前手动选一次文件
在 10-File/Upload Avatar 请求 → Body → form-data → 把 file 行里的 value 从 Text 切到 File，然后点 Select Files 选本地 jpg
原因：环境变量 avatarFile只是一个字符串，例如 C:\Users\xxx\avatar.jpg。form-data file 字段	需要 Postman 在 UI 里 真正选中文件，内部才会生成一个 File 对象并附带 filename、Content-Type、binary 等信息

## import

{
  "info": {
    "name": "Action Backend Full E2E (All APIs by OpenAPI)",
    "_postman_id": "b5b2e2f1-4f27-4c0e-9d2f-action-e2e",
    "description": "Auth -> User -> File -> Relations -> Groups -> Recruit -> Room -> Messaging -> Record -> Like (ordered for front-end flow testing). Use environment variables: baseUrl, phone, code, accessToken, userId, avatarFile, backstoryId, otherUserId(optional).",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json     "
  },
  "auth": {
    "type": "bearer",
    "bearer": [
      { "key": "token", "value": "{{accessToken}}", "type": "string" }
    ]
  },
  "item": [
    {
      "name": "00-Healthz (public)",
      "auth": { "type": "noauth" },
      "request": { "method": "GET", "url": "{{baseUrl}}/healthz" }
    },
    {
      "name": "01-Auth/SendCode (public)",
      "auth": { "type": "noauth" },
      "request": {
        "method": "POST",
        "header": [ { "key": "Content-Type", "value": "application/json" } ],
        "url": "{{baseUrl}}/api/user/send_code",
        "body": { "mode": "raw", "raw": "{\n  \"phone\": \"{{phone}}\"\n}" }
      }
    },
    {
      "name": "02-Auth/Login (public -> save accessToken)",
      "auth": { "type": "noauth" },
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json();",
          "if(j&&j.data){",
          "  if(j.data.accessToken) pm.environment.set('accessToken', j.data.accessToken);",
          "  if(j.data.refreshToken) pm.environment.set('refreshToken', j.data.refreshToken);",
          "}"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method": "POST",
        "header": [ { "key": "Content-Type", "value": "application/json" } ],
        "url": "{{baseUrl}}/api/user/login",
        "body": { "mode": "raw", "raw": "{\n  \"phone\": \"{{phone}}\",\n  \"code\": \"{{code}}\"\n}" }
      }
    },
    {
      "name": "03-Auth/Refresh",
      "request": {
        "method": "POST",
        "header": [ { "key": "Content-Type", "value": "application/json" } ],
        "url": "{{baseUrl}}/api/auth/refresh",
        "body": { "mode": "raw", "raw": "{\n  \"refreshToken\": \"{{refreshToken}}\"\n}" }
      }
    },
    {
      "name": "04-Auth/OneClickLogin (public, optional)",
      "auth": { "type": "noauth" },
      "request": {
        "method": "POST",
        "header": [ { "key": "Content-Type", "value": "application/json" } ],
        "url": "{{baseUrl}}/api/user/oneclick_login",
        "body": { "mode": "raw", "raw": "{\n  \"phone\":\"{{phone}}\",\n  \"device_id\":\"dev-001\",\n  \"platform\":\"web\"\n}" }
      }
    },
    {
      "name": "05-User/GetMe (save userId)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json();",
          "if(j&&j.data&&j.data.user_id){ pm.environment.set('userId', j.data.user_id); }"
        ],"type":"text/javascript"}}
      ],
      "request": { "method": "GET", "url": "{{baseUrl}}/api/user/me" }
    },
    {
      "name": "06-User/UpdateMe",
      "request": {
        "method": "PUT",
        "header": [ { "key": "Content-Type", "value": "application/json" } ],
        "url": "{{baseUrl}}/api/user/me",
        "body": { "mode": "raw", "raw": "{\n  \"nickname\":\"Tester\",\n  \"avatar\":\"\",\n  \"gender\":\"\",\n  \"bio\":\"Hello from Postman\"\n}" }
      }
    },
    {
      "name": "07-User/Heartbeat",
      "request": { "method": "POST", "url": "{{baseUrl}}/api/user/heartbeat" }
    },
    {
      "name": "08-User/Profile",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/user/profile/{{userId}}" }
    },
    {
      "name": "09-User/Activities",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/user/activities/{{userId}}?last_id=&limit=20" }
    },
    {
      "name": "10-File/Upload Avatar (form-data -> save avatarFileId)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json();",
          "if(j&&j.data&&j.data.avatar_url){",
          "  pm.environment.set('avatarUrl', j.data.avatar_url);",
          "  try{ pm.environment.set('avatarFileId', j.data.avatar_url.split('/').pop()); }catch(e){}",
          "}"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method": "POST",
        "url": "{{baseUrl}}/api/file/avatar",
        "body": { "mode":"formdata", "formdata":[ { "key":"file", "type":"file", "src":"{{avatarFile}}" } ] }
      }
    },
    {
      "name": "11-File/Get Avatar (public)",
      "auth": { "type": "noauth" },
      "request": { "method": "GET", "url": "{{baseUrl}}/api/file/{{avatarFileId}}" }
    },
    {
      "name": "12-Follow/Follow (optional, set {{otherUserId}})",
      "request": { "method": "POST", "url": "{{baseUrl}}/api/relation/follow/{{otherUserId}}" }
    },
    {
      "name": "13-Follow/Unfollow (optional)",
      "request": { "method": "DELETE", "url": "{{baseUrl}}/api/relation/follow/{{otherUserId}}" }
    },
    {
      "name": "14-Follow/Status (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/follow/status/{{otherUserId}}" }
    },
    {
      "name": "15-Follow/Followers (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/followers" }
    },
    {
      "name": "16-Follow/Following (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/following" }
    },
    {
      "name": "17-Friends/Create Request (optional, set {{otherUserId}})",
      "request": {
        "method": "POST",
        "header": [ { "key":"Content-Type","value":"application/json" } ],
        "url": "{{baseUrl}}/api/relation/friend/request",
        "body": { "mode":"raw", "raw":"{\n  \"user_id\":\"{{otherUserId}}\",\n  \"greeting\":\"hi\"\n}" }
      }
    },
    {
      "name": "18-Friends/Respond (optional, set {{requestId}})",
      "request": {
        "method": "POST",
        "header": [ { "key":"Content-Type","value":"application/json" } ],
        "url": "{{baseUrl}}/api/relation/friend/respond",
        "body": { "mode":"raw", "raw":"{\n  \"request_id\":\"{{requestId}}\",\n  \"action\":\"accept\"\n}" }
      }
    },
    {
      "name": "19-Friends/Requests (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/friend/requests" }
    },
    {
      "name": "20-Friends/List (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/friends" }
    },
    {
      "name": "21-Friends/Delete (optional)",
      "request": { "method": "DELETE", "url": "{{baseUrl}}/api/relation/friend/{{otherUserId}}" }
    },
    {
      "name": "22-Blocks/Block (optional)",
      "request": { "method": "POST", "url": "{{baseUrl}}/api/relation/block/{{otherUserId}}" }
    },
    {
      "name": "23-Blocks/Unblock (optional)",
      "request": { "method": "DELETE", "url": "{{baseUrl}}/api/relation/block/{{otherUserId}}" }
    },
    {
      "name": "24-Blocks/List (optional)",
      "request": { "method": "GET", "url": "{{baseUrl}}/api/relation/blocks" }
    },
    {
      "name": "25-Groups/Create (save groupId)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json(); if(j&&j.data&&j.data.group_id){ pm.environment.set('groupId', j.data.group_id); }"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method": "POST",
        "header": [ { "key":"Content-Type","value":"application/json" } ],
        "url": "{{baseUrl}}/api/group",
        "body": { "mode":"raw", "raw":"{\n  \"name\":\"My Group\",\n  \"avatar\":\"\"\n}" }
      }
    },
    {
      "name": "26-Groups/Add Members (optional)",
      "request": {
        "method": "POST",
        "header": [ { "key":"Content-Type","value":"application/json" } ],
        "url": "{{baseUrl}}/api/group/{{groupId}}/members",
        "body": { "mode":"raw", "raw":"{\n  \"user_ids\":[\"{{otherUserId}}\"]\n}" }
      }
    },
    {
      "name": "27-Groups/Remove Member (optional)",
      "request": { "method":"DELETE", "url":"{{baseUrl}}/api/group/{{groupId}}/members/{{otherUserId}}" }
    },
    {
      "name": "28-Groups/My Groups (optional)",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/group/my" }
    },
    {
      "name": "29-Groups/Detail (optional)",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/group/{{groupId}}" }
    },
    {
      "name": "30-Recruit/List",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/recruit/list?page=1&size=20" }
    },
    {
      "name": "31-Recruit/Create (save recruitId) - requires {{backstoryId}}",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json(); if(j&&j.data&&j.data.id){ pm.environment.set('recruitId', j.data.id); }"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/recruit/create",
        "body":{"mode":"raw","raw":"{\n  \"backstory_id\":\"{{backstoryId}}\",\n  \"mode\":\"couple\",\n  \"myCharacters\":[\"A\"],\n  \"targetCharacters\":[\"B\"],\n  \"title\":\"Recruit from Postman\"\n}"}
      }
    },
    {
      "name": "32-Recruit/Detail",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/recruit/detail/{{recruitId}}" }
    },
    {
      "name": "33-Recruit/Accept (save roomId)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json(); if(j&&j.data&&j.data.room_id){ pm.environment.set('roomId', j.data.room_id); pm.environment.set('conversationId', j.data.room_id); }"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/recruit/{{recruitId}}/accept",
        "body":{"mode":"raw","raw":"{\n  \"character_id\":\"B\"\n}"}
      }
    },
    {
      "name": "34-Recruit/Delete (optional cleanup)",
      "request": { "method":"DELETE", "url":"{{baseUrl}}/api/recruit/{{recruitId}}" }
    },
    {
      "name": "35-Room/Send Message #1",
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/room/{{roomId}}/message",
        "body":{"mode":"raw","raw":"{\n  \"message_type\":\"user\",\n  \"element\":{ \"type\":\"text\", \"text\":\"hello room\" }\n}"}
      }
    },
    {
      "name": "36-Room/Send Message #2",
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/room/{{roomId}}/message",
        "body":{"mode":"raw","raw":"{\n  \"message_type\":\"user\",\n  \"element\":{ \"type\":\"text\", \"text\":\"another msg\" }\n}"}
      }
    },
    {
      "name": "37-Room/Get Messages (save msgId1,msgId2)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json();",
          "if(j&&j.data&&Array.isArray(j.data.messages)){",
          "  if(j.data.messages[0]&&j.data.messages[0].id) pm.environment.set('msgId1', j.data.messages[0].id);",
          "  if(j.data.messages[1]&&j.data.messages[1].id) pm.environment.set('msgId2', j.data.messages[1].id);",
          "}"
        ],"type":"text/javascript"}}
      ],
      "request": { "method":"GET", "url":"{{baseUrl}}/api/room/{{roomId}}/messages?lastSeq=0&limit=50" }
    },
    {
      "name": "38-Message/Send (generic) - dm/group/room",
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/message/send",
        "body":{"mode":"raw","raw":"{\n  \"conversation_type\":\"room\",\n  \"conversation_id\":\"{{roomId}}\",\n  \"message_type\":\"user\",\n  \"element\":{ \"type\":\"text\", \"text\":\"via generic send\" }\n}"}
      }
    },
    {
      "name": "39-Message/History (generic) - dm/group/room",
      "request": {
        "method":"GET",
        "url":"{{baseUrl}}/api/message/history?conversation_type=room&conversation_id={{roomId}}&lastSeq=0&limit=50"
      }
    },
    {
      "name": "40-Record/Create (save recordId)",
      "event": [
        { "listen":"test","script":{"exec":[
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json(); if(j&&j.data&&j.data.id){ pm.environment.set('recordId', j.data.id); }"
        ],"type":"text/javascript"}}
      ],
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/record/create",
        "body":{"mode":"raw","raw":"{\n  \"title\":\"Record from Postman\",\n  \"description\":\"generated from room messages\",\n  \"room_id\":\"{{roomId}}\",\n  \"message_ids\":[\"{{msgId1}}\",\"{{msgId2}}\"]\n}"}
      }
    },
    {
      "name": "41-Record/Detail",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/record/detail/{{recordId}}" }
    },
    {
      "name": "42-Record/Messages",
      "request": { "method":"GET", "url":"{{baseUrl}}/api/record/message/{{recordId}}" }
    },
    {
      "name": "43-Like/Toggle Record",
      "request": {
        "method":"POST",
        "header":[{"key":"Content-Type","value":"application/json"}],
        "url":"{{baseUrl}}/api/like",
        "body":{"mode":"raw","raw":"{\n  \"target_type\":\"record\",\n  \"target_id\":\"{{recordId}}\"\n}"}
      }
    },
{
  "name": "44-Record/List",
  "request": {
    "method": "GET",
    "url": "{{baseUrl}}/api/record/list?page=1&size=20"
  }
},
{
  "name": "45-Room/Join (explicit)",
  "event": [
    {
      "listen": "test",
      "script": {
        "exec": [
          "pm.test('200',()=>pm.response.to.have.status(200));",
          "let j=pm.response.json(); if(j&&j.data&&j.data.room_id){ pm.environment.set('roomId', j.data.room_id); pm.environment.set('conversationId', j.data.room_id); }"
        ],
        "type": "text/javascript"
      }
    }
  ],
  "request": {
    "method": "POST",
    "header": [ { "key": "Content-Type", "value": "application/json" } ],
    "url": "{{baseUrl}}/api/room/join",
    "body": {
      "mode": "raw",
      "raw": "{\n  \"recruit_id\":\"{{recruitId}}\",\n  \"character_id\":\"B\"\n}"
    }
  }
},
{
  "name": "46-Admin/UserList",
  "request": {
    "method": "GET",
    "url": "{{baseUrl}}/api/admin/userList"
  }
}
  ]
}


