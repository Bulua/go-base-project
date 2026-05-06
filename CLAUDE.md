# CLAUDE.md

本文件是 GoBaseProject 的项目级开发规则，也是本仓库给自动化助手和协作者的默认中文规范。修改代码、脚本、配置或文档前，请先阅读并遵守这里的约定。

英文版本见 [CLAUDE-EN.md](./CLAUDE-EN.md)。如中英文内容不一致，以本文件为准。

## 1. 项目定位

GoBaseProject 是一个 Go + Vue 3 的后台管理基础项目，目标是提供一套可以通过 Docker Compose 启动、便于二次开发的管理系统基座。

当前项目包含：

1. Go 后端服务
2. Vue 3 / TypeScript 前端应用
3. MySQL 初始化脚本
4. Redis 缓存配置
5. Docker Compose 本地和远程依赖启动方式

`development-plan/` 只作为本地计划和设计参考，已经被 `.gitignore` 忽略。正式源码、运行脚本、Docker 配置和可提交文档都不能依赖它。

## 2. 常用命令

从项目根目录复制环境变量文件：

```bash
cp .env.example .env
```

Windows PowerShell：

```powershell
Copy-Item .env.example .env
```

使用远程 MySQL / Redis：

```bash
docker compose up -d --build
```

使用本地 MySQL / Redis：

```bash
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

后端测试：

```bash
cd server
go test ./...
```

前端构建检查：

```bash
cd web
npm install
npm run build
```

Compose 配置检查：

```bash
docker compose config
```

如果当前环境缺少 Docker、Go 或 Node.js，无法执行相关验证时，需要在最终说明中明确写出未验证的项目。

## 3. 根目录约定

```text
.
|-- AGENTS.md                 # Agent 入口说明
|-- CLAUDE.md                 # 中文主规范
|-- CLAUDE-EN.md              # 英文规范
|-- README.md                 # 面向使用者的快速说明
|-- docker-compose.yaml       # 默认 Compose 配置
|-- docker-compose.local.yaml # 本地 MySQL / Redis 覆盖配置
|-- deploy                    # 部署、初始化、Nginx 等资源
|-- docs                      # 可提交的详细文档
|-- server                    # Go 后端
|-- web                       # Vue 前端
`-- development-plan          # ignored，仅本地计划和设计参考
```

可提交的运行文档放在 `README.md` 和 `docs/` 下。部署脚本、初始化脚本、Nginx 配置等放在 `deploy/` 下。

## 4. 敏感信息规则

1. `.env` 只用于本地环境，不提交。
2. `.env.example` 只能放占位值，例如 `change-me`。
3. 不要在代码、测试、日志、Markdown、SQL、Dockerfile 或配置模板中写真实密码、Token、密钥。
4. 默认演示账号可以写入文档，但生产账号和真实凭据不能写入文档。
5. 修改前后可用 `rg` 检查是否误写敏感信息：

```bash
rg "<REAL_SECRET_TO_CHECK>" -g "!.env" -g "!development-plan/**" -g "!web/node_modules/**" -g "!web/dist/**" .
```

命令没有输出才表示该关键词未泄露。

## 5. 后端规范

后端目录是 `server/`。

当前后端使用 Go 1.22、标准库 `net/http`、MySQL、JWT、bcrypt、YAML 配置加载。新增依赖前应确认确实能降低复杂度，并保持现有代码风格。

推荐结构：

```text
server
|-- cmd
|   `-- gobaseproject          # 主程序入口
|-- configs                    # 配置模板
|-- internal
|   |-- handler                # HTTP 处理层
|   |-- service                # 业务逻辑层
|   |-- repository             # 数据访问层
|   |-- model                  # 领域类型、DTO、错误
|   |-- middleware             # 权限、审计、跨域等中间件
|   `-- infra                  # 配置、数据库、缓存等基础设施
|-- pkg                        # 可复用通用包
`-- tests                      # 后端测试
```

分层职责：

```text
handler: 解析请求、校验参数、调用 service、封装响应
service: 业务规则、权限判断、事务编排、跨 repository 协调
repository: 数据库访问、查询条件构建、持久化细节
model: 领域模型、DTO、枚举、错误定义
middleware: 请求链路中的鉴权、审计、上下文处理
infra: 配置、数据库、Redis、日志等基础设施
pkg: 与具体业务弱相关的可复用能力
```

后端规则：

1. 新业务优先按 `handler -> service -> repository -> model` 落位。
2. 跨模块调用优先通过 service，不要直接跨模块访问 repository。
3. 所有接口响应应保持统一格式。
4. 列表接口默认分页，并限制 `page_size` 最大值。
5. 密码只保存哈希值，推荐 bcrypt 或 argon2id。
6. JWT 黑名单只保存 Token 摘要，不保存原始 Token。
7. 登录失败、退出登录、重置密码、权限变更和写操作应记录审计日志。
8. 审计日志必须对密码、Token、密钥等字段脱敏。
9. 提交前运行 `gofmt`，涉及后端逻辑时运行 `go test ./...`。

统一响应格式示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {},
  "trace_id": "request-trace-id"
}
```

