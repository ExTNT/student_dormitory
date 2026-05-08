# 宿舍管理系统后端设计方案（Go）

> **版本**：v1.0  
> **关联文档**：数据库设计文档 v2.4  
> **技术栈**：Go 1.22+, PostgreSQL 12+, JWT  
> **设计原则**：RESTful、分层架构、事务性保障、安全可控

---

## 1. 项目概述

本系统为 B/S 架构的宿舍管理系统，后端采用 Go 语言实现，提供 RESTful API 服务。核心功能覆盖新生分配、离校/晚归/换寝/校外居住申请、维修保洁工单、水电缴费与低余额提醒、宿舍及舍友查询等 14 项业务。图片以二进制方式存入 PostgreSQL，通过专用附件表管理。

### 目标
- 高可用、可扩展的 API 服务
- 清晰的模块划分与低耦合设计
- 通过数据库触发器与事务协同保证数据一致性
- 安全认证与细粒度权限控制

---

## 2. 技术选型

| 组件 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.22+ | 高性能、强类型、并发支持 |
| Web 框架 | Gin | 轻量、高性能、中间件生态丰富 |
| 数据库驱动 | pgx v5 | 纯 Go 实现，高性能，原生支持 PostgreSQL |
| 数据访问 | sqlx + pgx 适配 | 编写原生 SQL，利用 sqlx 的结构体映射 |
| 认证 | JWT（golang-jwt）+ bcrypt | 无状态认证，密码 hash 存储 |
| 参数校验 | go-playground/validator | 结构体标签验证 |
| 配置管理 | Viper | 多源配置（文件、环境变量） |
| 日志 | Zap | 高性能结构化日志 |
| 图片处理 | 标准库 `image` / `disintegration/imaging`（可选） | 用于可能需要的缩略图生成或格式校验 |
| 通知（实时） | WebSocket（gorilla/websocket） | 低余额、工单状态实时推送（可选） |

---

## 3. 系统架构

采用经典分层架构（Controller → Service → Repository），结合中间件实现横切关注点。

