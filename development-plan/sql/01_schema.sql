SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `go_base_project`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `go_base_project`;

-- =========================================================
-- GoBaseProject - 第一阶段数据库结构
-- MySQL 5.7+/8.x/MariaDB 通用版本。
-- =========================================================

DROP TABLE IF EXISTS `gbp_operation_audit_logs`;
DROP TABLE IF EXISTS `gbp_login_audit_logs`;
DROP TABLE IF EXISTS `gbp_system_params`;
DROP TABLE IF EXISTS `gbp_dictionary_items`;
DROP TABLE IF EXISTS `gbp_dictionaries`;
DROP TABLE IF EXISTS `gbp_role_data_scopes`;
DROP TABLE IF EXISTS `gbp_permission_policies`;
DROP TABLE IF EXISTS `gbp_api_skip_rules`;
DROP TABLE IF EXISTS `gbp_api_resources`;
DROP TABLE IF EXISTS `gbp_role_actions`;
DROP TABLE IF EXISTS `gbp_menu_actions`;
DROP TABLE IF EXISTS `gbp_role_menus`;
DROP TABLE IF EXISTS `gbp_menu_route_params`;
DROP TABLE IF EXISTS `gbp_menus`;
DROP TABLE IF EXISTS `gbp_jwt_blocklist`;
DROP TABLE IF EXISTS `gbp_user_roles`;
DROP TABLE IF EXISTS `gbp_users`;
DROP TABLE IF EXISTS `gbp_roles`;

CREATE TABLE `gbp_roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `role_code` VARCHAR(64) NOT NULL COMMENT '角色编码，例如 super_admin/system_admin/operator',
  `role_name` VARCHAR(128) NOT NULL COMMENT '角色名称',
  `parent_role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父角色ID，0表示无父级',
  `default_route` VARCHAR(128) NOT NULL DEFAULT 'dashboard' COMMENT '默认进入路由',
  `sort_no` INT NOT NULL DEFAULT 0 COMMENT '排序号',
  `role_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  `remark` VARCHAR(255) NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_roles_role_code` (`role_code`),
  KEY `idx_gbp_roles_parent_role_id` (`parent_role_id`),
  KEY `idx_gbp_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';

CREATE TABLE `gbp_users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `user_uuid` VARCHAR(64) NOT NULL COMMENT '用户UUID',
  `login_name` VARCHAR(128) NOT NULL COMMENT '登录账号',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希，推荐 bcrypt/argon2id',
  `display_name` VARCHAR(128) NOT NULL DEFAULT '系统用户' COMMENT '显示名称',
  `avatar_url` VARCHAR(512) NULL COMMENT '头像地址',
  `primary_role_id` BIGINT UNSIGNED NULL COMMENT '默认角色ID',
  `phone_number` VARCHAR(32) NULL COMMENT '手机号',
  `email_address` VARCHAR(128) NULL COMMENT '邮箱地址',
  `user_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1正常，2冻结',
  `must_change_password` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否必须修改密码：0否，1是',
  `last_login_at` DATETIME(3) NULL COMMENT '最后登录时间',
  `profile_config` JSON NULL COMMENT '用户个性化配置',
  `remark` VARCHAR(255) NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_users_user_uuid` (`user_uuid`),
  UNIQUE KEY `uk_gbp_users_login_name` (`login_name`),
  KEY `idx_gbp_users_primary_role_id` (`primary_role_id`),
  KEY `idx_gbp_users_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_gbp_users_primary_role` FOREIGN KEY (`primary_role_id`) REFERENCES `gbp_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE `gbp_user_roles` (
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '绑定时间',
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_gbp_user_roles_role_id` (`role_id`),
  CONSTRAINT `fk_gbp_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `gbp_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_gbp_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `gbp_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色关联表';

CREATE TABLE `gbp_jwt_blocklist` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `token_hash` CHAR(64) NOT NULL COMMENT 'JWT SHA-256 摘要，不存原始JWT',
  `expires_at` DATETIME(3) NOT NULL COMMENT 'Token 原始过期时间',
  `logout_reason` VARCHAR(255) NULL COMMENT '拉黑原因：主动退出、强制下线、密码修改等',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_jwt_blocklist_token_hash` (`token_hash`),
  KEY `idx_gbp_jwt_blocklist_expires_at` (`expires_at`),
  KEY `idx_gbp_jwt_blocklist_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='JWT黑名单表';

CREATE TABLE `gbp_menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示根节点',
  `menu_type` TINYINT NOT NULL DEFAULT 2 COMMENT '类型：1目录，2菜单，3隐藏路由，4外链',
  `route_path` VARCHAR(255) NULL COMMENT '前端路由路径',
  `route_name` VARCHAR(128) NULL COMMENT '前端路由名称',
  `component_path` VARCHAR(255) NULL COMMENT '前端组件路径',
  `redirect_path` VARCHAR(255) NULL COMMENT '重定向路径',
  `menu_title` VARCHAR(128) NOT NULL COMMENT '菜单标题',
  `menu_icon` VARCHAR(128) NULL COMMENT '菜单图标',
  `sort_no` INT NOT NULL DEFAULT 0 COMMENT '排序号',
  `is_hidden` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏：0显示，1隐藏',
  `is_keep_alive` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否缓存：0否，1是',
  `is_affix` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否固定标签页：0否，1是',
  `active_route` VARCHAR(128) NULL COMMENT '当前路由高亮菜单',
  `transition_name` VARCHAR(64) NULL COMMENT '路由切换动画',
  `external_url` VARCHAR(512) NULL COMMENT '外链地址',
  `menu_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menus_parent_id` (`parent_id`),
  KEY `idx_gbp_menus_route_name` (`route_name`),
  KEY `idx_gbp_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单与动态路由表';

CREATE TABLE `gbp_menu_route_params` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '参数ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '菜单ID',
  `param_mode` VARCHAR(32) NOT NULL DEFAULT 'query' COMMENT '参数类型：query 或 params',
  `param_key` VARCHAR(128) NOT NULL COMMENT '参数名',
  `param_value` VARCHAR(255) NULL COMMENT '参数值',
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menu_route_params_menu_id` (`menu_id`),
  KEY `idx_gbp_menu_route_params_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_gbp_menu_route_params_menu` FOREIGN KEY (`menu_id`) REFERENCES `gbp_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单路由参数表';

CREATE TABLE `gbp_role_menus` (
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '菜单ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '授权时间',
  PRIMARY KEY (`role_id`, `menu_id`),
  KEY `idx_gbp_role_menus_menu_id` (`menu_id`),
  CONSTRAINT `fk_gbp_role_menus_role` FOREIGN KEY (`role_id`) REFERENCES `gbp_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_gbp_role_menus_menu` FOREIGN KEY (`menu_id`) REFERENCES `gbp_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色菜单权限表';

