-- GoBaseProject 数据库初始化脚本
-- 数据库：go_base_project
-- 执行顺序：建库 → 建表 → 基础数据 → 菜单 → API资源 → 权限策略

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `go_base_project`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `go_base_project`;

-- ── 清理旧表（依赖顺序从叶到根）────────────────────────────────────────────

DROP TABLE IF EXISTS `gbp_files`;
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

-- ── 表结构 ──────────────────────────────────────────────────────────────────

CREATE TABLE `gbp_roles` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`     DATETIME(3)     NULL,
  `role_code`      VARCHAR(64)     NOT NULL,
  `role_name`      VARCHAR(128)    NOT NULL,
  `parent_role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `default_route`  VARCHAR(128)    NOT NULL DEFAULT 'dashboard',
  `sort_no`        INT             NOT NULL DEFAULT 0,
  `role_status`    TINYINT         NOT NULL DEFAULT 1,
  `remark`         VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_roles_role_code` (`role_code`),
  KEY `idx_gbp_roles_parent_role_id` (`parent_role_id`),
  KEY `idx_gbp_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='roles';

CREATE TABLE `gbp_users` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`           DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`           DATETIME(3)     NULL,
  `user_uuid`            VARCHAR(64)     NOT NULL,
  `login_name`           VARCHAR(128)    NOT NULL,
  `password_hash`        VARCHAR(255)    NOT NULL,
  `display_name`         VARCHAR(128)    NOT NULL DEFAULT '系统用户',
  `avatar_url`           VARCHAR(512)    NULL,
  `primary_role_id`      BIGINT UNSIGNED NULL,
  `phone_number`         VARCHAR(32)     NULL,
  `email_address`        VARCHAR(128)    NULL,
  `user_status`          TINYINT         NOT NULL DEFAULT 1,
  `must_change_password` TINYINT(1)      NOT NULL DEFAULT 0,
  `last_login_at`        DATETIME(3)     NULL,
  `profile_config`       JSON            NULL,
  `remark`               VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_users_user_uuid` (`user_uuid`),
  UNIQUE KEY `uk_gbp_users_login_name` (`login_name`),
  KEY `idx_gbp_users_primary_role_id` (`primary_role_id`),
  KEY `idx_gbp_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='users';

