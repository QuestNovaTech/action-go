# Action GO 后端

面向语C扮演场景的后端服务，使用 Go + MongoDB 构建，内置用户鉴权、关系链、群组/房间聊天、招募与戏文（Record）、文件存储（GridFS）等能力。

> 适合新手快速上手：按文档一步步操作，即可在本地跑通“登录 → 发布招募 → 接取招募（进房） → 发消息 → 生成戏文 → 点赞”完整链路。

---

## 1. 环境准备

- Go 1.21+（`go version` 验证）
- MongoDB 6.0+（本地或远程均可）
- Git / PowerShell / curl 或 Postman

### 1.1 安装 Go
- 访问 `https://golang.org/dl` 下载并安装
- 安装后在终端输入：`go version`

国内可选代理（可任选其一）：
```pwsh
# 七牛
go env -w GOPROXY=https://goproxy.cn,direct
# 阿里
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

### 1.2 安装与启动 MongoDB
- 访问 `https://www.mongodb.com/try/download/community` 安装社区版
- 启动方式 A（推荐）：直接运行 MongoDB 服务（默认 27017 端口）
- 启动方式 B（使用仓库自带配置）：
```bash
mkdir -p tmp/mongodb
mongod --config ./mongod.yaml
```
> 说明：`mongod.yaml` 是 MongoDB 进程的启动配置；仅 MongoDB 会读取它，与应用配置不冲突。

---

## 2. 应用配置（必须）

应用读取 `action-go/configs/config.yaml`。若没有，请新建：

```yaml
server:
  port: 8080

jwt:
  # 生产请替换为足够复杂的随机串
  secret: "dev-secret-change-me"
  access_ttl_minutes: 30
  refresh_ttl_days: 14

mongo:
  uri: "mongodb://localhost:27017"
  database: "roleplay"

sms:
  enabled: true
  mock_code: "123456"
```

- `jwt.secret`：用于签名/校验 JWT，必须非空（生产请改为安全随机值）
- `mongo.uri`：与 MongoDB 实际监听一致即可

> 提示：当前代码未做 `${ENV}` 占位符自动展开，如需使用环境变量请告知，我们可补充 BindEnv 支持。

---

## 3. 一键跑通（本地）

```bash
# 安装依赖
cd action-go
go mod tidy

# 第一次使用：生成示例数据
go run ./cmd/seed

# 测试数据库是否正常浏览器访问： 
http://localhost:27017/
# 正常结果：It looks like you are trying to access MongoDB over HTTP on the native driver port.

# 启动服务
go run ./cmd/server
```

看到日志中有：Configuration loaded、MongoDB connected、indexes ensured、Server Information 即表示启动成功。

健康检查：浏览器打开 `http://localhost:8080/healthz` 应显示：`{"status":"ok"}`

---

## 4. 登录与鉴权（新手友好）

1) 发送验证码（Mock）
```bash
curl -X POST http://localhost:8080/api/user/send_code \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000000"}'
```
2) 使用 `configs/config.yaml` 里的 `sms.mock_code` 登录，获取 accessToken
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800000000","code":"123456"}'
```
3) 鉴权头写法（二选一）
- `Authorization: Bearer <ACCESS_TOKEN>`（推荐）
- `Authentication: Bearer <ACCESS_TOKEN>` 或直接 token

示例：
```bash
curl http://localhost:8080/api/user/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

## 5. 文件与头像（GridFS）

- 上传头像（服务端裁剪为正方形、压缩、存 GridFS）：
```bash
curl -X POST http://localhost:8080/api/file/avatar \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -F "file=@/path/to/local/avatar.jpg"
```
- 返回的 `avatar_url`/`thumbnail_url` 形如：`/api/file/{id}`
- 获取文件：
```bash
curl http://localhost:8080/api/file/<FILE_ID> --output avatar.jpg
```
> 当前固定 `image/jpeg` 返回类型；需要多类型 MIME 时可后续增强。

---

## 6. 关系链 / 招募 / 房间 / 消息 / 戏文（端到端）

