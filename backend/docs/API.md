# 宿舍管理系统后端接口说明

基础地址：`http://localhost:8080/api`

认证方式：除登录、刷新、登出、健康检查外，其余接口都需要请求头：

```http
Authorization: Bearer <access_token>
```

通用 JSON 错误响应：

```json
{
  "error": "bad_request",
  "message": "错误详情"
}
```

常见状态码：

- `200`：查询成功
- `201`：创建成功
- `204`：更新成功，无响应体
- `400`：参数错误
- `401`：未登录或 token 无效
- `403`：无权限或资源不属于当前用户
- `404`：资源不存在或当前状态不允许操作
- `409`：业务冲突

## 角色

系统角色字段：

- `student`：学生
- `repair_staff`：维修人员
- `cleaning_staff`：保洁人员
- `dormitory_manager`：宿舍管理员
- `system_admin`：系统管理员

## 认证

### 登录

`POST /auth/login`

请求：

```json
{
  "username": "student001",
  "password": "123456"
}
```

响应：

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### 刷新 Token

`POST /auth/refresh`

请求：

```json
{
  "refresh_token": "..."
}
```

响应同登录。刷新成功后旧 refresh token 会失效，请前端替换本地保存的 token。

### 登出

`POST /auth/logout`

请求：

```json
{
  "refresh_token": "..."
}
```

成功响应：`204`

### 健康检查

`GET /health`

响应：

```json
{
  "status": "ok"
}
```

## 通用查询

### 当前用户信息

`GET /students/me`

权限：任意已登录用户

响应字段来自 `users`，不包含密码哈希。学生若已分配床位，会额外返回当前宿舍/床位信息；所有用户都会返回 `has_survey` 和 `has_bed`，便于前端控制入口显示。

补充字段：

- `has_survey`：是否提交过生活习惯调查
- `has_bed`：是否拥有已占用床位
- `building_id`
- `building_name`
- `room_id`
- `room_number`
- `bed_id`
- `bed_label`
- `avatar_attachment_id`：当前头像附件 ID，若未上传头像则不存在

### 更新当前用户信息

`PUT /students/me`

权限：任意已登录用户

当前用于修改本人电话。请求：

```json
{
  "phone": "13800000000"
}
```

响应：更新后的当前用户信息，字段同 `GET /students/me`。

### 楼栋列表

`GET /buildings`

权限：任意已登录用户

### 房间水电余额

`GET /rooms/{id}/balance`

权限：任意已登录用户

响应示例：

```json
{
  "id": 1,
  "building_id": 1,
  "room_number": "201",
  "floor": 2,
  "total_beds": 4,
  "water_balance": 10,
  "electricity_balance": 20
}
```

### 空闲床位

`GET /beds/available?building_id=1&floor=2`

权限：任意已登录用户

查询参数均可选：

- `building_id`
- `floor`

## 学生端

以下接口权限：`student`

### 提交生活习惯调查

`POST /students/me/survey`

请求：

```json
{
  "sleep_time": "23:30:00",
  "smoking": 0,
  "snoring": 0,
  "study_habit": "晚上学习",
  "remarks": "希望安静"
}
```

### 获取最新生活习惯调查

`GET /students/me/survey`

### 我的申请总览

`GET /students/me/requests`

返回维修、保洁、离校、换寝、晚归、校外居住、分配申请和缴费记录的统一列表。

### 我的舍友

`GET /students/me/roommates`

### 创建新生分配申请

`POST /allocations`

请求体为空。学生必须先提交生活习惯调查，否则返回 `400`。

后端会根据最新生活习惯调查推荐床位并创建 pending 申请。推荐规则：

- 排除已占用床位。
- 排除已被其他 pending 分配申请推荐的床位。
- 使用事务锁 `FOR UPDATE SKIP LOCKED` 避免并发重复推荐。
- 优先匹配同房间已入住学生的生活习惯：
  - 吸烟习惯一致加分，不一致强扣分。
  - 打鼾习惯一致加分，不一致扣分。
  - 就寝时间越接近分数越高。
  - 学习习惯关键词有交集加分。
- 分数相同时，优先选择已有兼容室友的房间，再按楼栋、楼层、房号、床位排序。

响应：

```json
{
  "id": 1
}
```

### 提交离校申请

`POST /leaves`

请求：

```json
{
  "type": "normal",
  "destination": "上海",
  "emergency_contact": "张三 13800000000",
  "return_time": "2026-05-10T10:00:00Z",
  "reason": "回家"
}
```

