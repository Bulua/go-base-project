# 06-API权限模块

## 1. 模块职责

负责维护系统 API 资源、配置鉴权规则，并为角色分配 API 权限；请求进入时由接口鉴权中间件校验。

## 2. 涉及数据表

```text
gbp_api_resources
gbp_api_skip_rules
gbp_permission_policies
gbp_roles
```

## 3. 核心接口

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /api/v1/apis | API 资源列表 |
| POST | /api/v1/apis | 创建 API 资源 |
| PUT | /api/v1/apis/:id | 修改 API 资源 |
| DELETE | /api/v1/apis/:id | 删除 API 资源 |
| GET | /api/v1/api-skip-rules | API 白名单列表 |
| POST | /api/v1/api-skip-rules | 创建白名单 |
| DELETE | /api/v1/api-skip-rules/:id | 删除白名单 |
| GET | /api/v1/roles/:id/apis | 查询角色 API 权限 |
| PUT | /api/v1/roles/:id/apis | 分配角色 API 权限 |

## 4. 关键流程

```text
请求进入 -> 命中白名单则放行 -> 解析 JWT 角色 -> 生成 resource_key -> 查询权限策略 -> allow 放行 / deny 拒绝 / 未配置 403
```

## 5. 开发注意事项

1. resource_key 推荐格式：METHOD:/api/v1/users。
2. deny 策略优先级高于 allow。
3. 第一阶段可以手动维护 API，第二阶段可做路由扫描自动注册。
