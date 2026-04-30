SET NAMES utf8mb4;
USE `go_base_project`;

-- =========================================================
-- GoBaseProject - 初始化数据
-- 注意：password_hash 仅为开发占位；生产环境必须重新生成。
-- =========================================================

INSERT INTO `gbp_roles`
  (`id`, `role_code`, `role_name`, `parent_role_id`, `default_route`, `sort_no`, `role_status`, `remark`)
VALUES
  (1, 'super_admin', '超级管理员', 0, 'dashboard', 1, 1, '系统内置最高权限角色'),
  (2, 'system_admin', '系统管理员', 1, 'dashboard', 2, 1, '负责用户、角色、菜单、权限等配置'),
  (3, 'operator', '运营人员', 2, 'dashboard', 3, 1, '普通运营人员角色');

INSERT INTO `gbp_users`
  (`id`, `user_uuid`, `login_name`, `password_hash`, `display_name`, `primary_role_id`, `user_status`, `must_change_password`, `remark`)
VALUES
  (1, '00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$REPLACE_WITH_REAL_BCRYPT_HASH', '超级管理员', 1, 1, 1, '系统初始化管理员');

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
  (15, 10, 2, 'dictionary', 'systemDictionary', 'views/system/dictionary/index.vue', NULL, '字典管理', 'Collection', 5, 0, 0, 1),
  (16, 10, 2, 'param', 'systemParam', 'views/system/param/index.vue', NULL, '系统参数', 'Tools', 6, 0, 0, 1),
  (20, 0, 1, '/audit', 'audit', 'layouts/default', '/audit/login-log', '审计日志', 'Document', 20, 0, 0, 1),
  (21, 20, 2, 'login-log', 'auditLoginLog', 'views/audit/login-log/index.vue', NULL, '登录日志', 'Tickets', 1, 0, 0, 1),
  (22, 20, 2, 'operation-log', 'auditOperationLog', 'views/audit/operation-log/index.vue', NULL, '操作日志', 'Notebook', 2, 0, 0, 1);

INSERT INTO `gbp_role_menus` (`role_id`, `menu_id`)
SELECT 1, `id` FROM `gbp_menus`;

INSERT INTO `gbp_role_menus` (`role_id`, `menu_id`) VALUES
  (2, 1), (2, 2), (2, 10), (2, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 20), (2, 21), (2, 22),
  (3, 1), (3, 2);

INSERT INTO `gbp_menu_actions`
  (`id`, `menu_id`, `action_code`, `action_name`, `action_desc`, `sort_no`)
VALUES
  (1, 11, 'add', '新增用户', '用户管理-新增', 1),
  (2, 11, 'edit', '编辑用户', '用户管理-编辑', 2),
  (3, 11, 'delete', '删除用户', '用户管理-删除', 3),
  (4, 11, 'reset_password', '重置密码', '用户管理-重置密码', 4),
  (5, 12, 'add', '新增角色', '角色管理-新增', 1),
  (6, 12, 'edit', '编辑角色', '角色管理-编辑', 2),
  (7, 12, 'delete', '删除角色', '角色管理-删除', 3),
  (8, 12, 'assign_menu', '分配菜单', '角色管理-菜单授权', 4),
  (9, 12, 'assign_api', '分配API', '角色管理-API授权', 5),
  (10, 13, 'add', '新增菜单', '菜单管理-新增', 1),
  (11, 13, 'edit', '编辑菜单', '菜单管理-编辑', 2),
  (12, 13, 'delete', '删除菜单', '菜单管理-删除', 3),
  (13, 14, 'add', '新增API', 'API管理-新增', 1),
  (14, 14, 'edit', '编辑API', 'API管理-编辑', 2),
  (15, 14, 'delete', '删除API', 'API管理-删除', 3),
  (16, 15, 'add', '新增字典', '字典管理-新增', 1),
  (17, 15, 'edit', '编辑字典', '字典管理-编辑', 2),
  (18, 15, 'delete', '删除字典', '字典管理-删除', 3);

INSERT INTO `gbp_role_actions` (`role_id`, `menu_id`, `action_id`)
SELECT 1, `menu_id`, `id` FROM `gbp_menu_actions`;

