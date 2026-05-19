# 宿舍管理系统数据库脚本说明

本目录用于存放宿舍管理系统的 PostgreSQL 数据库初始化脚本。

## 文件说明

| 文件 | 说明 |
|------|------|
| `001_create_student_dormitory_schema.sql` | 创建宿舍管理系统的表、索引、视图、函数和触发器 |
| `002_seed_student_dormitory_test_data.sql` | 插入前后端联调用测试数据，不包含工单和图片 |
| `003_truncate_student_dormitory_data.sql` | 清空所有业务数据并重置自增序列，保留表结构、视图和触发器 |
| `004_seed_roommate_recommendation_test_data.sql` | 插入自动舍友/床位推荐专项测试数据，构造可预测的兼容与不兼容候选宿舍 |

## 数据库前提

脚本默认数据库已提前创建：

- 数据库名：`student_dormitory`
- 数据库用户：`admin`
- 用户密码：`passwd`
- 数据库类型：PostgreSQL 12+

如果尚未创建数据库，可先使用具备建库权限的账号执行：

```sql
CREATE DATABASE student_dormitory OWNER admin;
```

## 执行方式

在项目根目录下执行：

```bash
PGPASSWORD=passwd psql -U admin -d student_dormitory -f sql/001_create_student_dormitory_schema.sql
```


插入测试数据：

```bash
PGPASSWORD=passwd psql -U admin -d student_dormitory -f sql/002_seed_student_dormitory_test_data.sql
```

插入自动舍友推荐专项测试数据：

```bash
PGPASSWORD=passwd psql -U admin -d student_dormitory -f sql/004_seed_roommate_recommendation_test_data.sql
```

清空所有数据：

```bash
PGPASSWORD=passwd psql -U admin -d student_dormitory -f sql/003_truncate_student_dormitory_data.sql
```

测试账号默认密码均为 `123456`。常用账号：

| 用户名 | 角色 | 说明 |
|--------|------|------|
| `admin001` | system_admin | 系统管理员 |
| `manager001` | dormitory_manager | 宿舍管理员 |
| `repair001` | repair_staff | 维修人员 |
| `cleaner001` | cleaning_staff | 保洁人员 |
| `student001` - `student008` | student | 学生 |
| `student101` - `student106` | student | 自动舍友推荐专项测试学生 |

## 自动舍友推荐专项数据

`004_seed_roommate_recommendation_test_data.sql` 会创建独立的 `自动推荐测试楼`，包含 701、702、703 三个测试房间：

- `student101`：待分配新生，已有生活习惯问卷，无床位、无 pending 分配申请。
- `student102`、`student103`：入住 701 A/B，作息、吸烟、打鼾和学习习惯与 `student101` 高度匹配。
- `student104`、`student105`：入住 702 A/B，与 `student101` 明显不匹配。
- `student106`：入住 703 A，无生活习惯问卷，用于验证无画像舍友的低权重评分。

执行脚本后，以 `student101` 登录并调用 `POST /api/allocations`，预期生成的 pending 分配请求推荐 `自动推荐测试楼 701 C床`。

该脚本会清理自身测试用户和测试房间相关的分配请求、问卷和床位状态，因此可以重复执行以恢复同一测试场景。

## 脚本内容

`001_create_student_dormitory_schema.sql` 会创建以下对象：

- 业务表：`users`、`buildings`、`rooms`、`beds`、`lifestyle_surveys`、`allocation_requests`、`leave_applications`、`utility_payments`、`repair_requests`、`cleaning_requests`、`late_return_records`、`room_change_requests`、`notifications`、`attachments`、`off_campus_living_applications`
- 查询视图：`v_dormitory_summary`、`v_available_beds`、`v_student_roommates`、`v_low_balance_rooms`、`v_pending_repairs`、`v_pending_cleanings`、`v_my_requests`、`v_attachment_metadata`
- 触发器：自动状态时间戳、床位一致性校验、换寝目标校验、低余额通知、水电缴费后余额更新
- 索引：外键列索引、工单状态索引、通知查询索引、附件关联索引等

## 重复执行说明

脚本按便于调试的方式编写：

- 表和索引使用 `IF NOT EXISTS`
- 视图和函数使用 `CREATE OR REPLACE`
- 触发器会先 `DROP TRIGGER IF EXISTS`，再重新创建

因此脚本可以重复执行，用于补建对象或更新视图、函数、触发器。

注意：如果已经存在同名表且字段结构与脚本不一致，`CREATE TABLE IF NOT EXISTS` 不会修改已有表结构。此时需要手动编写迁移脚本处理结构变更。

## 附件存储说明

系统图片统一存储在 `attachments.file_data` 字段中，类型为 `BYTEA`。列表查询应优先使用 `v_attachment_metadata`，避免不必要地读取二进制图片数据。

支持的附件关联方式：

| 业务对象 | `owner_type` | `category` |
|----------|--------------|------------|
| 用户头像 | `user_avatar` | `avatar` |
| 维修后照片 | `repair` | `after` |
| 保洁前照片 | `cleaning` | `before` |
| 保洁后照片 | `cleaning` | `after` |

## 低余额通知说明

当 `rooms.water_balance` 或 `rooms.electricity_balance` 更新后低于 5 元时，触发器会向该宿舍所有已入住学生写入 `notifications` 记录。

当前实现会在余额低于 5 元且余额值发生变化时发送通知。若业务需要限制通知频率，例如一天只提醒一次，可在触发器函数 `fn_low_balance_notification` 中增加时间窗口判断。

## 水电缴费说明

向 `utility_payments` 插入缴费记录后，触发器会自动更新对应宿舍余额：

- `payment_type = 'water'`：增加水费余额
- `payment_type = 'electricity'`：增加电费余额
- `payment_type = 'both'`：金额平均分配到水费和电费余额

如果实际业务需要前端分别传入水费金额和电费金额，建议后续将缴费表拆分为两个金额字段，或新增明细表。
