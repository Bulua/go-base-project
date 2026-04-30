# CLAUDE.md

本文件是 GoBaseProject 的项目级开发规则。任何自动化助手、协作者或后续维护者在修改本项目时，都应优先遵守这里的约定。

## 1. 项目定位

GoBaseProject 是一个 Go + Vue3 的后台管理基座项目。

目标是提供一套可通过 Docker Compose 启动的标准工程，包括：

1. Go 后端服务
2. Vue3 前端应用
3. Redis 缓存配置
4. MySQL 初始化脚本
5. 可复制的运行规范

`development-plan/` 是本地开发计划和设计参考，已经被 `.gitignore` 忽略。不要让正式源码、运行脚本、Docker 配置或文档依赖 `development-plan/`。

## 2. 启动规范

仓库 clone 后应从项目根目录启动：

```bash
cp .env.example .env
docker compose up -d --build
```

Windows PowerShell：

```powershell
Copy-Item .env.example .env
docker compose up -d --build
```

默认远程服务：

```text
MySQL: 120.53.251.75:13307
Redis: 120.53.251.75:16379
```

真实密码只允许写入本地 `.env`，不要写入 `.env.example`、README、代码、Dockerfile、SQL 或其他可提交文件。

本地 MySQL：

```bash
docker compose --profile local-db up -d --build
```

本地 Redis：

```bash
docker compose --profile local-cache up -d --build
```

## 3. 根目录约定

```text
.
|-- docker-compose.yaml
|-- deploy
|-- docs
|-- server
|-- web
`-- development-plan   # ignored, local planning only
```

可提交运行文档放在 `README.md` 和 `docs/` 下。部署、初始化、Nginx、脚本等放在 `deploy/` 下。

## 4. 敏感信息规则

1. `.env` 不提交。
2. `.env.example` 只能放占位密码，例如 `change-me`。
3. 不要在代码、测试、日志、Markdown 或 SQL 中写真实密码、Token、密钥。
4. 修改前后可用下面的命令检查泄露：

```bash
rg "<REAL_SECRET_TO_CHECK>" -g "!.env" -g "!development-plan/**" -g "!web/node_modules/**" -g "!web/dist/**" .
```

命令没有输出才算通过。

## 5. 后端规范

后端目录：`server/`。

技术栈：Go，当前基础服务使用标准库 HTTP，后续可逐步接入 Gin、GORM、Viper、Zap、JWT、Redis 客户端等。

后端采用“类型优先、模块次级”的结构：

```text
server/internal
|-- app
|-- infra
|-- handler
|   |-- auth
|   |-- user
|   |-- role
|   `-- ...
|-- service
|   |-- auth
|   |-- user
|   |-- role
|   `-- ...
|-- repository
|   |-- auth
|   |-- user
|   |-- role
|   `-- ...
`-- model
    |-- auth
    |-- user
    |-- role
    `-- ...
```

分层职责：

```text
handler: HTTP 入参解析、参数校验、响应封装
service: 业务规则、事务编排、权限判断
repository: 数据库访问和查询条件构建
model: GORM 模型、DTO、领域常量
infra: 数据库、Redis、日志、配置等基础设施
app: 启动、路由、中间件、生命周期
pkg: 可复用通用包
```

后端规则：

1. 新业务优先按 `handler -> service -> repository -> model` 落位。
2. 跨模块调用优先通过 service，不要直接跨模块访问 repository。
3. 统一响应格式应保持：

```json
{
  "code": 0,
  "message": "ok",
  "data": {},
  "trace_id": "request-trace-id"
}
```

4. 所有列表接口默认分页，并限制 `page_size` 最大值。
5. 密码只保存哈希值，推荐 bcrypt 或 argon2id。
6. JWT 黑名单只保存 Token 摘要，不保存原始 Token。
7. 登录失败、退出登录、重置密码和权限变更必须记录审计日志。
8. 操作审计必须对密码、Token、密钥等字段脱敏。

后端验证：

```bash
cd server
go test ./...
```

提交前应运行 `gofmt`。

