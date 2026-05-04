SET NAMES utf8mb4;
USE `go_base_project`;

INSERT INTO `gbp_roles`
  (`id`, `role_code`, `role_name`, `parent_role_id`, `default_route`, `sort_no`, `role_status`, `remark`)
VALUES
  (1, 'super_admin', '超级管理员', 0, 'dashboard', 1, 1, '系统内置最高权限角色'),
  (2, 'system_admin', '系统管理员', 1, 'dashboard', 2, 1, '负责用户、角色、菜单、权限等配置'),
  (3, 'operator', '运营人员', 2, 'dashboard', 3, 1, '普通运营人员角色');

INSERT INTO `gbp_users`
  (`id`, `user_uuid`, `login_name`, `password_hash`, `display_name`, `primary_role_id`, `user_status`, `must_change_password`, `remark`)
VALUES
  (1, '00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$26jSKICOGGsyliE4t2rHNuCrD2Hc9pCd/RSQ4mIx7NPMEEkDAs5kK', '超级管理员', 1, 1, 1, '系统初始化管理员');

INSERT INTO `gbp_user_roles` (`user_id`, `role_id`) VALUES
  (1, 1);

INSERT INTO `gbp_menus`
  (`id`, `parent_id`, `menu_type`, `route_path`, `route_name`, `component_path`, `redirect_path`, `menu_title`, `menu_icon`, `sort_no`, `is_hidden`, `is_keep_alive`, `menu_status`)
VALUES
  (1, 0, 1, '/dashboard', 'dashboard', 'layouts/default', '/dashboard/workbench', '首页', 'HomeFilled', 1, 0, 1, 1),
  (2, 1, 2, 'workbench', 'dashboardWorkbench', 'views/dashboard/workbench/index.vue', NULL, '工作台', 'DataBoard', 1, 0, 1, 1),
  (10, 0, 1, '/system', 'system', 'layouts/default', '/system/user', '系统管理', 'Setting', 10, 0, 0, 1),
  (11, 10, 2, 'user', 'systemUser', 'views/system/user/index.vue', NULL, '用户管理', 'User', 1, 0, 0, 1),
  (12, 10, 2, 'role', 'systemRole', 'views/system/role/index.vue', NULL, '角色管理', 'Avatar', 2, 0, 0, 1),
  (13, 10, 2, 'menu', 'systemMenu', 'views/system/menu/index.vue', NULL, '菜单管理', 'Menu', 3, 0, 0, 1),
  (14, 10, 2, 'api', 'systemApi', 'views/system/api/index.vue', NULL, 'API管理', 'Link', 4, 0, 0, 1),
  (15, 10, 2, 'dictionary', 'systemDictionary', 'views/system/dict/index.vue', NULL, '字典管理', 'Collection', 5, 0, 0, 1),
  (16, 10, 2, 'param', 'systemParam', 'views/system/param/index.vue', NULL, '系统参数', 'Tools', 6, 0, 0, 1),
  (20, 0, 1, '/audit', 'audit', 'layouts/default', '/audit/login-log', '审计日志', 'Document', 20, 0, 0, 1),
  (21, 20, 2, 'login-log', 'auditLoginLog', 'views/audit/login-log/index.vue', NULL, '登录日志', 'Tickets', 1, 0, 0, 1),
  (22, 20, 2, 'operation-log', 'auditOperationLog', 'views/audit/operation-log/index.vue', NULL, '操作日志', 'Notebook', 2, 0, 0, 1);

INSERT INTO `gbp_role_menus` (`role_id`, `menu_id`)
SELECT 1, `id` FROM `gbp_menus`;

INSERT INTO `gbp_api_skip_rules` (`api_path`, `api_method`, `skip_reason`) VALUES
  ('/api/v1/auth/login', 'POST', '登录接口'),
  ('/api/v1/auth/refresh', 'POST', '刷新Token接口'),
  ('/api/v1/health', 'GET', '健康检查接口');

INSERT INTO `gbp_dictionaries`
  (`id`, `dict_name`, `dict_code`, `dict_status`, `remark`)
VALUES
  (1, '用户状态', 'user_status', 1, '用户启用/冻结状态'),
  (2, '角色状态', 'role_status', 1, '角色启用/禁用状态'),
  (3, '菜单类型', 'menu_type', 1, '目录、菜单、隐藏路由、外链'),
  (4, '参数类型', 'param_type', 1, '系统参数值类型');

INSERT INTO `gbp_dictionary_items`
  (`dict_id`, `item_label`, `item_value`, `sort_no`)
VALUES
  (1, '正常', '1', 1),
  (1, '冻结', '2', 2),
  (2, '启用', '1', 1),
  (2, '禁用', '2', 2),
  (3, '目录', '1', 1),
  (3, '菜单', '2', 2),
  (3, '隐藏路由', '3', 3),
  (3, '外链', '4', 4),
  (4, '字符串', 'string', 1),
  (4, '数字', 'number', 2),
  (4, '布尔值', 'boolean', 3),
  (4, 'JSON', 'json', 4);

INSERT INTO `gbp_system_params`
  (`param_name`, `param_key`, `param_value`, `param_type`, `is_encrypted`, `param_status`, `remark`)
VALUES
  ('系统名称', 'system.name', 'GoBaseProject', 'string', 0, 1, '前端页面展示的系统名称'),
  ('访问Token有效期', 'auth.access_token_ttl_minutes', '120', 'number', 0, 1, '单位：分钟'),
  ('刷新Token有效期', 'auth.refresh_token_ttl_minutes', '10080', 'number', 0, 1, '单位：分钟'),
  ('登录失败锁定开关', 'security.login_lock_enabled', 'true', 'boolean', 0, 1, '是否启用登录失败锁定'),
  ('操作日志记录响应体', 'audit.record_response_body', 'false', 'boolean', 0, 1, '是否记录接口响应体');
