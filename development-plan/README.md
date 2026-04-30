# GoBaseProject 开发方案包

本目录包含 GoBaseProject 自研后台基座第一阶段开发方案、分模块文档、数据库 SQL、后端落地规范和前端设计示例。

## 文件结构

```text
development-plan
├── docs
├── server-skeleton
├── sql
│   ├── 01_schema.sql
│   └── 02_seed.sql
└── web-design
```

## 使用顺序

```bash
mysql -uroot -p < sql/01_schema.sql
mysql -uroot -p go_base_project < sql/02_seed.sql
```

开发建议顺序：

```text
认证鉴权 -> 用户管理 -> 角色管理 -> 菜单路由 -> API权限 -> 按钮权限 -> 数据权限 -> 字典参数 -> 审计日志
```

## 开发文档

后端开发先看 `docs/00-项目总览.md` 和 `docs/11-后端开发规范.md`，前端开发先看 `docs/12-前端开发规范.md`，再按模块文档拆分实现。接口总览见 `docs/13-接口清单.md`，前端页面原型参考 `web-design/*`。
