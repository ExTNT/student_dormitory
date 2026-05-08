# 宿舍管理系统数据库设计文档

> **版本**：v2.4  
> **主要更新**：新增数据库触发器建议，覆盖自动时间戳、低余额通知、床位一致性约束等场景  
> **图片存储**：数据库二进制存储（BYTEA），统一附件表  
> **适用范围**：覆盖全部14项业务事件，提供便捷查询视图与自动化触发器

---

## 1. 引言

### 1.1 编写目的
本文档为宿舍管理系统提供完整的数据库设计方案，包括实体定义、表结构、字段约束、图片存储方案、业务视图、数据库触发器以及实施要点，指导后续开发与维护。

### 1.2 系统角色与职责

| 角色 | 职责概要 |
|------|----------|
| 学生 | 填写生活习惯问卷、申请换寝/离校/晚归/维修/保洁/校外居住、缴纳水电费、查询工单状态及舍友 |
| 宿舍维修人员 | 接收维修工单，处理并上传维修后照片 |
| 宿舍保洁人员 | 接收保洁工单，清洁并上传清洁后照片 |
| 宿舍管理人员 | 审批离校、换寝、晚归、维修工单、保洁工单、校外居住申请 |
| 系统管理员 | 审批宿舍分配请求，管理用户账号及基础配置 |

### 1.3 图片存储策略
所有图片（用户头像、维修后照片、保洁前后照片）以二进制格式存入 `attachments` 表（BYTEA 类型），通过 `owner_type`、`owner_id`、`category` 与业务记录关联。数据库中不存储任何外部 URL 或文件路径。

---

## 2. 概念结构设计

### 2.1 核心实体与关系
- **用户** 与 **床位**：一个学生角色占用零个或一个床位；床位隶属于某个宿舍。
- **宿舍楼** 包含多个 **宿舍**；宿舍包含多个 **床位**。
- **学生** 可以提交多份生活习惯调查、离校申请、水电缴费、维修申请、保洁申请、晚归记录、换寝申请、校外居住申请。
- **维修工单** 由一名 **维修人员** 处理，可包含多张维修后照片。
- **保洁工单** 由一名 **保洁人员** 处理，可包含多张保洁前/后照片。
- **宿舍管理人员** 审批多种申请。
- **系统管理员** 审批宿舍分配。
- **附件** 通过多对一方式与不同业务实体关联。

---

## 3. 表结构设计

### 3.1 表清单

| 序号 | 表名 | 说明 |
|------|------|------|
| 1 | `users` | 系统用户（所有角色） |
| 2 | `buildings` | 宿舍楼 |
| 3 | `rooms` | 宿舍房间（含水/电余额） |
| 4 | `beds` | 床位（占用状态与入住学生） |
| 5 | `lifestyle_surveys` | 新生生活习惯调查 |
| 6 | `allocation_requests` | 宿舍分配请求（推荐→审批） |
| 7 | `leave_applications` | 离校/节假日申请 |
| 8 | `utility_payments` | 水电缴费记录 |
| 9 | `repair_requests` | 维修工单 |
| 10 | `cleaning_requests` | 保洁工单 |
| 11 | `late_return_records` | 晚归记录 |
| 12 | `room_change_requests` | 换寝申请 |
| 13 | `notifications` | 通知（低余额提醒等） |
| 14 | `attachments` | 附件表（所有图片二进制） |
| 15 | `off_campus_living_applications` | 校外居住申请 |

---

### 3.2 详细字段说明

#### 3.2.1 `users`（用户表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 用户ID |
| `username` | VARCHAR(32) | 非空，唯一 | 登录账号 |
| `password_hash` | VARCHAR(128) | 非空 | 加密密码（bcrypt） |
| `role` | VARCHAR(20) | 非空 | 角色：student, repair_staff, cleaning_staff, dormitory_manager, system_admin |
| `name` | VARCHAR(32) | 非空 | 真实姓名 |
| `phone` | VARCHAR(20) | 可空 | 联系电话 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 创建时间 |

- **头像** 存储于附件表，通过 `owner_type='user_avatar'`, `category='avatar'` 关联。

#### 3.2.2 `buildings`（宿舍楼表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | INT | 主键，自增 | 楼栋ID |
| `name` | VARCHAR(64) | 非空 | 楼栋名称 |
| `location` | VARCHAR(255) | 可空 | 位置描述 |