### 6.1 招募流程
- 创建招募：
```bash
curl -X POST http://localhost:8080/api/recruit/create \
  -H "Authorization: Bearer <ACCESS_TOKEN>" -H "Content-Type: application/json" \
  -d '{
    "backstory_id":"<BACKSTORY_ID>",
    "mode":"couple",
    "myCharacters":["A"],
    "targetCharacters":["B"],
    "title":"招募标题"
  }'
```
- 招募列表/详情：
```bash
curl "http://localhost:8080/api/recruit/list?page=1&size=20" -H "Authorization: Bearer <ACCESS_TOKEN>"
curl http://localhost:8080/api/recruit/detail/<RECRUIT_ID> -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- 接取招募（返回 room_id）：
```bash
curl -X POST http://localhost:8080/api/recruit/<RECRUIT_ID>/accept \
  -H "Authorization: Bearer <ACCESS_TOKEN>" -H "Content-Type: application/json" \
  -d '{"character_id":"B"}'
```

### 6.2 房间发消息 / 拉历史
```bash
# 发消息
curl -X POST http://localhost:8080/api/room/<ROOM_ID>/message \
  -H "Authorization: Bearer <ACCESS_TOKEN>" -H "Content-Type: application/json" \
  -d '{"message_type":"user","element":{"type":"text","text":"hello room"}}'

# 拉历史（按 seq 游标）
curl "http://localhost:8080/api/room/<ROOM_ID>/messages?lastSeq=0&limit=50" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
> 权限：仅房间参与者可发/拉；群聊仅群成员可发/拉；DM 若任一方拉黑另一方则禁止。

### 6.3 生成戏文（Record）
```bash
curl -X POST http://localhost:8080/api/record/create \
  -H "Authorization: Bearer <ACCESS_TOKEN>" -H "Content-Type: application/json" \
  -d '{
    "title":"示例戏文",
    "description":"从若干条消息生成",
    "room_id":"<ROOM_ID>",
    "message_ids":["<MSG_ID_1>","<MSG_ID_2>"]
  }'

# 列表 / 详情 / 消息
curl "http://localhost:8080/api/record/list?page=1&size=20" -H "Authorization: Bearer <ACCESS_TOKEN>"
curl http://localhost:8080/api/record/detail/<RECORD_ID> -H "Authorization: Bearer <ACCESS_TOKEN>"
curl http://localhost:8080/api/record/message/<RECORD_ID> -H "Authorization: Bearer <ACCESS_TOKEN>"
```

### 6.4 点赞（Record/Backstory）
```bash
curl -X POST http://localhost:8080/api/like \
  -H "Authorization: Bearer <ACCESS_TOKEN>" -H "Content-Type: application/json" \
  -d '{"target_type":"record","target_id":"<RECORD_ID>"}'
```

---

## 7. API 文档
- 打开 `docs/openapi.yaml`（IDEA/VSCode 均可预览）
- 文档涵盖用户、关系链、群组、消息、房间、文件、招募、戏文等模块

---

## 8. 常见问题（FAQ）

- 端口被占用：改 `configs/config.yaml` 的 `server.port` 或释放端口
- Mongo 连接失败：
  - 确认 MongoDB 已启动
  - `mongo.uri` 与实际监听匹配（例如 `mongodb://localhost:27017`）
- 种子（seed）执行报错：
  - `scheme must be "mongodb" or "mongodb+srv"` → `mongo.uri` 为空或格式错误，请检查 config
- 登录无效：
  - 确认使用 `sms.mock_code` 作为验证码
  - 鉴权头务必携带 `Authorization: Bearer <token>`
- 文件下载 404：
  - 确认上传成功并使用返回的 `/api/file/{id}` 访问

---

## 9. 生产部署建议（简要）
- 使用强随机 `jwt.secret`
- 将 MongoDB 与应用分离部署，并设置访问控制
- 关闭调试与 Mock，启用必要的速率限制与审计日志
- 结合 Nginx/Traefik 做反向代理与 TLS 终端

---

## 10. 目录结构（关键）
- `cmd/server`：主服务入口
- `cmd/seed`：示例数据生成
- `internal/controller`：各模块控制器（鉴权、用户、关系链、群组、消息、房间、招募、戏文、文件）
- `internal/model`：数据模型
- `internal/router`：路由注册
- `internal/middleware`：鉴权中间件
- `internal/repository`：Mongo/ GridFS 访问
- `internal/indexer`：索引初始化
- `configs/`：应用配置
- `docs/openapi.yaml`：API 文档

## 11. Postman测试
import -> raw text
详见：docs/postman.md
> 有任何问题，先查看日志与本 README 的“常见问题”部分；仍无法解决，欢迎在仓库提出 Issue。