INSERT INTO `gbp_role_actions` (`role_id`, `menu_id`, `action_id`)
SELECT 2, `menu_id`, `id` FROM `gbp_menu_actions`;

INSERT INTO `gbp_role_data_scopes` (`role_id`, `visible_role_id`) VALUES
  (1, 1), (1, 2), (1, 3),
  (2, 2), (2, 3),
  (3, 3);

INSERT INTO `gbp_api_resources`
  (`id`, `api_path`, `api_method`, `api_group`, `api_desc`)
VALUES
  (1, '/api/v1/auth/login', 'POST', '认证鉴权', '登录'),
  (2, '/api/v1/auth/logout', 'POST', '认证鉴权', '退出'),
  (3, '/api/v1/auth/refresh', 'POST', '认证鉴权', '刷新Token'),
  (4, '/api/v1/auth/profile', 'GET', '认证鉴权', '当前用户信息'),
  (5, '/api/v1/auth/routes', 'GET', '认证鉴权', '当前用户动态路由'),
  (6, '/api/v1/auth/actions', 'GET', '认证鉴权', '当前用户按钮权限'),
  (10, '/api/v1/users', 'GET', '用户管理', '用户分页列表'),
  (11, '/api/v1/users', 'POST', '用户管理', '创建用户'),
  (12, '/api/v1/users/:id', 'GET', '用户管理', '用户详情'),
  (13, '/api/v1/users/:id', 'PUT', '用户管理', '修改用户'),
  (14, '/api/v1/users/:id', 'DELETE', '用户管理', '删除用户'),
  (15, '/api/v1/users/:id/status', 'PUT', '用户管理', '修改状态'),
  (16, '/api/v1/users/:id/password', 'PUT', '用户管理', '重置密码'),
  (17, '/api/v1/users/:id/roles', 'PUT', '用户管理', '分配角色'),
  (20, '/api/v1/roles', 'GET', '角色管理', '角色列表'),
  (21, '/api/v1/roles', 'POST', '角色管理', '创建角色'),
  (22, '/api/v1/roles/:id', 'GET', '角色管理', '角色详情'),
  (23, '/api/v1/roles/:id', 'PUT', '角色管理', '修改角色'),
  (24, '/api/v1/roles/:id', 'DELETE', '角色管理', '删除角色'),
  (25, '/api/v1/roles/tree', 'GET', '角色管理', '角色树'),
  (26, '/api/v1/roles/:id/menus', 'PUT', '角色管理', '菜单授权'),
  (27, '/api/v1/roles/:id/actions', 'PUT', '角色管理', '按钮授权'),
  (28, '/api/v1/roles/:id/apis', 'PUT', '角色管理', 'API授权'),
  (29, '/api/v1/roles/:id/data-scopes', 'PUT', '角色管理', '数据权限'),
  (30, '/api/v1/menus', 'GET', '菜单路由', '菜单列表'),
  (31, '/api/v1/menus/tree', 'GET', '菜单路由', '菜单树'),
  (32, '/api/v1/menus', 'POST', '菜单路由', '创建菜单'),
  (33, '/api/v1/menus/:id', 'PUT', '菜单路由', '修改菜单'),
  (34, '/api/v1/menus/:id', 'DELETE', '菜单路由', '删除菜单'),
  (40, '/api/v1/audit/login-logs', 'GET', '审计日志', '登录日志分页'),
  (41, '/api/v1/audit/operation-logs', 'GET', '审计日志', '操作日志分页');

INSERT INTO `gbp_api_skip_rules` (`api_path`, `api_method`, `skip_reason`) VALUES
  ('/api/v1/auth/login', 'POST', '登录接口'),
  ('/api/v1/auth/refresh', 'POST', '刷新Token接口'),
  ('/api/v1/health', 'GET', '健康检查接口');

INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`)
SELECT 'role', 1, 'api', `id`, CONCAT(`api_method`, ':', `api_path`), '*', 'allow'
FROM `gbp_api_resources`;

INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`)
SELECT 'role', 2, 'api', `id`, CONCAT(`api_method`, ':', `api_path`), '*', 'allow'
FROM `gbp_api_resources`
WHERE `api_group` IN ('认证鉴权', '用户管理', '角色管理', '菜单路由', '审计日志');

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