CREATE TABLE `gbp_user_roles` (
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `role_id`    BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_gbp_user_roles_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='user roles';

CREATE TABLE `gbp_jwt_blocklist` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`    DATETIME(3)     NULL,
  `token_hash`    CHAR(64)        NOT NULL,
  `expires_at`    DATETIME(3)     NOT NULL,
  `logout_reason` VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_jwt_blocklist_token_hash` (`token_hash`),
  KEY `idx_gbp_jwt_blocklist_expires_at` (`expires_at`),
  KEY `idx_gbp_jwt_blocklist_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='jwt blocklist';

CREATE TABLE `gbp_menus` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`      DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`      DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`      DATETIME(3)     NULL,
  `parent_id`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `menu_type`       TINYINT         NOT NULL DEFAULT 2,
  `route_path`      VARCHAR(255)    NULL,
  `route_name`      VARCHAR(128)    NULL,
  `component_path`  VARCHAR(255)    NULL,
  `redirect_path`   VARCHAR(255)    NULL,
  `menu_title`      VARCHAR(128)    NOT NULL,
  `menu_icon`       VARCHAR(128)    NULL,
  `sort_no`         INT             NOT NULL DEFAULT 0,
  `is_hidden`       TINYINT(1)      NOT NULL DEFAULT 0,
  `is_keep_alive`   TINYINT(1)      NOT NULL DEFAULT 0,
  `is_affix`        TINYINT(1)      NOT NULL DEFAULT 0,
  `active_route`    VARCHAR(128)    NULL,
  `transition_name` VARCHAR(64)     NULL,
  `external_url`    VARCHAR(512)    NULL,
  `menu_status`     TINYINT         NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menus_parent_id` (`parent_id`),
  KEY `idx_gbp_menus_route_name` (`route_name`),
  KEY `idx_gbp_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menus';

CREATE TABLE `gbp_menu_route_params` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL,
  `menu_id`     BIGINT UNSIGNED NOT NULL,
  `param_mode`  VARCHAR(32)     NOT NULL DEFAULT 'query',
  `param_key`   VARCHAR(128)    NOT NULL,
  `param_value` VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_menu_route_params_menu_id` (`menu_id`),
  KEY `idx_gbp_menu_route_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menu route params';

CREATE TABLE `gbp_role_menus` (
  `role_id`    BIGINT UNSIGNED NOT NULL,
  `menu_id`    BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `menu_id`),
  KEY `idx_gbp_role_menus_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role menus';

CREATE TABLE `gbp_menu_actions` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`    DATETIME(3)     NULL,
  `menu_id`       BIGINT UNSIGNED NOT NULL,
  `action_code`   VARCHAR(64)     NOT NULL,
  `action_name`   VARCHAR(128)    NOT NULL,
  `action_desc`   VARCHAR(255)    NULL,
  `sort_no`       INT             NOT NULL DEFAULT 0,
  `action_status` TINYINT         NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_menu_actions_menu_code` (`menu_id`, `action_code`),
  KEY `idx_gbp_menu_actions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='menu actions';

CREATE TABLE `gbp_role_actions` (
  `role_id`    BIGINT UNSIGNED NOT NULL,
  `menu_id`    BIGINT UNSIGNED NOT NULL,
  `action_id`  BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `menu_id`, `action_id`),
  KEY `idx_gbp_role_actions_action_id` (`action_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role actions';

CREATE TABLE `gbp_api_resources` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3)     NULL,
  `api_path`   VARCHAR(255)    NOT NULL,
  `api_method` VARCHAR(16)     NOT NULL DEFAULT 'POST',
  `api_group`  VARCHAR(128)    NULL,
  `api_desc`   VARCHAR(255)    NULL,
  `api_status` TINYINT         NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_resources_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_resources_group` (`api_group`),
  KEY `idx_gbp_api_resources_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='api resources';

CREATE TABLE `gbp_api_skip_rules` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL,
  `api_path`    VARCHAR(255)    NOT NULL,
  `api_method`  VARCHAR(16)     NOT NULL DEFAULT 'POST',
  `skip_reason` VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_api_skip_rules_path_method` (`api_path`, `api_method`),
  KEY `idx_gbp_api_skip_rules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='api skip rules';

CREATE TABLE `gbp_permission_policies` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`     DATETIME(3)     NULL,
  `subject_type`   VARCHAR(32)     NOT NULL DEFAULT 'role',
  `subject_id`     BIGINT UNSIGNED NOT NULL,
  `resource_type`  VARCHAR(32)     NOT NULL DEFAULT 'api',
  `resource_id`    BIGINT UNSIGNED NULL,
  `resource_key`   VARCHAR(255)    NOT NULL,
  `action`         VARCHAR(64)     NOT NULL DEFAULT '*',
  `effect`         VARCHAR(16)     NOT NULL DEFAULT 'allow',
  `condition_expr` VARCHAR(512)    NULL,
  `policy_status`  TINYINT         NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_permission_policies` (`subject_type`, `subject_id`, `resource_type`, `resource_key`, `action`),
  KEY `idx_gbp_permission_policies_subject` (`subject_type`, `subject_id`),
  KEY `idx_gbp_permission_policies_resource` (`resource_type`, `resource_id`),
  KEY `idx_gbp_permission_policies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='permission policies';

CREATE TABLE `gbp_role_data_scopes` (
  `role_id`         BIGINT UNSIGNED NOT NULL,
  `visible_role_id` BIGINT UNSIGNED NOT NULL,
  `created_at`      DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_id`, `visible_role_id`),
  KEY `idx_gbp_role_data_scopes_visible_role_id` (`visible_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='role data scopes';

CREATE TABLE `gbp_dictionaries` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL,
  `dict_name`   VARCHAR(128)    NOT NULL,
  `dict_code`   VARCHAR(128)    NOT NULL,
  `dict_status` TINYINT         NOT NULL DEFAULT 1,
  `parent_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `remark`      VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionaries_dict_code` (`dict_code`),
  KEY `idx_gbp_dictionaries_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionaries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='dictionaries';

CREATE TABLE `gbp_dictionary_items` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL,
  `dict_id`     BIGINT UNSIGNED NOT NULL,
  `item_label`  VARCHAR(128)    NOT NULL,
  `item_value`  VARCHAR(128)    NOT NULL,
  `item_extra`  VARCHAR(255)    NULL,
  `item_status` TINYINT         NOT NULL DEFAULT 1,
  `sort_no`     INT             NOT NULL DEFAULT 0,
  `parent_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `tree_level`  INT             NOT NULL DEFAULT 1,
  `tree_path`   VARCHAR(512)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_dictionary_items_dict_value` (`dict_id`, `item_value`),
  KEY `idx_gbp_dictionary_items_parent_id` (`parent_id`),
  KEY `idx_gbp_dictionary_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='dictionary items';

CREATE TABLE `gbp_system_params` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`   DATETIME(3)     NULL,
  `param_name`   VARCHAR(128)    NOT NULL,
  `param_key`    VARCHAR(128)    NOT NULL,
  `param_value`  TEXT            NULL,
  `param_type`   VARCHAR(32)     NOT NULL DEFAULT 'string',
  `is_encrypted` TINYINT(1)      NOT NULL DEFAULT 0,
  `param_status` TINYINT         NOT NULL DEFAULT 1,
  `remark`       VARCHAR(255)    NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_gbp_system_params_param_key` (`param_key`),
  KEY `idx_gbp_system_params_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='system params';

CREATE TABLE `gbp_login_audit_logs` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`    DATETIME(3)     NULL,
  `user_id`       BIGINT UNSIGNED NULL,
  `login_name`    VARCHAR(128)    NULL,
  `source_ip`     VARCHAR(64)     NULL,
  `login_success` TINYINT(1)      NOT NULL DEFAULT 0,
  `fail_reason`   VARCHAR(255)    NULL,
  `user_agent`    VARCHAR(512)    NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_login_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_login_audit_logs_login_name` (`login_name`),
  KEY `idx_gbp_login_audit_logs_source_ip` (`source_ip`),
  KEY `idx_gbp_login_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_login_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='login audit logs';

CREATE TABLE `gbp_operation_audit_logs` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`     DATETIME(3)     NULL,
  `user_id`        BIGINT UNSIGNED NULL,
  `source_ip`      VARCHAR(64)     NULL,
  `request_method` VARCHAR(16)     NULL,
  `request_path`   VARCHAR(255)    NULL,
  `status_code`    INT             NULL,
  `cost_ms`        BIGINT          NULL,
  `user_agent`     TEXT            NULL,
  `error_message`  VARCHAR(512)    NULL,
  `request_body`   LONGTEXT        NULL,
  `response_body`  LONGTEXT        NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gbp_operation_audit_logs_user_id` (`user_id`),
  KEY `idx_gbp_operation_audit_logs_request_path` (`request_path`),
  KEY `idx_gbp_operation_audit_logs_status_code` (`status_code`),
  KEY `idx_gbp_operation_audit_logs_created_at` (`created_at`),
  KEY `idx_gbp_operation_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='operation audit logs';

CREATE TABLE `gbp_files` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `original_name` VARCHAR(255)    NOT NULL COMMENT '原始文件名',
  `storage_key`   VARCHAR(512)    NOT NULL COMMENT '磁盘存储相对路径，格式 YYYY/MM/DD/<uuid>.<ext>',
  `file_size`     BIGINT          NOT NULL COMMENT '文件大小（字节）',
  `mime_type`     VARCHAR(128)    NOT NULL DEFAULT '' COMMENT 'MIME 类型',
  `uploader_id`   BIGINT UNSIGNED NULL COMMENT '上传用户 ID',
  `created_at`    DATETIME(3)     NOT NULL DEFAULT NOW(3),
  `deleted_at`    DATETIME(3)     NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_storage_key` (`storage_key`),
  KEY `idx_uploader` (`uploader_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件记录';

-- ── 基础角色与用户 ──────────────────────────────────────────────────────────

INSERT INTO `gbp_roles`
  (`id`, `role_code`, `role_name`, `parent_role_id`, `default_route`, `sort_no`, `role_status`, `remark`)
VALUES
  (1, 'super_admin',   '超级管理员', 0, 'dashboard', 1, 1, '系统内置最高权限角色'),
  (2, 'system_admin',  '系统管理员', 1, 'dashboard', 2, 1, '负责用户、角色、菜单、权限等配置'),
  (3, 'operator',      '运营人员',   2, 'dashboard', 3, 1, '普通运营人员角色');

-- 初始密码：Admin@123456
INSERT INTO `gbp_users`
  (`id`, `user_uuid`, `login_name`, `password_hash`, `display_name`, `primary_role_id`, `user_status`, `must_change_password`, `remark`)
VALUES
  (1, '00000000-0000-0000-0000-000000000001', 'admin',
   '$2a$10$26jSKICOGGsyliE4t2rHNuCrD2Hc9pCd/RSQ4mIx7NPMEEkDAs5kK',
   '超级管理员', 1, 1, 1, '系统初始化管理员');

INSERT INTO `gbp_user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- ── 菜单 ─────────────────────────────────────────────────────────────────────

INSERT INTO `gbp_menus`
  (`id`, `parent_id`, `menu_type`, `route_path`, `route_name`, `component_path`,
   `redirect_path`, `menu_title`, `menu_icon`, `sort_no`, `is_hidden`, `is_keep_alive`, `menu_status`)
VALUES
  -- 首页
  (1,  0,  1, '/dashboard', 'dashboard',        'layouts/default',                   '/dashboard/workbench', '首页',     'HomeFilled',  1,  0, 1, 1),
  (2,  1,  2, 'workbench',  'dashboardWorkbench','views/dashboard/workbench/index.vue', NULL,                '工作台',   'DataBoard',   1,  0, 1, 1),
  -- 系统管理
  (10, 0,  1, '/system',    'system',            'layouts/default',                   '/system/user',        '系统管理', 'Setting',    10,  0, 0, 1),
  (11, 10, 2, 'user',       'systemUser',        'views/system/user/index.vue',        NULL,                 '用户管理', 'User',        1,  0, 0, 1),
  (12, 10, 2, 'role',       'systemRole',        'views/system/role/index.vue',        NULL,                 '角色管理', 'Avatar',      2,  0, 0, 1),
  (13, 10, 2, 'menu',       'systemMenu',        'views/system/menu/index.vue',        NULL,                 '菜单管理', 'Menu',        3,  0, 0, 1),
  (14, 10, 2, 'api',        'systemApi',         'views/system/api/index.vue',         NULL,                 'API管理',  'Link',        4,  0, 0, 1),
  (15, 10, 2, 'dictionary', 'systemDictionary',  'views/system/dict/index.vue',        NULL,                 '字典管理', 'Collection',  5,  0, 0, 1),
  (16, 10, 2, 'param',      'systemParam',       'views/system/param/index.vue',       NULL,                 '系统参数', 'Tools',       6,  0, 0, 1),
  -- 审计日志
  (20, 0,  1, '/audit',     'audit',             'layouts/default',                   '/audit/login-log',    '审计日志', 'Document',   20,  0, 0, 1),
  (21, 20, 2, 'login-log',  'auditLoginLog',     'views/audit/login-log/index.vue',    NULL,                 '登录日志', 'Tickets',     1,  0, 0, 1),
  (22, 20, 2, 'operation-log','auditOperationLog','views/audit/operation-log/index.vue',NULL,                 '操作日志', 'Notebook',    2,  0, 0, 1),
  -- 文件管理
  (30, 0,  1, '/file',      'file',              'layouts/default',                   '/file/list',          '文件管理', 'FolderOpened',30, 0, 0, 1),
  (31, 30, 2, 'list',       'fileList',          'views/system/file/index.vue',        NULL,                 '文件列表', 'Files',       1,  0, 0, 1);

INSERT INTO `gbp_role_menus` (`role_id`, `menu_id`)
SELECT 1, `id` FROM `gbp_menus`;

-- ── 权限白名单（无需鉴权的接口）────────────────────────────────────────────

INSERT INTO `gbp_api_skip_rules` (`api_path`, `api_method`, `skip_reason`) VALUES
  ('/api/v1/auth/login',                    'POST', '登录接口'),
  ('/api/v1/auth/refresh',                  'POST', '刷新Token接口'),
  ('/api/v1/health',                        'GET',  '健康检查接口'),
  ('/api/v1/dictionaries/{dictCode}/items', 'GET',  '字典项查询，所有登录用户均可访问');

-- ── 字典 ─────────────────────────────────────────────────────────────────────

INSERT INTO `gbp_dictionaries` (`id`, `dict_name`, `dict_code`, `dict_status`, `remark`) VALUES
  (1, '用户状态', 'user_status', 1, '用户启用/冻结状态'),
  (2, '角色状态', 'role_status', 1, '角色启用/禁用状态'),
  (3, '菜单类型', 'menu_type',   1, '目录、菜单、隐藏路由、外链'),
  (4, '参数类型', 'param_type',  1, '系统参数值类型');

INSERT INTO `gbp_dictionary_items` (`dict_id`, `item_label`, `item_value`, `sort_no`) VALUES
  (1, '正常',   '1',       1),
  (1, '冻结',   '2',       2),
  (2, '启用',   '1',       1),
  (2, '禁用',   '2',       2),
  (3, '目录',   '1',       1),
  (3, '菜单',   '2',       2),
  (3, '隐藏路由','3',      3),
  (3, '外链',   '4',       4),
  (4, '字符串', 'string',  1),
  (4, '数字',   'number',  2),
  (4, '布尔值', 'boolean', 3),
  (4, 'JSON',   'json',    4);

-- ── 系统参数 ─────────────────────────────────────────────────────────────────

INSERT INTO `gbp_system_params`
  (`param_name`, `param_key`, `param_value`, `param_type`, `is_encrypted`, `param_status`, `remark`)
VALUES
  ('系统名称',             'system.name',                    'GoBaseProject', 'string',  0, 1, '前端页面展示的系统名称'),
  ('访问Token有效期',      'auth.access_token_ttl_minutes',  '120',           'number',  0, 1, '单位：分钟'),
  ('刷新Token有效期',      'auth.refresh_token_ttl_minutes', '10080',         'number',  0, 1, '单位：分钟'),
  ('登录失败锁定开关',     'security.login_lock_enabled',    'true',          'boolean', 0, 1, '是否启用登录失败锁定'),
  ('操作日志记录响应体',   'audit.record_response_body',     'false',         'boolean', 0, 1, '是否记录接口响应体');

-- ── API 资源 ─────────────────────────────────────────────────────────────────
-- path 中的路径参数统一用 {id}，与中间件的 normalizePath 对应

INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  -- auth
  ('/api/v1/auth/logout',                   'POST',   'auth',   '退出登录',           1),
  ('/api/v1/auth/profile',                  'GET',    'auth',   '当前用户信息',       1),
  ('/api/v1/auth/routes',                   'GET',    'auth',   '菜单路由',           1),
  ('/api/v1/auth/actions',                  'GET',    'auth',   '按钮权限',           1),
  -- user
  ('/api/v1/users',                         'GET',    'user',   '用户列表',           1),
  ('/api/v1/users',                         'POST',   'user',   '创建用户',           1),
  ('/api/v1/users/role-options',            'GET',    'user',   '角色选项',           1),
  ('/api/v1/users/{id}',                    'GET',    'user',   '用户详情',           1),
  ('/api/v1/users/{id}',                    'PUT',    'user',   '修改用户',           1),
  ('/api/v1/users/{id}',                    'DELETE', 'user',   '删除用户',           1),
  ('/api/v1/users/{id}/status',             'PUT',    'user',   '修改用户状态',       1),
  ('/api/v1/users/{id}/password',           'PUT',    'user',   '重置密码',           1),
  ('/api/v1/users/{id}/roles',              'PUT',    'user',   '分配角色',           1),
  -- role
  ('/api/v1/roles',                         'GET',    'role',   '角色列表',           1),
  ('/api/v1/roles',                         'POST',   'role',   '创建角色',           1),
  ('/api/v1/roles/tree',                    'GET',    'role',   '角色树',             1),
  ('/api/v1/roles/resources',               'GET',    'role',   '权限资源目录',       1),
  ('/api/v1/roles/{id}',                    'GET',    'role',   '角色详情',           1),
  ('/api/v1/roles/{id}',                    'PUT',    'role',   '修改角色',           1),
  ('/api/v1/roles/{id}',                    'DELETE', 'role',   '删除角色',           1),
  ('/api/v1/roles/{id}/menus',              'GET',    'role',   '角色菜单',           1),
  ('/api/v1/roles/{id}/menus',              'PUT',    'role',   '分配菜单',           1),
  ('/api/v1/roles/{id}/actions',            'GET',    'role',   '角色按钮权限',       1),
  ('/api/v1/roles/{id}/actions',            'PUT',    'role',   '分配按钮权限',       1),
  ('/api/v1/roles/{id}/apis',               'GET',    'role',   '角色API权限',        1),
  ('/api/v1/roles/{id}/apis',               'PUT',    'role',   '分配API权限',        1),
  ('/api/v1/roles/{id}/data-scopes',        'GET',    'role',   '角色数据范围',       1),
  ('/api/v1/roles/{id}/data-scopes',        'PUT',    'role',   '分配数据范围',       1),
  -- menu
  ('/api/v1/menus',                         'GET',    'menu',   '菜单列表',           1),
  ('/api/v1/menus',                         'POST',   'menu',   '创建菜单',           1),
  ('/api/v1/menus/tree',                    'GET',    'menu',   '菜单树',             1),
  ('/api/v1/menus/{id}',                    'GET',    'menu',   '菜单详情',           1),
  ('/api/v1/menus/{id}',                    'PUT',    'menu',   '修改菜单',           1),
  ('/api/v1/menus/{id}',                    'DELETE', 'menu',   '删除菜单',           1),
  ('/api/v1/menus/{id}/params',             'GET',    'menu',   '菜单路由参数',       1),
  ('/api/v1/menus/{id}/params',             'POST',   'menu',   '新增路由参数',       1),
  ('/api/v1/menus/params/{id}',             'DELETE', 'menu',   '删除路由参数',       1),
  ('/api/v1/menus/{id}/actions',            'GET',    'menu',   '菜单按钮权限',       1),
  ('/api/v1/menus/{id}/actions',            'POST',   'menu',   '新增按钮权限',       1),
  ('/api/v1/menu-actions/{id}',             'PUT',    'menu',   '修改按钮权限',       1),
  ('/api/v1/menu-actions/{id}',             'DELETE', 'menu',   '删除按钮权限',       1),
  -- system
  ('/api/v1/system/runtime',                'GET',    'system', '系统运行信息',       1),
  -- dict
  ('/api/v1/dictionaries',                  'GET',    'dict',   '字典列表',           1),
  ('/api/v1/dictionaries',                  'POST',   'dict',   '创建字典',           1),
  ('/api/v1/dictionaries/{id}',             'PUT',    'dict',   '修改字典',           1),
  ('/api/v1/dictionaries/{id}',             'DELETE', 'dict',   '删除字典',           1),
  ('/api/v1/dictionaries/{id}/items',       'POST',   'dict',   '新增字典项',         1),
  ('/api/v1/dictionaries/{dictCode}/items', 'GET',    'dict',   '查询字典项',         1),
  ('/api/v1/dictionary-items/{id}',         'PUT',    'dict',   '修改字典项',         1),
  ('/api/v1/dictionary-items/{id}',         'DELETE', 'dict',   '删除字典项',         1),
  -- audit
  ('/api/v1/audit/login-logs',              'GET',    'audit',  '登录日志列表',       1),
  ('/api/v1/audit/login-logs/cleanup',      'DELETE', 'audit',  '登录日志清理',       1),
  ('/api/v1/audit/operation-logs',          'GET',    'audit',  '操作日志列表',       1),
  ('/api/v1/audit/operation-logs/{id}',     'GET',    'audit',  '操作日志详情',       1),
  ('/api/v1/audit/operation-logs/cleanup',  'DELETE', 'audit',  '操作日志清理',       1),
  -- api management
  ('/api/v1/apis/groups',                   'GET',    'api',    'API分组列表',        1),
  ('/api/v1/apis',                          'GET',    'api',    'API资源列表',        1),
  ('/api/v1/apis',                          'POST',   'api',    '创建API资源',        1),
  ('/api/v1/apis/{id}',                     'PUT',    'api',    '修改API资源',        1),
  ('/api/v1/apis/{id}',                     'DELETE', 'api',    '删除API资源',        1),
  ('/api/v1/api-skip-rules',                'GET',    'api',    '白名单列表',         1),
  ('/api/v1/api-skip-rules',                'POST',   'api',    '创建白名单',         1),
  ('/api/v1/api-skip-rules/{id}',           'DELETE', 'api',    '删除白名单',         1),
  -- file
  ('/api/v1/files',                         'POST',   'file',   '上传文件',           1),
  ('/api/v1/files',                         'GET',    'file',   '文件列表',           1),
  ('/api/v1/files/{id}',                    'DELETE', 'file',   '删除文件',           1),
  ('/api/v1/files/{id}/raw',                'GET',    'file',   '下载/预览文件',      1);

-- ── 超级管理员 API 权限策略（覆盖全部接口）──────────────────────────────────

INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT 'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `deleted_at` IS NULL;

SET FOREIGN_KEY_CHECKS = 1;