#### 3.2.3 `rooms`（宿舍房间表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | INT | 主键，自增 | 房间ID |
| `building_id` | INT | 非空，外键→buildings.id | 所属楼栋 |
| `room_number` | VARCHAR(16) | 非空 | 房间号（如201） |
| `floor` | SMALLINT | 非空 | 楼层 |
| `total_beds` | SMALLINT | 非空 | 床位总数 |
| `water_balance` | NUMERIC(8,2) | 默认0.00 | 水费余额（元） |
| `electricity_balance` | NUMERIC(8,2) | 默认0.00 | 电费余额（元） |

- 组合唯一：(building_id, room_number)

#### 3.2.4 `beds`（床位表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | INT | 主键，自增 | 床位ID |
| `room_id` | INT | 非空，外键→rooms.id | 所属房间 |
| `bed_label` | VARCHAR(8) | 非空 | 床位标签（如A、1床） |
| `status` | VARCHAR(10) | 非空，默认'available' | 状态：available(空闲), occupied(占用) |
| `student_id` | BIGINT | 外键→users.id，可空 | 入住学生ID |
| `occupied_since` | DATE | 可空 | 入住日期 |

- 组合唯一：(room_id, bed_label)

#### 3.2.5 `lifestyle_surveys`（生活习惯调查表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 问卷ID |
| `student_id` | BIGINT | 非空，外键→users.id | 填表新生 |
| `sleep_time` | TIME | 可空 | 就寝时间 |
| `smoking` | SMALLINT | 默认0，CHECK(0,1) | 是否吸烟 |
| `snoring` | SMALLINT | 默认0，CHECK(0,1) | 是否打鼾 |
| `study_habit` | VARCHAR(255) | 可空 | 学习习惯 |
| `remarks` | TEXT | 可空 | 备注 |
| `submitted_at` | TIMESTAMPTZ | 默认当前时间 | 提交时间 |

#### 3.2.6 `allocation_requests`（宿舍分配请求表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 请求ID |
| `student_id` | BIGINT | 非空，外键→users.id | 新生ID |
| `recommended_room_id` | INT | 非空，外键→rooms.id | 推荐房间 |
| `recommended_bed_id` | INT | 非空，外键→beds.id | 推荐床位 |
| `status` | VARCHAR(10) | 非空，默认'pending' | 状态：pending, approved, rejected |
| `admin_id` | BIGINT | 外键→users.id，可空 | 审批管理员 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 创建时间 |
| `resolved_at` | TIMESTAMPTZ | 可空 | 处理时间 |

#### 3.2.7 `leave_applications`（离校申请表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 申请ID |
| `student_id` | BIGINT | 非空，外键→users.id | 学生 |
| `type` | VARCHAR(10) | 非空，默认'normal' | 类型：normal, holiday |
| `destination` | VARCHAR(128) | 非空 | 目的地 |
| `emergency_contact` | VARCHAR(64) | 非空 | 紧急联系人及电话 |
| `return_time` | TIMESTAMPTZ | 非空 | 预计返校时间 |
| `reason` | TEXT | 非空 | 原因 |
| `status` | VARCHAR(10) | 非空，默认'pending' | 状态：pending, approved, rejected |
| `manager_id` | BIGINT | 外键→users.id，可空 | 审批的宿舍管理人员 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `resolved_at` | TIMESTAMPTZ | 可空 | 审批时间 |

#### 3.2.8 `utility_payments`（水电缴费记录表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 缴费ID |
| `room_id` | INT | 非空，外键→rooms.id | 缴费宿舍 |
| `payer_id` | BIGINT | 非空，外键→users.id | 付款学生 |
| `amount` | NUMERIC(8,2) | 非空 | 金额（元） |
| `payment_type` | VARCHAR(12) | 非空 | 类型：water, electricity, both |
| `paid_at` | TIMESTAMPTZ | 默认当前时间 | 缴费时间 |

#### 3.2.9 `repair_requests`（维修工单表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 工单ID |
| `student_id` | BIGINT | 非空，外键→users.id | 申请学生 |
| `room_id` | INT | 非空，外键→rooms.id | 所在宿舍 |
| `description` | TEXT | 非空 | 维修问题描述 |
| `status` | VARCHAR(12) | 非空，默认'pending' | 状态：pending→accepted→repaired→completed，或rejected |
| `repair_staff_id` | BIGINT | 外键→users.id，可空 | 维修人员 |
| `repair_description` | TEXT | 可空 | 维修说明 |
| `reviewer_id` | BIGINT | 外键→users.id，可空 | 审核人（宿舍管理人员） |
| `review_comment` | TEXT | 可空 | 审核意见 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `accepted_at` | TIMESTAMPTZ | 可空 | 接单时间 |
| `repaired_at` | TIMESTAMPTZ | 可空 | 维修完成时间（上传照片） |
| `reviewed_at` | TIMESTAMPTZ | 可空 | 审核时间 |