`type` 可选：`normal`、`holiday`。不传默认为 `normal`。

### 提交晚归记录

`POST /late-returns`

请求：

```json
{
  "return_date": "2026-05-08",
  "reason": "实验结束较晚"
}
```

### 提交换寝申请

`POST /room-changes`

请求，可不指定目标床位：

```json
{
  "reason": "作息不合"
}
```

请求，指定目标床位：

```json
{
  "target_room_id": 2,
  "target_bed_id": 8,
  "reason": "与同专业同学同住"
}
```

若指定目标，`target_room_id` 和 `target_bed_id` 必须同时提供。

### 提交校外居住申请

`POST /off-campus`

请求：

```json
{
  "retain_bed": 0,
  "reason": "家庭原因",
  "destination": "校外地址"
}
```

`retain_bed`：`0` 不保留床位，`1` 保留床位。

### 提交维修工单

`POST /repairs`

请求：

```json
{
  "room_id": 1,
  "description": "水龙头漏水"
}
```

后端会校验当前学生是否入住该房间。

### 提交保洁工单

`POST /cleanings`

请求：

```json
{
  "building_id": 1,
  "location_desc": "二楼公共洗手间"
}
```

### 水电缴费

`POST /payments`

请求：

```json
{
  "room_id": 1,
  "amount": 50,
  "payment_type": "both"
}
```

`payment_type` 可选：`water`、`electricity`、`both`。后端会校验当前学生是否入住该房间。

## 维修人员

以下接口权限：`repair_staff`

### 维修工单列表

`GET /repairs`

权限：`student`、`repair_staff`、`dormitory_manager`

- 学生：返回自己提交的全部维修工单。
- 维修人员：返回自己接单过的全部维修工单。
- 宿舍管理员：返回全部维修工单。

返回包含学生、房间、维修人员、审核人、维修说明、审核意见和各状态时间；照片通过附件接口按 `owner_type=repair&owner_id=工单ID&category=after` 查询。

### 待处理维修工单

`GET /repairs/pending`

权限：`repair_staff`、`dormitory_manager`

宿舍管理员使用该接口查看 `pending`、`accepted`、`repaired` 状态的维修工单。返回包含学生、房间、维修人员、审核人、维修说明、审核意见和各状态时间；照片通过附件接口按 `owner_type=repair&owner_id=工单ID&category=after` 查询。

### 接单

`PUT /repairs/{id}/accept`

成功响应：`204`

### 维修完成

`PUT /repairs/{id}/repair`

请求：

```json
{
  "repair_description": "已更换水龙头阀芯"
}
```

成功响应：`204`

调用维修完成接口前，必须先通过附件接口上传至少一张维修后照片，使用：

- `owner_type=repair`
- `owner_id=<维修工单ID>`
- `category=after`

## 保洁人员

以下接口权限：`cleaning_staff`

### 保洁工单列表

`GET /cleanings`

权限：`student`、`cleaning_staff`、`dormitory_manager`

- 学生：返回自己提交的全部保洁工单。
- 保洁人员：返回自己接单过的全部保洁工单。
- 宿舍管理员：返回全部保洁工单。

返回包含学生、楼栋、保洁人员、审核人、审核意见和各状态时间；照片通过附件接口按 `owner_type=cleaning&owner_id=工单ID&category=before/after` 查询。

### 待处理保洁工单

`GET /cleanings/pending`

权限：`cleaning_staff`、`dormitory_manager`

宿舍管理员使用该接口查看 `pending`、`accepted`、`cleaned` 状态的保洁工单。返回包含学生、楼栋、保洁人员、审核人、审核意见和各状态时间；照片通过附件接口按 `owner_type=cleaning&owner_id=工单ID&category=before/after` 查询。

### 接单

`PUT /cleanings/{id}/accept`

### 清洁完成

`PUT /cleanings/{id}/clean`

成功响应：`204`

调用清洁完成接口前，必须先通过附件接口上传至少一张保洁后照片，使用：

- `owner_type=cleaning`
- `owner_id=<保洁工单ID>`
- `category=after`

## 宿舍管理员

以下接口权限：`dormitory_manager`

### 待审批离校申请

`GET /leaves/pending`

返回 `status=pending` 的离校/节假日离校申请，包含学生姓名、目的地、紧急联系人、返校时间和原因。

### 审批离校申请

`PUT /leaves/{id}/review`

请求：

```json
{
  "status": "approved"
}
```

`status` 可选：`approved`、`rejected`。

### 待审批晚归记录