```
┌─────────────────────────────────────────────┐
│                 HTTP Clients                │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────┴──────────────────────────┐
│           Gin Router & Middleware           │
│  ┌───────────────────────────────────────┐  │
│  │    Auth / CORS / Logger / Recovery    │  │
│  ├───────────────────┬───────────────────┤  │
│  │  Handler Layer    │   WebSocket Hub   │  │
│  ├───────────────────┴───────────────────┤  │
│  │         Service / Business Layer      │  │
│  ├───────────────────────────────────────┤  │
│  │         Repository / Data Access      │  │
│  └───────────────┬───────────────────────┘  │
└──────────────────┼──────────────────────────┘
                   │
┌──────────────────┴──────────────────────────┐
│               PostgreSQL DB                 │
│  ┌───────────────────────────────────────┐  │
│  │ Tables, Views, Triggers, Functions    │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

- **Handler**：解析请求、校验、调用 Service、封装响应。
- **Service**：核心业务逻辑，跨多个 Repository 操作，编排事务。
- **Repository**：封装数据库查询，使用 sqlx + pgx，返回数据对象或错误。
- **Middleware**：身份认证、权限校验、请求日志等。

---

## 4. 模块划分

| 模块 | 职责 | 主要 API 路径 |
|------|------|--------------|
| auth | 登录、登出、Token 刷新 | `/api/auth/*` |
| user | 用户管理、头像上传/获取 | `/api/users/*` |
| building | 楼栋信息查询 | `/api/buildings` |
| room | 房间信息、水电余额查询 | `/api/rooms/*` |
| bed | 床位查询、入住分配 | `/api/beds/*` |
| allocation | 新生宿舍分配请求 | `/api/allocations` |
| leave | 离校/节假日申请 | `/api/leaves` |
| late | 晚归记录 | `/api/late-returns` |
| room-change | 换寝申请 | `/api/room-changes` |
| off-campus | 校外居住申请 | `/api/off-campus` |
| repair | 维修工单 | `/api/repairs` |
| cleaning | 保洁工单 | `/api/v1/cleanings` |
| payment | 水电缴费 | `/api/payments` |
| notification | 通知查询、标记已读 | `/api/notifications` |
| attachment | 附件上传、下载 | `/api/attachments` |
| dashboard | 统计视图（宿舍楼汇总等） | `/api/dashboard/*` |
| roommate | 舍友信息 | `/api/roommates` |
| survey | 生活习惯调查 | `/api/surveys` |

每个模块对应独立的 Service 和 Repository 包。

---

## 5. API 设计（RESTful）

### 5.1 通用规范
- 基础路径：`/api`
- 请求/响应格式：JSON
- 认证：Header `Authorization: Bearer <token>`
- 分页：查询参数 `page`、`page_size`，响应包含 `total`、`page`、`page_size`。
- 时间格式：ISO 8601 UTC
- 状态码：标准 HTTP 状态码
- 错误响应体：
  ```json
  {
    "error": "short description",
    "message": "detailed message"
  }
  ```

### 5.2 核心端点示例

#### 认证
- `POST /auth/login` – 登录返回 JWT
- `POST /auth/refresh` – 刷新 Token

#### 学生端
- `GET /students/me` – 个人信息
- `GET /students/me/survey` – 我的生活习惯调查
- `POST /students/me/survey` – 提交/更新调查
- `GET /students/me/requests` – 我的所有申请（调用 `v_my_requests`）
- `GET /students/me/roommates` – 查看舍友（`v_student_roommates`）
- `POST /leaves` – 提交离校申请
- `POST /late-returns` – 提交晚归记录
- `POST /room-changes` – 提交换寝申请
- `POST /off-campus` – 提交校外居住申请
- `POST /repairs` – 提交维修工单
- `POST /cleanings` – 提交保洁申请
- `POST /payments` – 缴纳水电费
- `GET /rooms/{id}/balance` – 查询宿舍水电余额
- `GET /beds/available` – 查询空闲床位（支持按楼栋、楼层过滤）

#### 维修人员
- `GET /repairs/pending` – 待处理维修工单列表
- `PUT /repairs/{id}/accept` – 接单
- `PUT /repairs/{id}/repair` – 维修完成（提交维修说明及上传照片）
- `POST /repairs/{id}/photos` – 上传维修后照片（亦可合并到维修完成接口）

#### 保洁人员
- `GET /cleanings/pending` – 待处理保洁工单
- `PUT /cleanings/{id}/accept` – 接单
- `PUT /cleanings/{id}/clean` – 清洁完成（含上传 after 照片）

#### 宿舍管理人员
- `GET /leaves/pending` – 待审批离校申请
- `PUT /leaves/{id}/approve` – 审批离校
- 同理处理晚归、换寝、校外居住、维修/保洁审核等

#### 系统管理员
- `GET /allocations/pending` – 待审批分配请求
- `PUT /allocations/{id}/approve` – 批准分配（自动占用床位）
- 用户管理 CRUD（创建、禁用等）

#### 附件
- `POST /attachments` – 上传图片（关联 owner_type, owner_id, category）
- `GET /attachments/{id}` – 获取原始图片（返回二进制流，设置 Content-Type）
- `GET /attachments/{id}/thumbnail` – 获取缩略图（可选）

#### 仪表盘/统计
- `GET /dashboard/summary` – 楼栋汇总视图（`v_dormitory_summary`）
- `GET /dashboard/low-balance` – 低余额宿舍列表（`v_low_balance_rooms`）

详细 API 文档可通过 Swagger 自动生成（swaggo/gin-swagger），便于前端对接。

---

## 6. 数据库交互层

### 6.1 数据库连接
使用 `pgxpool` 创建连接池，配置最大连接数、超时等。结合 `sqlx` 扩展方法，使用原生 SQL 并方便地映射结构体。

```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Open("pgx", dsn)
    // 配置连接池参数
    return db, err
}
```

### 6.2 查询方式
- 简单查询：`db.Get/Select` 配合 `db:"column"` 标签。
- 多表关联：命名查询 `db.NamedExec` 或手动 SQL 拼接（防注入）。
- 统计视图查询：直接映射结构体。
  ```go
  type DormitorySummary struct {
      BuildingID   int    `db:"building_id"`
      BuildingName string `db:"building_name"`
      TotalRooms   int    `db:"total_rooms"`
      ...
  }
  summaries := []DormitorySummary{}
  err := db.Select(&summaries, "SELECT * FROM v_dormitory_summary")
  ```

### 6.3 事务管理
关键业务（分配床位、缴费、状态流转）必须在事务中执行。Service 层控制事务：

```go
func (s *AllocationService) Approve(ctx context.Context, reqID int64, adminID int64) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    // 1. 更新 allocation_request 状态
    // 2. 更新 beds 表为 occupied
    // 3. 发送通知（可选）
    return tx.Commit()
}
```

### 6.4 触发器协同
数据库已定义：
- 时间戳触发器（自动设置 `accepted_at`、`repaired_at` 等）
- 床位一致性触发器（`trg_bed_consistency`）
- 低余额通知触发器（`trg_low_balance_notification`）
- 缴费后更新余额触发器（`trg_apply_utility_payment`）

后端职责：
- 执行 UPDATE 语句即可，时间戳由数据库自动填写。
- 通知由触发器写入，后端提供查询接口；若需实时推送，可监听 PostgreSQL `LISTEN/NOTIFY` 或结合 WebSocket 轮询。

---

## 7. 认证与授权

### 7.1 JWT 设计
- 登录成功生成 JWT，包含 UserID、Role、过期时间。
- 算法：HS256（密钥可配置）。
- 有效期：Access Token 15分钟，Refresh Token 7天。
- Refresh Token 可存储于 Redis 或数据库，支持撤销。

### 7.2 权限中间件
自定义 Gin 中间件 `AuthRequired(roles ...string)`：
1. 解析 Header 中的 Token 并验证。
2. 检查角色是否在允许列表中。
3. 将 UserID、Role 注入 `context`，后续 Handler 可通过 `c.Get("user")` 获取。
4. 资源级权限（如只能操作自己的申请）在 Service 层二次校验。

### 7.3 密码存储
使用 `golang.org/x/crypto/bcrypt`，cost 设置为 12，存储哈希，登录时比对。

---

## 8. 关键业务流程实现

### 8.1 新生宿舍分配
1. 学生提交生活习惯调查。
2. 系统根据调查结果和空闲床位（`v_available_beds`）生成推荐（算法可置入 Service）。
3. 创建 `allocation_requests`（pending）。
4. 管理员审批，事务中：更新请求状态为 approved → 占用床位（status='occupied', student_id, occupied_since）→ 触发时间戳。

### 8.2 换寝流程
1. 学生提交 `room_change_requests`，可指定目标或留空。
2. 若留空，系统通过 `v_available_beds` 推荐同楼栋床位，填入 `recommended_bed_id`。
3. 管理员审批：数据库触发器校验目标床位合法性（同楼栋、空闲）。事务中释放原床位、占用新床位，更新申请状态。

### 8.3 维修/保洁工单状态机
- 状态流转：`pending → accepted → repaired/cleaned → completed/rejected`
- Handler 调用 Service，Service 更新状态字段；数据库自动记录时间戳。
- 上传照片：调用附件接口，插入 `attachments` 关联工单 ID。

### 8.4 水电缴费与余额提醒
- 缴费：写入 `utility_payments`，触发器自动增加 `rooms` 余额。
- 低余额：`rooms` 表的触发器检测到余额<5时，向该宿舍所有学生发送通知（写入 `notifications`）。
- 后端提供未读通知列表，可配合 WebSocket 推送。

### 8.5 校外居住申请
- 提交时选择 `retain_bed`。
- 审批通过且不保留床位：事务中释放对应床位。触发器保证床位状态一致性。

所有涉及床位变更的操作均受 `trg_bed_consistency` 保护。

---

## 9. 图片处理方案

### 9.1 上传
- 客户端通过 `multipart/form-data` 提交，附带 `owner_type`, `owner_id`, `category`, `sort_order`。
- 后端校验文件大小（例如 ≤5MB）、MIME 类型（`image/jpeg`, `image/png`）。
- 读取文件为 `[]byte`，插入 `attachments` 表。头像更新时，先删旧再插新（事务）。
- 可选生成缩略图（存储为单独字段或记录，简化可暂不做）。

### 9.2 下载
- `GET /attachments/{id}` 查询 `file_data` 为 `[]byte`。
- 设置 `Content-Type` 为记录的 `content_type`，`Content-Length` 为数据长度，直接写入响应体。
- 使用 `c.Data(http.StatusOK, contentType, data)` 返回。

### 9.3 性能优化
- 列表查询避免加载 `file_data`，使用 `v_attachment_metadata` 或仅查询元数据列。

---

## 10. 通知与实时推送

- 通知存储于 `notifications` 表，触发器写入低余额、工单状态变化消息。
- API：`GET /notifications`（按时间倒序）、`PUT /notifications/{id}/read`。
- 实时推送：建立 WebSocket 连接（认证后），后端可从 `notifications` 表扫描新增记录或监听 PostgreSQL `NOTIFY` 频道推送至对应用户。
- 工单状态变化：Service 层在更新状态后，通过 WebSocket Hub 发布事件。

---

## 11. 配置管理

使用 Viper 加载 `config.yaml` 和环境变量，示例：
```yaml
server:
  port: 8080
  mode: release
database:
  host: localhost
  port: 5432
  user: turing
  password: 10928
  name: student_dormitory
  max_open_conns: 25
jwt:
  secret: your-secret
  access_expiry: 15m
  refresh_expiry: 168h
upload:
  max_size: 5242880  # 5MB
```

在 `main.go` 中初始化配置，并注入到各层。

---

## 12. 项目结构建议

```
/backend
├── cmd
│   └── server
│       └── main.go              # 入口
├── config
│   ├── config.go
│   └── config.yaml
├── internal
│   ├── handler                   # HTTP 处理函数（按模块）
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── repair.go
│   │   └── ...
│   ├── service                  # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── allocation_service.go
│   │   └── ...
│   ├── repository               # 数据访问层（原生SQL）
│   │   ├── user_repo.go
│   │   ├── bed_repo.go
│   │   └── ...
│   ├── model                    # 数据库实体模型
│   │   ├── user.go
│   │   ├── request.go
│   │   └── ...
│   ├── middleware                # 中间件（认证、日志、CORS）
│   ├── dto                      # 请求/响应 DTO
│   ├── errs                     # 自定义错误
│   └── utils                    # 工具函数
├── go.mod
└── go.sum
```

---

## 13. 安全注意事项

- 所有 SQL 使用参数化查询，杜绝拼接。
- JWT 密钥不硬编码，通过环境变量提供。
- 密码 bcrypt 加密，cost=12。
- 附件上传校验 MIME 类型，防恶意文件。
- 角色权限强制中间件校验，资源归属在 Service 层二次检查。
- 生产环境强制 HTTPS。
- 设置请求限流中间件。
- 日志中过滤密码、Token 等敏感字段。

---

## 14. 总结

本设计方案充分利用 Go 语言的并发能力、Gin 的生态以及 PostgreSQL 的高级特性（触发器、视图），构建一个健壮、安全、可维护的宿舍管理系统后端。通过清晰的分层架构、事务管理和数据库触发器的协同，保证了业务数据的一致性，同时减少了应用层冗余代码。结构化的模块划分和详细的 API 设计为后续开发提供了明确的指导。