- 维修后照片存储于附件表：`owner_type='repair'`, `category='after'`。

#### 3.2.10 `cleaning_requests`（保洁工单表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 工单ID |
| `student_id` | BIGINT | 非空，外键→users.id | 申请学生 |
| `building_id` | INT | 非空，外键→buildings.id | 宿舍楼 |
| `location_desc` | VARCHAR(255) | 非空 | 公共区域位置描述 |
| `status` | VARCHAR(12) | 非空，默认'pending' | 状态：pending→accepted→cleaned→completed，或 rejected |
| `cleaner_id` | BIGINT | 外键→users.id，可空 | 保洁人员 |
| `reviewer_id` | BIGINT | 外键→users.id，可空 | 审核人 |
| `review_comment` | TEXT | 可空 | 审核意见 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `accepted_at` | TIMESTAMPTZ | 可空 | 接单时间 |
| `cleaned_at` | TIMESTAMPTZ | 可空 | 清洁完成时间 |
| `reviewed_at` | TIMESTAMPTZ | 可空 | 审核时间 |

- 保洁前照片：附件表 `category='before'`；保洁后照片：`category='after'`。

#### 3.2.11 `late_return_records`（晚归记录表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 记录ID |
| `student_id` | BIGINT | 非空，外键→users.id | 学生 |
| `return_date` | DATE | 非空 | 晚归日期 |
| `reason` | TEXT | 非空 | 原因 |
| `status` | VARCHAR(10) | 非空，默认'pending' | 状态：pending, approved, rejected |
| `manager_id` | BIGINT | 外键→users.id，可空 | 审批人 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `resolved_at` | TIMESTAMPTZ | 可空 | 审批时间 |

#### 3.2.12 `room_change_requests`（换寝申请表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 申请ID |
| `student_id` | BIGINT | 非空，外键→users.id | 学生 |
| `from_bed_id` | INT | 非空，外键→beds.id | 原床位 |
| `target_room_id` | INT | 外键→rooms.id，可空 | 学生指定的目标宿舍 |
| `target_bed_id` | INT | 外键→beds.id，可空 | 学生指定的目标床位 |
| `recommended_bed_id` | INT | 外键→beds.id，可空 | 系统推荐的床位（目标为空时） |
| `reason` | TEXT | 非空 | 换寝原因 |
| `status` | VARCHAR(10) | 非空，默认'pending' | 状态：pending, approved, rejected |
| `manager_id` | BIGINT | 外键→users.id，可空 | 审批人 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `resolved_at` | TIMESTAMPTZ | 可空 | 审批时间 |

- 约束：若指定目标，其楼栋必须与原床位一致，且目标床位必须为空闲。

#### 3.2.13 `notifications`（通知表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 通知ID |
| `recipient_id` | BIGINT | 非空，外键→users.id | 接收学生 |
| `room_id` | INT | 外键→rooms.id，可空 | 关联宿舍（水电通知时使用） |
| `message` | TEXT | 非空 | 通知内容 |
| `type` | VARCHAR(12) | 非空，默认'general' | 类型：low_balance, general |
| `is_read` | SMALLINT | 默认0，CHECK(0,1) | 是否已读 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 通知时间 |

#### 3.2.14 `attachments`（附件表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 附件ID |
| `owner_type` | VARCHAR(32) | 非空 | 业务类型：user_avatar, repair, cleaning |
| `owner_id` | BIGINT | 非空 | 业务主键ID |
| `category` | VARCHAR(32) | 非空 | 分类：avatar, before, after |
| `sort_order` | SMALLINT | 默认0 | 同组排序 |
| `file_name` | VARCHAR(255) | 可空 | 原始文件名 |
| `content_type` | VARCHAR(100) | 非空 | MIME类型（image/jpeg等） |
| `file_data` | BYTEA | 非空 | 图片二进制数据 |
| `uploaded_at` | TIMESTAMPTZ | 默认当前时间 | 上传时间 |

- 索引：(owner_type, owner_id, category)

