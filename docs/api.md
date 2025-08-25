





---

附加说明（对接与规范）
- 鉴权请求头：兼容两种写法
  - Authorization: Bearer <accessToken>
  - Authentication: Bearer <accessToken> 或直接 <accessToken>
- 文件与头像
  - 上传：POST /api/file/avatar （multipart/form-data: file），服务端裁剪压缩并写入 MongoDB GridFS
  - 访问：GET /api/file/{id}（返回 image/jpeg），用户表中 avatar/thumbnail 保存为对应 API URL
- 消息权限
  - group：仅群成员可发/拉取
  - room：仅房间 participants 可发/拉取
  - dm：若会话已存在且非 participants，拒绝访问
  - 黑名单：后续补充，影响 DM 与可见性
- Recruit / Record（刘茂负责）
  - Recruit：列表/详情/创建/删除，入房间 JoinRoom；前端链路为：剧本详情 -> 招募发布 -> 招募详情（房间）
  - Record：从会话/房间选取消息生成戏文，列表/详情/消息；点赞可选
- 分页：消息使用 seq 游标；列表使用 id/时间游标或页码（与前端确认）

- 入房接口使用注意事项（Recruit/Accept 与 Room/Join）
  - Recruit/Accept（POST /api/recruit/{id}/accept）
    - 语义：在“招募详情页”内接取，路径上携带 recruitId，Body 仅需 character_id。
    - 适用：招募流标准入口；推荐前端主链路使用。
  - Room/Join（POST /api/room/join）
    - 语义：通用入房入口，Body 携带 recruit_id + character_id。
    - 适用：从非招募详情页（如活动、通知）直接入房。
  - 两者都会：按 recruitId 查找/创建房间（theaters），将当前用户追加到 participants 并返回 room_id。
  - 测试链路避免重复：二者功能可互换，联调时二选一即可（建议优先使用 Recruit/Accept）。 

---

当前项目支持的 46 个 API（作用说明）

公共与鉴权
- GET /healthz：健康检查
- POST /api/user/send_code：发送登录验证码（Mock/真实通道）
- POST /api/user/login：手机号+验证码登录，自动注册/签发 token
- POST /api/auth/refresh：用刷新令牌换新访问令牌
- POST /api/user/oneclick_login：一键登录（本地模拟）

用户
- GET /api/user/me：获取当前用户资料（通过 token）
- PUT /api/user/me：更新当前用户资料（昵称、头像、性别、签名）
- GET /api/user/profile/{user_id}：用户主页（在线状态/粉丝/关注统计等）
- GET /api/user/activities/{user_id}：用户最近活动（游标分页）
- POST /api/user/heartbeat：心跳上报（更新 lastSeenAt，用于在线状态）

文件
- POST /api/file/avatar：上传头像（multipart），服务端裁剪压缩并存入 GridFS
- GET /api/file/{id}：按文件ID下载（当前固定 image/jpeg）

关系链-好友
- POST /api/relation/friend/request：发起好友申请
- POST /api/relation/friend/respond：处理好友申请（accept/reject）
- GET /api/relation/friend/requests：我的申请（我发起/我收到）
- GET /api/relation/friends：我的好友列表（用户ID集合）
- DELETE /api/relation/friend/{user_id}：解除好友关系

关系链-拉黑
- POST /api/relation/block/{user_id}：拉黑用户
- DELETE /api/relation/block/{user_id}：取消拉黑
- GET /api/relation/blocks：我的黑名单列表

关系链-关注
- POST /api/relation/follow/{user_id}：关注用户（唯一索引去重）
- DELETE /api/relation/follow/{user_id}：取消关注
- GET /api/relation/follow/status/{user_id}：我是否关注目标用户
- GET /api/relation/followers：我的粉丝用户ID列表
- GET /api/relation/following：我关注的用户ID列表

群组
- POST /api/group：创建群组（当前用户为群主）
- POST /api/group/{group_id}/members：添加群成员（示例仅校验群主）
- DELETE /api/group/{group_id}/members/{user_id}：移除成员/退群
- GET /api/group/my：我加入的群
- GET /api/group/{group_id}：群详情与成员列表
(校验；邀请同意；加群审批)

消息
- POST /api/message/send：统一发消息（dm/group/room），使用 counters 自增 seq
- GET /api/message/history：查询历史消息（按 seq 游标，支持 lastSeq/limit）

房间
- POST /api/room/join：根据 recruit_id 入房（创建/复用 theater 并写 participants）
- GET /api/room/{id}/messages：房间消息列表（内部转发到 message/history，支持分页）
- POST /api/room/{id}/message：房间发消息（内部复用统一发送逻辑）

招募（Recruit）
- GET /api/recruit/list：招募列表（分页/筛选：mode/status/backstory/keyword）
- GET /api/recruit/detail/{id}：招募详情
- POST /api/recruit/create：发布招募
- DELETE /api/recruit/{id}：删除招募（仅发布者/管理员）
- POST /api/recruit/{id}/accept：接取招募并入房（返回 room_id）

戏文（Record/Cassette）
- POST /api/record/create：从若干消息生成戏文（推导 participants）
- GET /api/record/list：戏文列表（分页/关键字）
- GET /api/record/detail/{id}：戏文详情
- GET /api/record/message/{id}：戏文关联消息列表

点赞
- POST /api/like：点赞/取消点赞（record/backstory，幂等切换并更新计数）

管理
- GET /api/admin/userList：管理端用户列表（受保护，需加权限控制） 