CREATE TABLE `gbp_menu_actions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '按钮动作ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '所属菜单ID',
  `action_code` VARCHAR(64) NOT NULL COMMENT '按钮权限标识，例如 add/edit/delete/export',
  `action_name` VARCHAR(128) NOT NULL COMMENT '按钮显示名称',
  `action_desc` VARCHAR(255) NULL COMMENT '按钮说明',
  `sort_no` INT NOT NULL DEFAULT 0 COMMENT '排序号',
  `action_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_menu_actions_menu_code` (`menu_id`, `action_code`),
  KEY `idx_gbp_menu_actions_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_gbp_menu_actions_menu` FOREIGN KEY (`menu_id`) REFERENCES `gbp_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单按钮动作定义表';

CREATE TABLE `gbp_role_actions` (
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '菜单ID',
  `action_id` BIGINT UNSIGNED NOT NULL COMMENT '按钮动作ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '授权时间',
  PRIMARY KEY (`role_id`, `menu_id`, `action_id`),
  KEY `idx_gbp_role_actions_action_id` (`action_id`),
  CONSTRAINT `fk_gbp_role_actions_role` FOREIGN KEY (`role_id`) REFERENCES `gbp_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_gbp_role_actions_menu` FOREIGN KEY (`menu_id`) REFERENCES `gbp_menus` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_gbp_role_actions_action` FOREIGN KEY (`action_id`) REFERENCES `gbp_menu_actions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色按钮权限表';

CREATE TABLE `gbp_api_resources` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'API资源ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `api_path` VARCHAR(255) NOT NULL COMMENT 'API路径，例如 /api/v1/users',
  `api_method` VARCHAR(16) NOT NULL DEFAULT 'POST' COMMENT '请求方法',
  `api_group` VARCHAR(128) NULL COMMENT 'API分组，例如 用户管理',
  `api_desc` VARCHAR(255) NULL COMMENT 'API描述',
  `api_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_resources_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_resources_group` (`api_group`),
  KEY `idx_gbp_api_resources_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='API资源表';

CREATE TABLE `gbp_api_skip_rules` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `api_path` VARCHAR(255) NOT NULL COMMENT '忽略鉴权的API路径',
  `api_method` VARCHAR(16) NOT NULL DEFAULT 'POST' COMMENT '请求方法',
  `skip_reason` VARCHAR(255) NULL COMMENT '忽略原因',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_skip_rules_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_skip_rules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='API鉴权忽略规则表';

CREATE TABLE `gbp_permission_policies` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '策略ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `subject_type` VARCHAR(32) NOT NULL DEFAULT 'role' COMMENT '主体类型：role/user',
  `subject_id` BIGINT UNSIGNED NOT NULL COMMENT '主体ID，通常为角色ID',
  `resource_type` VARCHAR(32) NOT NULL DEFAULT 'api' COMMENT '资源类型：api/menu/button',
  `resource_id` BIGINT UNSIGNED NULL COMMENT '资源ID，例如 API ID',
  `resource_key` VARCHAR(255) NOT NULL COMMENT '资源标识，例如 METHOD:/api/v1/users',
  `action` VARCHAR(64) NOT NULL DEFAULT '*' COMMENT '动作，例如 GET/POST/read/write/*',
  `effect` VARCHAR(16) NOT NULL DEFAULT 'allow' COMMENT '效果：allow/deny',
  `condition_expr` VARCHAR(512) NULL COMMENT '条件表达式，第一阶段可为空',
  `policy_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_permission_policies` (`subject_type`, `subject_id`, `resource_type`, `resource_key`, `action`),
  KEY `idx_gbp_permission_policies_subject` (`subject_type`, `subject_id`),
  KEY `idx_gbp_permission_policies_resource` (`resource_type`, `resource_id`),
  KEY `idx_gbp_permission_policies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='权限策略表';

CREATE TABLE `gbp_role_data_scopes` (
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '当前角色ID',
  `visible_role_id` BIGINT UNSIGNED NOT NULL COMMENT '可见数据角色ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '配置时间',
  PRIMARY KEY (`role_id`, `visible_role_id`),
  KEY `idx_gbp_role_data_scopes_visible_role_id` (`visible_role_id`),
  CONSTRAINT `fk_gbp_role_data_scopes_role` FOREIGN KEY (`role_id`) REFERENCES `gbp_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_gbp_role_data_scopes_visible_role` FOREIGN KEY (`visible_role_id`) REFERENCES `gbp_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色数据权限范围表';

CREATE TABLE `gbp_dictionaries` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `dict_name` VARCHAR(128) NOT NULL COMMENT '字典名称',
  `dict_code` VARCHAR(128) NOT NULL COMMENT '字典编码',
  `dict_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级字典ID',
  `remark` VARCHAR(255) NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionaries_dict_code` (`dict_code`),
  KEY `idx_gbp_dictionaries_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionaries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='字典主表';

CREATE TABLE `gbp_dictionary_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典项ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `dict_id` BIGINT UNSIGNED NOT NULL COMMENT '所属字典ID',
  `item_label` VARCHAR(128) NOT NULL COMMENT '显示名称',
  `item_value` VARCHAR(128) NOT NULL COMMENT '字典值',
  `item_extra` VARCHAR(255) NULL COMMENT '扩展值',
  `item_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  `sort_no` INT NOT NULL DEFAULT 0 COMMENT '排序号',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级字典项ID',
  `tree_level` INT NOT NULL DEFAULT 1 COMMENT '树层级',
  `tree_path` VARCHAR(512) NULL COMMENT '树路径',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionary_items_dict_value` (`dict_id`, `item_value`),
  KEY `idx_gbp_dictionary_items_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionary_items_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_gbp_dictionary_items_dict` FOREIGN KEY (`dict_id`) REFERENCES `gbp_dictionaries` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='字典项表';

CREATE TABLE `gbp_system_params` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '参数ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `param_name` VARCHAR(128) NOT NULL COMMENT '参数名称',
  `param_key` VARCHAR(128) NOT NULL COMMENT '参数键',
  `param_value` TEXT NULL COMMENT '参数值',
  `param_type` VARCHAR(32) NOT NULL DEFAULT 'string' COMMENT '参数类型：string/number/boolean/json',
  `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否加密存储：0否，1是',
  `param_status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，2禁用',
  `remark` VARCHAR(255) NULL COMMENT '说明',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_system_params_param_key` (`param_key`),
  KEY `idx_gbp_system_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统参数表';

CREATE TABLE `gbp_login_audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `user_id` BIGINT UNSIGNED NULL COMMENT '用户ID，登录失败时可能为空',
  `login_name` VARCHAR(128) NULL COMMENT '登录账号',
  `source_ip` VARCHAR(64) NULL COMMENT '来源IP',
  `login_success` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否成功：0失败，1成功',
  `fail_reason` VARCHAR(255) NULL COMMENT '失败原因',
  `user_agent` VARCHAR(512) NULL COMMENT 'User-Agent',
  PRIMARY KEY (`id`),
  KEY `idx_gbp_login_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_login_audit_logs_login_name` (`login_name`),
  KEY `idx_gbp_login_audit_logs_source_ip` (`source_ip`),
  KEY `idx_gbp_login_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_login_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='登录审计日志表';

CREATE TABLE `gbp_operation_audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL COMMENT '删除时间，软删除标记',
  `user_id` BIGINT UNSIGNED NULL COMMENT '操作用户ID',
  `source_ip` VARCHAR(64) NULL COMMENT '来源IP',
  `request_method` VARCHAR(16) NULL COMMENT '请求方法',
  `request_path` VARCHAR(255) NULL COMMENT '请求路径',
  `status_code` INT NULL COMMENT '响应状态码',
  `cost_ms` BIGINT NULL COMMENT '请求耗时，单位毫秒',
  `user_agent` TEXT NULL COMMENT 'User-Agent',
  `error_message` VARCHAR(512) NULL COMMENT '错误信息',
  `request_body` LONGTEXT NULL COMMENT '请求体',
  `response_body` LONGTEXT NULL COMMENT '响应体',
  PRIMARY KEY (`id`),
  KEY `idx_gbp_operation_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_operation_audit_logs_request_path` (`request_path`),
  KEY `idx_gbp_operation_audit_logs_status_code` (`status_code`),
  KEY `idx_gbp_operation_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_operation_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='操作审计日志表';

SET FOREIGN_KEY_CHECKS = 1;
