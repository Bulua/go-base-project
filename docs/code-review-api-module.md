# API 管理模块代码审查报告

审查范围：`server/internal/{handler,service,repository,model}/api/`、
`web/src/{api,types,views/system}/api/`、`deploy/mysql/init/007_seed_api_mgmt_apis.sql`

参考基准：现有 User / Role / Menu 模块的实现风格。

---

## 1. Handler 层

### P0 — pathUint64 返回 0 时未拒绝请求

`server/internal/handler/api/handler.go` 中路径参数解析失败时返回 `id=0`，
下游直接使用该值，最终触发 `record not found` 而非 `invalid param`，
错误语义错误，且可能在某些边界情况下绕过存在性检查。

**修复方式**

```go
id, err := pathUint64(r, "id")
if err != nil || id == 0 {
    apperror.WriteDefinition(w, r, apperror.InvalidParams)
    return
}
```

### P1 — Request 类型未定义到 model 层

`createAPI` / `updateAPI` 使用匿名 struct 解析请求体，校验逻辑因此流入 Service 层，
Handler 与 Service 的边界模糊。对比 User 模块有 `model.CreateUserRequest` 等正式类型。

**修复方式**：在 `server/internal/model/api/model.go` 中补充：

```go
type CreateAPIRequest struct {
    APIPath   string `json:"api_path"`
    APIMethod string `json:"api_method"`
    APIGroup  string `json:"api_group"`
    APIDesc   string `json:"api_desc"`
}

type UpdateAPIRequest struct {
    APIPath   string `json:"api_path"`
    APIMethod string `json:"api_method"`
    APIGroup  string `json:"api_group"`
    APIDesc   string `json:"api_desc"`
    APIStatus int    `json:"api_status"`
}
```

Handler 只负责解析到该类型并做基础格式校验，业务校验留在 Service。

### P2 — intParam 与其他模块工具函数重复

Handler 内定义了局部的 `intParam()`，而其他模块已有类似实现（`parseIntDefault` 等）。
应提取到公共 handler util 包，避免多处散落。

---

## 2. Service 层

### P0 — 缺少审计日志

`server/internal/service/api/service.go` 的所有写操作（Create / Update / Delete）
均未调用审计日志接口。对比 User 模块，每个写操作都有 `s.logSuccess()` / `s.logFailure()`。

API 权限管理是高敏感操作，无日志则无法追溯变更者，同时违反 CLAUDE.md 规则第 7 条。

**修复方式**：为 `APIService` 注入 `AuditRepository`，在每个写操作的成功/失败路径上
记录操作者、操作类型、目标资源 ID。

### P0 — CreateAPI 存在竞态条件

"检查存在 → 插入" 两步之间没有事务或数据库唯一约束的正确翻译：

```
线程 A: ExistsAPIPathMethod -> false
线程 B: ExistsAPIPathMethod -> false
线程 A: CreateAPI -> 成功
线程 B: CreateAPI -> UNIQUE 冲突 -> Internal Server Error  (应为 AlreadyExists)
```

**修复方式**：捕获数据库的 duplicate key 错误并转换为业务错误码 `AlreadyExists`，
或在数据库层加唯一索引后在 Repository 层统一翻译错误类型。

### P0 — DeleteAPI 不检查权限策略引用

删除 API 时未检查 `gbp_permission_policies` 是否仍有该资源的记录，
删除后会遗留孤儿权限配置，导致权限策略引用不存在的资源。

**修复方式**：删除前查询关联的策略条目，有则拒绝删除并返回提示，
或在事务中级联删除关联策略（需与产品确认语义）。

### P1 — CreateSkipRule 多了一次不必要的数据库查询

插入后用 `Keyword: p.APIPath` 再查一次来"找回"刚插入的记录，逻辑复杂且存在竞态：

```go
// 当前：3 次 DB 操作
id, _ := s.repo.CreateSkipRule(ctx, p)
list, _ := s.repo.ListSkipRules(ctx, Query{Keyword: p.APIPath})  // 不稳定
return list.Items[0], nil
```

**修复方式**：用 `lastInsertId` + `GetSkipRuleByID` 代替：

```go
id, _ := s.repo.CreateSkipRule(ctx, p)
return s.repo.GetSkipRuleByID(ctx, id)
```

### P1 — DeleteSkipRule 未验证记录存在性

`DeleteSkipRule()` 直接执行 DELETE，删除 0 行时静默成功，
与 `DeleteAPI()` 先 Get 再 Delete 的行为不一致。

**修复方式**：删除前调用 `GetSkipRuleByID`，不存在则返回 `NotFound`。

### P2 — 分页限制使用魔法数字

`PageSize > 200` 重置为 `20`，这两个数字应定义为包级常量，
并复用项目已有的全局分页限制（如果存在）。

---

## 3. Repository 层

### P1 — LIKE 搜索未转义特殊字符