#### 3.2.15 `off_campus_living_applications`（校外居住申请表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `id` | BIGINT | 主键，自增 | 申请ID |
| `student_id` | BIGINT | 非空，外键→users.id | 学生 |
| `retain_bed` | SMALLINT | 非空，默认0，CHECK(0,1) | 是否保留床位 |
| `reason` | TEXT | 非空 | 校外居住原因 |
| `destination` | VARCHAR(255) | 可空 | 居住地址（可选） |
| `status` | VARCHAR(10) | 非空，默认'pending' | 状态：pending, approved, rejected |
| `manager_id` | BIGINT | 外键→users.id，可空 | 审批人 |
| `review_comment` | TEXT | 可空 | 审批意见 |
| `created_at` | TIMESTAMPTZ | 默认当前时间 | 申请时间 |
| `resolved_at` | TIMESTAMPTZ | 可空 | 审批时间 |

- 审批通过且 `retain_bed=0` 时，释放对应床位。

---

## 4. 图片存储关联方式

| 业务对象 | `owner_type` | `category` | 说明 |
|----------|--------------|------------|------|
| 用户头像 | `user_avatar` | `avatar` | 每个用户最多一条，更换时先删旧后插新 |
| 维修后照片 | `repair` | `after` | 可多张，`sort_order` 排序 |
| 保洁前照片 | `cleaning` | `before` | 学生申请时上传，可多张 |
| 保洁后照片 | `cleaning` | `after` | 保洁完成后上传，可多张 |

---

## 5. 关键业务流程说明

1. **新生分配**：填写调查表→系统匹配空闲床位→生成分配请求，管理员审批后占用床位。
2. **换寝**：学生申请（可指定目标或留空）→系统推荐空位（若未指定）→宿管审批→释放旧床位、占用新床位。
3. **水电缴费与提醒**：缴费记录写入 `utility_payments`，更新 `rooms` 余额。余额<5元时，通过 `notifications` 表通知该宿舍所有学生。
4. **维修工单**：`pending` → `accepted` → 维修完成上传照片 (`repaired`) → 宿管审核 (`completed`/`rejected`)。
5. **保洁工单**：`pending`（含 before 照片）→ `accepted` → 清洁完成上传 after 照片 (`cleaned`) → 宿管审核 (`completed`/`rejected`)。
6. **晚归**：填写晚归原因→宿舍管理人员审批。
7. **离校/节假日离校**：提交详细信息（目的地、紧急联系人、返校时间）→宿管审批。
8. **校外居住**：申请并选择是否保留床位→宿管审批→若不保留床位则自动释放。

---

## 6. 推荐数据库视图

以下视图用于简化业务查询，均为只读，不存储数据。

### 6.1 宿舍楼汇总视图 (`v_dormitory_summary`)
- **用途**：统计各楼栋房间数、总床位、已占用床位、空闲床位（对应事件8）
- **数据来源**：buildings, rooms, beds
- **输出列**：`building_id`, `building_name`, `total_rooms`, `total_beds`, `occupied_beds`, `free_beds`

### 6.2 空闲床位视图 (`v_available_beds`)
- **用途**：快速查询所有空闲床位，用于分配和换寝推荐
- **数据来源**：beds (status='available'), rooms, buildings
- **输出列**：`bed_id`, `room_id`, `room_number`, `bed_label`, `building_id`, `building_name`, `floor`

### 6.3 舍友视图 (`v_student_roommates`)
- **用途**：根据当前学生ID查找同宿舍所有其他学生（事件12）。返回多行，每行一位舍友，适用于任意人数的宿舍。
- **数据来源**：beds, users
- **实现逻辑**：通过输入学生ID获取其房间，查询该房间所有占用床位，排除本人，并关联用户信息与头像附件ID。
- **输出列**：
  - `student_id`：查询者ID
  - `roommate_id`：舍友用户ID
  - `roommate_name`：舍友姓名
  - `roommate_phone`：舍友电话
  - `bed_label`：舍友床位标签
  - `avatar_attachment_id`：舍友头像的附件ID（便于前端调用图片接口）

### 6.4 低水电余额宿舍视图 (`v_low_balance_rooms`)
- **用途**：获取水费或电费低于5元的房间，用于触发通知（事件10）
- **数据来源**：rooms
- **输出列**：`room_id`, `building_id`, `room_number`, `water_balance`, `electricity_balance`
- **过滤条件**：`water_balance < 5 OR electricity_balance < 5`