## 6. 前端规范

前端目录：`web/`。

技术栈：Vue3 + TypeScript + Vite + Element Plus + Pinia + Vue Router + Axios。

前端样式参考来自 `development-plan/web-design/`，但正式实现必须落在 `web/` 下，不要在运行代码中引用 `development-plan/`。

前端同样采用“类型优先、模块次级”的结构：

```text
web/src
|-- api
|   |-- request
|   |-- auth
|   |-- user
|   `-- ...
|-- assets
|-- components
|   |-- common
|   |-- business
|   `-- layout
|-- composables
|-- constants
|-- directives
|-- layouts
|-- permissions
|-- router
|-- store
|   `-- modules
|-- types
|-- utils
`-- views
```

前端放置规则：

```text
api/request: Axios 实例、拦截器、错误处理
api/<module>: 业务接口请求
views: 页面级组件
store/modules/<module>: Pinia 模块状态
types/<module>: DTO、枚举、类型声明
components/common: 通用组件
components/business: 业务复合组件
components/layout: 布局组件
composables: useXxx 组合函数
directives: v-auth 等指令
permissions: 权限判断、动态路由转换
utils: 纯工具函数
assets/styles: 全局样式、主题变量
```

前端规则：

1. 页面组件入口使用 `index.vue`。
2. 业务组件使用 PascalCase，例如 `RoleSelect.vue`。
3. 组合函数使用 `useXxx.ts`。
4. API 文件按资源命名，例如 `api/user/index.ts`。
5. 前端权限只控制展示，不能替代后端鉴权。
6. 所有接口请求必须通过统一 Axios 实例。
7. UI 实现要贴近 `development-plan/web-design/` 的后台管理风格：信息密度适中、操作清晰、不要做营销页。

前端验证：

```bash
cd web
npm install
npm run build
```

## 7. 数据库规范

可提交的数据库初始化脚本放在：

```text
deploy/mysql/init
```

不要依赖 `development-plan/sql` 作为运行脚本。

当前数据库名：

```text
go_base_project
```

初始化脚本应保持可重复审查、顺序明确：

```text
001_schema.sql
002_seed.sql
```

生产数据库变更后续应引入迁移机制，迁移文件放在 `server/migrations/` 或统一迁移目录中，并在 README 中说明执行方式。

## 8. Docker 规范

`docker-compose.yaml` 是 clone 后的标准启动入口。

默认模式使用远程 MySQL 和远程 Redis，只启动：

```text
server
web
```

可选 profile：

```text
local-db: 启动本地 MySQL
local-cache: 启动本地 Redis
```

Compose 修改后应尽量验证：

```bash
docker compose config
```

如果当前机器没有 Docker，需要在最终说明中明确未执行 Docker 校验。

## 9. 文档规范

1. 根 README 只写 clone、配置、启动、常用入口。
2. 详细运行说明写在 `docs/RUNNING.md`。
3. 项目长期规范写在 `CLAUDE.md`。
4. `AGENTS.md` 只引用 `CLAUDE.md`，不要复制规则，避免两份文件分叉。
5. 文档尽量使用 ASCII 树形字符，避免在 Windows 终端出现编码乱码。

## 10. 编码和文件规则

1. 新文件默认 UTF-8。
2. 脚本、配置、Markdown 尽量使用 ASCII 标点和路径符号。
3. 不要提交构建产物：`web/dist/`、`web/node_modules/`、`server/bin/`。
4. 不要删除 `.gitkeep`，除非目录中已有真实文件。
5. 修改目录结构时，同步更新 README、docs 和 Docker 配置。

## 11. 开发顺序建议

优先级建议：

1. 后端配置加载、日志、统一响应
2. 数据库连接和迁移/初始化策略
3. Redis 连接
4. 认证鉴权
5. 用户、角色、菜单、权限
6. 审计日志
7. Vue3 登录页和基础布局
8. 动态路由、菜单、按钮权限
9. 各管理页面 CRUD

每次完成一块功能，都应补齐最小验证方式。