`buildAPIWhere()` 和 `buildSkipWhere()` 未对 Keyword 中的 `%`、`_` 做转义：

```go
// 用户输入 "50%" 会匹配所有含 "50" 且后跟任意字符的记录
args = append(args, "%"+q.Keyword+"%")
```

**修复方式**：

```go
func escapeLike(s string) string {
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `%`, `\%`)
    s = strings.ReplaceAll(s, `_`, `\_`)
    return s
}
// 使用时
args = append(args, "%"+escapeLike(q.Keyword)+"%")
// SQL 中加 ESCAPE '\'
```

### P2 — 空字符串与 NULL 处理不统一

`nullString()` 将空字符串转为 SQL NULL，但 ListAPIs 中用 `COALESCE(a.api_group, '')` 
把 NULL 再转回空字符串，一进一出抵消，且使空字符串语义丢失。
建议统一：数据库字段默认值设为 `''`，不存 NULL，去掉 `nullString()` 转换。

---

## 4. Model 层

### P1 — HTTP 方法常量三处重复定义

`validMethods` 在 Service 中定义，前端 `index.vue` 里硬编码，SQL 种子脚本里也有，
三处不同步时只有运行时才能发现。

**修复方式**：在 `server/internal/model/api/model.go` 中定义：

```go
var ValidHTTPMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
```

Service 引用此变量，前端通过专用接口获取或在 `web/src/types/api/index.ts` 中
以同一来源维护枚举。

### P2 — 状态常量前后端各自定义

`StatusActive=1 / StatusDisabled=2` 在后端 model 中定义，
`web/src/views/system/api/index.vue` 再次硬编码相同值。
后端修改时前端不会报错，只会展示错误的状态文本。

**修复方式**：后端在列表接口或元数据接口中返回可用状态枚举；
或至少将前端常量集中到 `web/src/types/api/index.ts` 并加注释标明与后端对应关系。

---

## 5. 前端 — API 层 & 类型层

### P1 — SaveAPIPayload 中 api_status 应为可选

Handler 在创建时默认将 `api_status` 设为 1，前端创建表单不传该字段，
但类型定义为必填，导致类型与实际使用不匹配。

```typescript
// web/src/types/api/index.ts
export interface SaveAPIPayload {
    api_path: string
    api_method: string
    api_group: string
    api_desc: string
    api_status?: number  // 创建时可选，更新时必填
}
```

---

## 6. 前端 — Views 层

### P1 — 表单验证不完整

`web/src/views/system/api/index.vue` 提交前只 trim 了 `api_path`，
未校验 `api_method` 是否已选，也未验证路径格式（应以 `/` 开头）。

**修复方式**：改用 `el-form` 的 `rules` 声明式验证：

```typescript
const formRules = {
    api_path: [
        { required: true, message: '请填写接口路径' },
        { pattern: /^\//, message: '路径必须以 / 开头' }
    ],
    api_method: [{ required: true, message: '请选择请求方法' }],
    api_group:  [{ required: true, message: '请填写接口分组' }]
}
```

### P2 — 时间格式依赖隐式假设

`new Date(row.created_at).toLocaleString('zh-CN')` 假设后端返回 ISO 8601 格式。
建议封装统一的日期格式化工具函数，集中处理格式兜底逻辑。

---

## 7. 数据库脚本

### P1 — resource_key 未规范化路径参数

`007_seed_api_mgmt_apis.sql` 用 `CONCAT(api_method, ':', api_path)` 生成 resource_key，
但路径中含参数占位符（如 `/api/v1/apis/{id}`）时，
`{id}` 会被当作字面量写入，导致与运行时实际路径不匹配。

**修复方式**：种子脚本中统一使用规范化路径（去掉 `{id}` 等占位符，或改用通配符模式），
与中间件实际使用的 resource_key 格式保持一致。

---

## 优先级汇总

| 优先级 | 问题 | 模块 |
|--------|------|------|
| P0 | pathUint64 零值未拒绝 | Handler |
| P0 | 缺少审计日志 | Service |
| P0 | CreateAPI 竞态条件 | Service |
| P0 | DeleteAPI 孤儿权限策略 | Service |
| P1 | Request 类型未定义到 model | Handler / Model |
| P1 | CreateSkipRule 冗余查询 | Service |
| P1 | DeleteSkipRule 未验证存在性 | Service |
| P1 | LIKE 搜索未转义 | Repository |
| P1 | HTTP 方法常量三处重复 | Model / 前端 |
| P1 | SaveAPIPayload api_status 应可选 | 前端 Types |
| P1 | 表单验证不完整 | 前端 Views |
| P1 | resource_key 路径参数未规范化 | SQL 脚本 |
| P2 | intParam 重复实现 | Handler |
| P2 | 分页魔法数字 | Service |
| P2 | 空字符串与 NULL 处理不统一 | Repository |
| P2 | 状态常量前后端各自定义 | Model / 前端 |
| P2 | 时间格式依赖隐式假设 | 前端 Views |