`GET /late-returns/pending`

返回 `status=pending` 的晚归记录，包含学生姓名、晚归日期和原因。

### 审批晚归记录

`PUT /late-returns/{id}/review`

请求同上。

### 待审批换寝申请

`GET /room-changes/pending`

返回 `status=pending` 的换寝申请，包含原床位、目标床位或推荐床位信息。

### 审批换寝申请

`PUT /room-changes/{id}/review`

请求同上。批准后后端会在事务中释放旧床位并占用目标或推荐床位。

### 待审批校外居住申请

`GET /off-campus/pending`

返回 `status=pending` 的校外居住申请，包含学生姓名、是否保留床位、原因和居住地址。

### 审批校外居住申请

`PUT /off-campus/{id}/review`

请求同上。若批准且申请不保留床位，后端会释放学生当前床位。

### 审核维修工单

列表接口：`GET /repairs/pending`，其中 `status=repaired` 的记录可审核。

`PUT /repairs/{id}/review`

请求：

```json
{
  "status": "completed",
  "comment": "维修合格"
}
```

`status` 可选：`completed`、`rejected`。

### 审核保洁工单

列表接口：`GET /cleanings/pending`，其中 `status=cleaned` 的记录可审核。

`PUT /cleanings/{id}/review`

请求同维修审核。

### 楼栋统计

`GET /dashboard/summary`

### 低余额房间

`GET /dashboard/low-balance`

## 系统管理员

以下接口权限：`system_admin`

### 创建用户

`POST /users`

请求：

```json
{
  "username": "student001",
  "password": "123456",
  "role": "student",
  "name": "张三",
  "phone": "13800000000"
}
```

### 待审批分配申请

`GET /allocations/pending`

### 审批分配申请

`PUT /allocations/{id}/review`

请求：

```json
{
  "status": "approved"
}
```

`status` 可选：`approved`、`rejected`。批准后后端会在事务中占用推荐床位。

### 统计接口

系统管理员也可访问：

- `GET /dashboard/summary`
- `GET /dashboard/low-balance`

## 通知

权限：任意已登录用户

### 通知列表

`GET /notifications`

只返回当前用户的通知。

### 标记已读

`PUT /notifications/{id}/read`

成功响应：`204`

## 附件

权限：任意已登录用户，但后端会按业务对象做资源权限校验。

支持图片类型：

- `image/jpeg`
- `image/png`

大小限制：默认 `5MB`，由 `upload.max_size` 配置。

### 上传附件

`POST /attachments`

请求格式：`multipart/form-data`

字段：

- `file`：图片文件，必填
- `owner_type`：`user_avatar`、`repair`、`cleaning`
- `owner_id`：业务对象 ID
- `category`：`avatar`、`before`、`after`
- `sort_order`：排序，可选

示例：

```bash
curl -X POST http://localhost:8080/api/attachments \
  -H "Authorization: Bearer <access_token>" \
  -F "owner_type=user_avatar" \
  -F "owner_id=1" \
  -F "category=avatar" \
  -F "file=@avatar.png"
```

合法组合：

| 业务对象 | owner_type | category |
| --- | --- | --- |
| 用户头像 | `user_avatar` | `avatar` |
| 维修后照片 | `repair` | `after` |
| 保洁前照片 | `cleaning` | `before` |
| 保洁后照片 | `cleaning` | `after` |

权限规则：

- 用户只能上传自己的头像。
- 维修人员只能给自己接单的维修工单上传 `after` 照片。
- 学生只能给自己创建的保洁工单上传 `before` 照片。
- 保洁人员只能给自己接单的保洁工单上传 `after` 照片。
- 宿舍管理员和系统管理员可访问管理场景所需附件。

### 附件元数据

`GET /attachments?owner_type=cleaning&owner_id=1&category=before`

`category` 可选。

响应不包含二进制内容，适合列表展示。

### 下载附件

`GET /attachments/{id}`

响应为图片二进制流，`Content-Type` 为图片 MIME 类型。

## 前端对接建议

- 登录后保存 `access_token` 和 `refresh_token`。
- API 返回 `401` 时，先调用 `/auth/refresh`；刷新成功后重试原请求。
- 刷新成功必须替换本地旧 `refresh_token`，旧 token 已被撤销。
- 调用 `/auth/logout` 后清空本地 token。
- 文件上传不要设置 JSON `Content-Type`，让浏览器自动生成 `multipart/form-data` boundary。
- 日期时间字段使用 ISO 8601，例如 `2026-05-10T10:00:00Z`。