### 6.5 待处理维修工单视图 (`v_pending_repairs`)
- **用途**：维修人员查看待接单/已接单未完工的工单；宿管审核待审核工单（事件5）
- **数据来源**：repair_requests, rooms, users
- **输出列**：`request_id`, `status`, `student_name`, `room_number`, `description`, `created_at`, `repair_staff_name`, `reviewer_name` 等
- **状态过滤**：通常 `status IN ('pending','accepted','repaired')`

### 6.6 待处理保洁工单视图 (`v_pending_cleanings`)
- **用途**：保洁人员及宿管查看待处理保洁工单（事件9）
- **数据来源**：cleaning_requests, buildings, users
- **输出列**：`request_id`, `status`, `student_name`, `building_name`, `location_desc`, `created_at`, `cleaner_name` 等
- **状态过滤**：通常 `status IN ('pending','accepted','cleaned')`

### 6.7 学生申请总览视图 (`v_my_requests`)
- **用途**：学生查看自己所有申请状态（事件11），统一维修、保洁、离校、换寝、晚归、校外居住等申请。
- **实现方式**：将各申请表 UNION ALL，包含类型标识。
- **输出列**：`student_id`, `request_type`（如'维修'、'保洁'）, `request_id`, `status`, `created_at`, `detail`（简要信息）

### 6.8 附件元数据视图 (`v_attachment_metadata`)
- **用途**：获取附件基本信息，不包含 `file_data` 二进制，避免大数据量传输。前端通过此视图获取附件列表，再根据附件ID调用接口下载具体图片。
- **数据来源**：attachments
- **输出列**：`id`, `owner_type`, `owner_id`, `category`, `sort_order`, `content_type`, `file_name`, `uploaded_at`

---

## 7. 推荐数据库触发器

利用触发器可以实现关键业务规则的自动化，减少应用层编码复杂度，并保证数据一致性。以下为建议创建的触发器列表，包含触发时机、用途和核心逻辑说明。

| 触发器名称 | 作用表 | 触发时机 | 功能描述 |
|-----------|--------|----------|----------|
| `trg_repair_status_timestamp` | `repair_requests` | BEFORE UPDATE OF `status` | 当工单状态更新为特定值时，自动填充对应时间戳：<br>• `status` 变为 `accepted` → 设置 `accepted_at` 为当前时间<br>• `status` 变为 `repaired` → 设置 `repaired_at`<br>• `status` 变为 `completed` 或 `rejected` → 设置 `reviewed_at`<br>若状态回退或超出范围则忽略。 |
| `trg_cleaning_status_timestamp` | `cleaning_requests` | BEFORE UPDATE OF `status` | 与维修工单类似：<br>• `status` 变为 `accepted` → 设置 `accepted_at`<br>• `status` 变为 `cleaned` → 设置 `cleaned_at`<br>• `status` 变为 `completed` 或 `rejected` → 设置 `reviewed_at` |
| `trg_bed_consistency` | `beds` | BEFORE INSERT OR UPDATE | 确保床位状态与学生ID的一致性：<br>• 若 `status = 'occupied'`，则 `student_id` 必须非空且 `occupied_since` 必须非空<br>• 若 `status = 'available'`，则 `student_id` 必须为 NULL 且 `occupied_since` 必须为 NULL<br>违反时抛出错误，防止脏数据。 |
| `trg_low_balance_notification` | `rooms` | AFTER UPDATE OF `water_balance, electricity_balance` | 当水费或电费余额更新后低于5元时，自动向该宿舍所有占用床位的学生发送通知：<br>• 若 `NEW.water_balance < 5`，为每个学生生成一条 `notifications` 记录（类型 `low_balance`，消息包含余额信息）<br>• 若 `NEW.electricity_balance < 5`，同理<br>• 多次更新余额且持续低于5元会重复通知，可根据业务需要调整防重逻辑（如一天内仅通知一次） |
| `trg_leave_resolved_at` | `leave_applications` | BEFORE UPDATE OF `status` | 当离校申请的状态变为 `approved` 或 `rejected` 时，自动设置 `resolved_at` 为当前时间。 |
| `trg_late_return_resolved_at` | `late_return_records` | BEFORE UPDATE OF `status` | 当晚归记录状态变为 `approved` 或 `rejected` 时，设置 `resolved_at`。 |
| `trg_room_change_resolved_at` | `room_change_requests` | BEFORE UPDATE OF `status` | 换寝申请状态变为 `approved` 或 `rejected` 时，设置 `resolved_at`。 |
| `trg_off_campus_resolved_at` | `off_campus_living_applications` | BEFORE UPDATE OF `status` | 校外居住申请状态变为 `approved` 或 `rejected` 时，设置 `resolved_at`。 |
| `trg_allocation_resolved_at` | `allocation_requests` | BEFORE UPDATE OF `status` | 分配请求状态变为 `approved` 或 `rejected` 时，设置 `resolved_at`。 |

