# 宿舍管理系统前端

基于 Vue 3、Vite、TypeScript、Vue Router、Pinia、Axios、Element Plus 和 ECharts 实现的宿舍管理系统前端，对接 Go 后端 REST API。

后端接口文档：`../backend/docs/API.md`

默认 API 基础地址：`http://localhost:8080/api`

## 功能范围

- 登录、登出、access token 自动注入、401 自动刷新 token。
- 按角色显示菜单和限制路由访问。
- 未登录访问业务页面自动跳转 `/login`。
- 无权限访问跳转 `/403`。
- 学生端：个人信息、生活习惯调查、新生分配、申请总览、舍友、离校、晚归、换寝、校外居住、维修、保洁、水电缴费、通知。
- 维修人员：待处理维修工单、接单、上传维修后照片、完成维修。
- 保洁人员：待处理保洁工单、接单、上传保洁后照片、完成清洁。
- 宿舍管理员：离校、晚归、换寝、校外居住审批，维修/保洁审核，楼栋统计，低余额房间。
- 系统管理员：创建用户，新生分配审批，楼栋统计，低余额房间。
- 附件：图片上传、图片预览、附件列表。

## 环境要求

- Node.js 18+。当前开发环境使用 Node.js 24。
- npm。
- 后端服务启动在 `http://localhost:8080`。

## 启动

```bash
cd frontend
npm install
npm run dev
```

开发服务器默认监听：

```text
http://localhost:5173/
```

生产构建：

```bash
npm run build
```

本地预览构建产物：

```bash
npm run preview
```

## 测试账号

数据库测试数据脚本提供以下账号，密码均为 `123456`：

```text
admin001
manager001
repair001
cleaner001
student001
student002
student003
student004
student005
student006
student007
student008
```

## 目录结构

```text
frontend/
  src/
    api/          REST API 封装和 Axios 实例
    components/   通用组件，包含附件和状态组件
    router/       路由定义与权限守卫
    stores/       Pinia 状态
    types/        TypeScript 类型定义
    utils/        格式化和映射工具
    views/        页面
```

## API 配置

API 基础地址定义在：

```text
src/api/http.ts
```

当前值：

```ts
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';
```

开发环境默认通过 Vite 代理转发到 `http://localhost:8080`。如果后端地址变化，可以设置 `VITE_API_BASE_URL` 覆盖默认值。

## 认证说明

登录调用：

```text
POST /auth/login
```

登录成功后保存：

- `access_token`
- `refresh_token`

所有业务请求都会自动注入：

```text
Authorization: Bearer <access_token>
```

当后端返回 `401` 时，Axios 响应拦截器会调用：

```text
POST /auth/refresh
```

刷新成功后会替换本地 `access_token` 和 `refresh_token`，并重试原请求。多个请求同时触发刷新时，共享同一个 refresh Promise，避免重复刷新导致旧 refresh token 失效。

刷新失败时会清空登录态并跳转登录页。

## 角色入口

登录成功后会调用：

```text
GET /students/me
```

根据当前用户角色跳转：

| 角色 | 路由 |
| --- | --- |
| student | `/student/dashboard` |
| repair_staff | `/repair/dashboard` |
| cleaning_staff | `/cleaning/dashboard` |
| dormitory_manager | `/manager/dashboard` |
| system_admin | `/admin/dashboard` |

## 附件限制

`ImageUploader` 使用 `multipart/form-data` 调用：

```text
POST /attachments
```

前端限制：

- 仅允许 `image/jpeg`、`image/png`
- 最大 5MB

## 常见问题

### 登录时报网络错误

确认 Go 后端是否已经启动，并且接口基础地址为：

```text
http://localhost:8080/api
```

### 访问业务页面跳到登录页

本地没有有效 `access_token`，或刷新 token 已失效。重新登录即可。

### 访问页面显示 403

当前登录用户角色不在该路由允许的角色列表中。

### 构建时出现 chunk size warning

ECharts 和 Element Plus 体积较大，Vite 可能提示 chunk 超过 500KB。这是构建体积提示，不影响运行。
