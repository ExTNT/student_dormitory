# 前端设计方案说明

## 设计目标

前端定位为后台管理系统，不做营销首页。用户登录后直接进入对应角色工作台，通过左侧菜单访问业务功能。页面风格以 Element Plus 的表格、表单、弹窗和 Tag 为主，强调清晰、稳定、易维护。

系统严格按 `backend/docs/API.md` 对接后端接口，不虚构接口，不硬编码当前用户 ID。当前用户信息统一通过 `GET /students/me` 获取。

## 技术选型

| 技术 | 用途 |
| --- | --- |
| Vue 3 | 页面和组件开发 |
| Vite | 开发服务器与构建 |
| TypeScript | 类型约束 |
| Vue Router | 路由与权限守卫 |
| Pinia | 登录态和通知状态 |
| Axios | REST API 请求 |
| Element Plus | 后台管理 UI |
| ECharts | 楼栋统计图表 |

## 模块划分

```text
src/api
```

按业务域封装接口：

- `http.ts`：Axios 实例、baseURL、认证头、401 refresh、统一错误提示。
- `auth.ts`：登录、登出、当前用户。
- `student.ts`：学生端接口和通用楼栋/房间查询。
- `repair.ts`：维修工单。
- `cleaning.ts`：保洁工单。
- `manager.ts`：宿管审批。
- `admin.ts`：管理员用户和分配审批。
- `attachment.ts`：附件上传、列表、下载。
- `dashboard.ts`：统计和低余额。

```text
src/stores
```

- `useAuthStore`：保存 token、当前用户、登录、登出、获取当前用户。
- `useNotificationStore`：通知列表和标记已读。

```text
src/router
```

集中定义所有路由，并在 `meta.roles` 中标记允许访问的角色。全局守卫负责登录检查和权限检查。

```text
src/components
```

封装通用 UI：

- `StatusTag`：统一申请和工单状态显示。
- `ImageUploader`：附件上传。
- `AttachmentImage`：blob 图片加载与预览。
- `AttachmentList`：附件元数据列表与图片预览。
- `ReviewButtons`：审批确认和带意见审核。

```text
src/views
```

按角色和业务组织页面：

- `student/`
- `repair/`
- `cleaning/`
- `manager/`
- `admin/`
- `dashboard/`

## 认证与 token 刷新

认证状态保存在 `localStorage`：

- `access_token`
- `refresh_token`

请求流程：

1. 请求拦截器读取 `access_token`。
2. 存在 token 时写入 `Authorization: Bearer <access_token>`。
3. 响应为 `401` 时，如果存在 `refresh_token`，调用 `/auth/refresh`。
4. refresh 成功后替换本地两个 token。
5. 使用新 access token 重试原请求。
6. refresh 失败则清空 token 并跳转 `/login`。

并发刷新处理：

`http.ts` 中使用共享 `refreshPromise`。多个请求同时收到 `401` 时只发起一次刷新，后续请求等待同一个 Promise，避免旧 refresh token 被重复使用。

## 权限设计

路由通过 `meta.roles` 声明权限，例如：

```ts
{
  path: 'manager/leaves',
  component: () => import('@/views/manager/LeavesReviewView.vue'),
  meta: { roles: ['dormitory_manager'] },
}
```

守卫规则：

- `meta.public` 页面直接放行，例如 `/login`、`/403`。
- 没有 token 的用户访问业务页面，跳转 `/login`。
- 有 token 但没有当前用户信息时，调用 `fetchMe()`。
- 当前用户角色不在 `meta.roles` 中时，跳转 `/403`。

左侧菜单也按 `auth.user.role` 动态生成，避免无关入口干扰用户。

## 页面交互设计

### 工作台

通用 `DashboardView` 根据角色显示当前用户姓名、角色和快捷入口。

### 状态展示

申请和工单状态通过 `StatusTag` 统一映射：

| 状态 | Tag 类型 |
| --- | --- |
| pending | warning |
| approved | success |
| rejected | danger |
| accepted | primary |
| repaired | primary |
| cleaned | primary |
| completed | success |
| paid | success |

### 日期时间

展示使用浏览器本地格式：

```ts
new Date(value).toLocaleString()
```

提交给后端的日期时间使用 ISO 8601：

```ts
new Date(value).toISOString()
```

晚归日期按接口要求提交 `YYYY-MM-DD`。

### 审批和危险操作

审批、接单、完成等会改变业务状态的操作均使用 Element Plus 弹窗确认。

成功后会重新请求对应列表，确保页面状态与后端一致。

## 角色功能设计

### 学生

学生端围绕个人申请和宿舍生活服务：

- 生活习惯调查作为新生分配申请前置说明。
- 分配申请只调用 `POST /allocations`，推荐逻辑由后端处理。
- 我的申请总览统一展示后端返回的混合申请列表。
- 换寝支持系统推荐或指定空床位；指定时同时提交 `target_room_id` 和 `target_bed_id`。
- 保洁申请创建成功后允许上传 `before` 图片。
- 缴费先查询房间余额，再提交缴费，成功后刷新余额。

### 维修人员

维修人员只处理 `/repairs/pending` 返回的工单：

- `pending` 状态显示接单按钮。
- `accepted` 状态允许上传 after 图片并完成维修。
- 完成前前端先检查 after 附件是否存在。

### 保洁人员

保洁人员只处理 `/cleanings/pending` 返回的工单：

- `pending` 状态显示接单按钮。
- `accepted` 状态允许上传 after 图片并完成清洁。
- 完成前前端先检查 after 附件是否存在。

### 宿舍管理员

宿管处理各类审批和审核：

- 离校审批。
- 晚归审批。
- 换寝审批。
- 校外居住审批。
- 维修审核，仅展示 `status=repaired` 的记录。
- 保洁审核，仅展示 `status=cleaned` 的记录。
- 楼栋统计和低余额房间。

### 系统管理员

管理员负责系统级操作：

- 创建用户。
- 新生分配审批。
- 楼栋统计和低余额房间。

## 附件设计

附件组件拆成三个层级：

### ImageUploader

负责上传：

- `owner_type`
- `owner_id`
- `category`
- `sort_order`
- `file`

上传成功后 emit `success`，由父页面刷新附件列表。

### AttachmentImage

通过 `GET /attachments/{id}` 以 blob 加载图片，并生成本地 object URL。组件销毁时释放 URL，避免内存泄漏。

### AttachmentList

通过 `GET /attachments?owner_type=...&owner_id=...&category=...` 获取元数据，列表中复用 `AttachmentImage` 展示和预览。

## 错误处理

后端通用错误格式：

```json
{
  "error": "bad_request",
  "message": "错误详情"
}
```

Axios 响应拦截器优先展示 `message`，其次展示 `error`，最后展示 Axios 错误信息。

`401` 单独处理，不直接弹普通错误，而是尝试刷新 token。

## 构建与部署

开发：

```bash
npm run dev
```

构建：

```bash
npm run build
```

构建产物输出到：

```text
dist/
```

`dist/`、`node_modules/` 和 `*.tsbuildinfo` 已在 `.gitignore` 中忽略。

部署时需要确保静态资源服务器支持 Vue Router history 模式，将未知路径回退到 `index.html`。

## 后续可优化点

- 将 API 基础地址改为 Vite 环境变量，例如 `VITE_API_BASE_URL`。
- 对 ECharts 单独拆 chunk，减少首包提示。
- 补充表单必填规则和更细粒度的校验提示。
- 引入端到端测试覆盖登录、角色跳转和核心审批流程。