**注意**：
- 时间戳触发器可通过统一函数减少重复，但为清晰起见单独列出。
- 复合状态逻辑（如换寝审批通过后的床位变更）仍建议在应用层实现，因其涉及跨表事务和复杂校验；触发器在此仅负责简单的时间戳自动记录。
- 低余额通知触发器可根据实际需求增加防骚扰逻辑，例如查询最近一条通知时间间隔。

---

## 8. 实施要点

- **数据库选型**：PostgreSQL 12+，利用其 TOAST 机制自动压缩和外部存储 BYTEA 数据，保持主表行体积合理。
- **索引策略**：
  - 所有外键列建立索引。
  - 工单表（`repair_requests`, `cleaning_requests`）建立 `(status, created_at)` 复合索引，加速待办查询。
  - `notifications` 建立 `(recipient_id, is_read)` 复合索引。
  - `attachments` 建立 `(owner_type, owner_id, category)` 复合索引。
- **备份方案**：由于图片入库，备份文件较大，建议采用 `pg_dump` 自定义格式压缩备份，并结合 WAL 连续归档进行时间点恢复。
- **性能优化**：
  - 查询图片时永远不要 `SELECT *` 包含 `file_data`，除非需要下载原图；列表场景使用 `v_attachment_metadata`。
  - 长列表分页展示。
- **安全设计**：
  - 密码采用 bcrypt 加密存储。
  - 应用层依赖 `role` 字段进行权限控制；数据库层必要时可配置行级安全策略（RLS）。
  - 所有图片读取接口必须校验当前用户是否有权访问（如本人头像、本人所在宿舍的工单照片等）。
- **数据报表**：宿舍人数、空闲床位等直接通过视图实时计算，避免在基本表中维护冗余统计字段。
- **触发器使用**：利用数据库触发器自动记录状态变更时间戳，强制床位数据一致性，以及低余额自动通知，简化应用逻辑并确保规则统一执行。

---

## 9. 业务需求最终对照

| 需求编号 | 需求描述 | 涉及核心表 / 视图 / 触发器 |
|---------|----------|---------------------------|
| 1 | 新生填写生活习惯调查 → 系统自动分配 → 管理员审批 | `lifestyle_surveys`, `allocation_requests`, `beds`, `v_available_beds`, `trg_allocation_resolved_at` |
| 2 | 学生提交离校申请 → 宿管审批 | `leave_applications`, `trg_leave_resolved_at` |
| 3 | 水电费缴费 | `utility_payments`, `rooms` |
| 4 | 宿舍水电费查询 | `rooms` |
| 5 | 维修申请 → 维修人员处理并上传照片 → 宿管审核 | `repair_requests`, `attachments`, `v_pending_repairs`, `trg_repair_status_timestamp` |
| 6 | 换寝申请（可指定或推荐）→ 宿管审批 | `room_change_requests`, `beds`, `v_available_beds`, `trg_room_change_resolved_at` |
| 7 | 晚归记录 → 宿管审批 | `late_return_records`, `trg_late_return_resolved_at` |
| 8 | 自动计算宿舍楼房间数、人数、空闲床位 | `v_dormitory_summary`, `v_available_beds` |
| 9 | 公共区域保洁申请 → 保洁处理并上传照片 → 宿管审核 | `cleaning_requests`, `attachments`, `v_pending_cleanings`, `trg_cleaning_status_timestamp` |
| 10 | 水电费低于5元通知宿舍所有学生 | `rooms`, `v_low_balance_rooms`, `notifications`, `v_student_roommates`, `trg_low_balance_notification` |
| 11 | 学生查询自身所有工单状态 | `v_my_requests` |
| 12 | 查看舍友信息 | `v_student_roommates` |
| 13 | 离校/节假日离校申请（含目的地、紧急联系人等） | `leave_applications` |
| 14 | 校外居住申请（可选保留床位）→ 宿管审批 | `off_campus_living_applications`, `beds`, `trg_off_campus_resolved_at` |