## 6. 前端规范

前端目录是 `web/`。

当前前端使用 Vue 3、TypeScript、Vite、Element Plus、Pinia、Vue Router、Axios。新增页面和组件时应贴近现有后台管理风格：信息密度适中、操作明确、不要做成营销页。

推荐结构：

```text
web/src
|-- api
|   |-- request                # Axios 实例、拦截器、错误处理
|   |-- auth
|   |-- user
|   `-- ...
|-- assets
|   `-- styles                 # 全局样式、主题变量
|-- components
|   |-- common                 # 通用组件
|   |-- business               # 业务复合组件
|   `-- layout                 # 布局组件
|-- composables                # useXxx 组合函数
|-- directives                 # v-auth 等指令
|-- layouts                    # 页面框架
|-- permissions                # 权限判断
|-- router                     # 静态和动态路由
|-- store
|   `-- modules                # Pinia 模块
|-- types                      # DTO、枚举、类型声明
|-- utils                      # 纯工具函数
`-- views                      # 页面级组件
```

前端规则：

1. 页面组件入口使用 `index.vue`。
2. 业务组件使用 PascalCase，例如 `RoleSelect.vue`。
3. 组合函数使用 `useXxx.ts`。
4. API 文件按资源命名，例如 `api/user/index.ts`。
5. 所有接口请求必须通过统一 Axios 实例。
6. 前端权限只控制展示，不能替代后端鉴权。
7. 组件优先复用 Element Plus 和项目内已有组件。
8. 涉及类型变更时，同步更新 `types/`、`api/`、`store/` 和页面使用处。
9. 涉及前端逻辑时至少运行 `npm run build`，有测试时补充 `npm run test`。

## 7. 数据库和初始化

数据库初始化脚本放在：

```text
deploy/mysql/init
```

当前可提交初始化文件是：

```text
deploy/mysql/init/init.sql
```

数据库名：

```text
go_base_project
```

规则：

1. 不要依赖 `development-plan/sql` 作为运行脚本。
2. 初始化脚本应可重复审查，表结构和种子数据变更需要说明影响。
3. 生产数据库变更后续应引入迁移机制，并在 README 或 `docs/` 中说明执行方式。
4. 修改 SQL 后，尽量用本地容器或测试数据库验证启动和初始化。

## 8. Docker 规范

`docker-compose.yaml` 是默认启动入口，用于连接 `.env` 中配置的外部 MySQL / Redis。

`docker-compose.local.yaml` 是本地开发覆盖文件，用于启动本地 MySQL / Redis。

规则：

1. 修改 Compose 后运行 `docker compose config`。
2. 修改容器端口、服务名或环境变量时，同步更新 `.env.example`、README 和 `docs/RUNNING.md`。
3. 不要把真实密码写进 Compose 文件。
4. 不要提交本地生成的 volume、构建产物或临时日志。

## 9. 文档规范

1. `README.md` 面向使用者，写快速启动、配置、入口和项目概览。
2. `docs/RUNNING.md` 写详细运行说明。
3. `CLAUDE.md` 写中文长期开发规范。
4. `CLAUDE-EN.md` 是英文翻译，供英文 Agent 或协作者参考。
5. `AGENTS.md` 保持轻量，只写 Agent 入口和快速上下文。
6. 文档中的树形结构尽量使用 ASCII 字符，减少 Windows 终端编码问题。

## 10. 文件和提交流程

1. 新文件默认 UTF-8。
2. 脚本、配置、Markdown 尽量使用 ASCII 标点和路径符号。
3. 不要提交构建产物：`web/dist/`、`web/node_modules/`、`server/bin/`。
4. 不要删除 `.gitkeep`，除非目录中已有真实文件。
5. 修改目录结构时，同步更新 README、docs 和 Docker 配置。
6. 工作区可能已有他人改动，修改前先确认 `git status`，不要回退自己没有创建的改动。
7. 完成后在回复中说明改了哪些文件、执行了哪些验证、哪些验证没有执行。

## 11. 开发优先级建议

优先级建议：

1. 配置加载、日志、统一响应
2. 数据库连接和初始化 / 迁移策略
3. Redis 连接
4. 认证鉴权
5. 用户、角色、菜单、权限
6. 审计日志
7. 文件、字典、系统参数等基础模块
8. Vue 登录页、基础布局、动态路由和按钮权限
9. 各管理页面 CRUD 和体验打磨

每完成一块功能，都应补齐对应的最小验证方式。
