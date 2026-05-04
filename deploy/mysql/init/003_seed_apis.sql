-- 003_seed_apis.sql
-- Seeds API resource catalog and grants super_admin (role_id=1) allow-all policies.
-- Path params are stored as {id} regardless of the actual route variable name,
-- so the middleware's normalizePath function can match them uniformly.

INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  -- auth
  ('/api/v1/auth/logout',               'POST',   'auth',   '退出登录',       1),
  ('/api/v1/auth/profile',              'GET',    'auth',   '当前用户信息',   1),
  ('/api/v1/auth/routes',               'GET',    'auth',   '菜单路由',       1),
  ('/api/v1/auth/actions',              'GET',    'auth',   '按钮权限',       1),
  -- user
  ('/api/v1/users',                     'GET',    'user',   '用户列表',       1),
  ('/api/v1/users',                     'POST',   'user',   '创建用户',       1),
  ('/api/v1/users/role-options',        'GET',    'user',   '角色选项',       1),
  ('/api/v1/users/{id}',                'GET',    'user',   '用户详情',       1),
  ('/api/v1/users/{id}',                'PUT',    'user',   '修改用户',       1),
  ('/api/v1/users/{id}',                'DELETE', 'user',   '删除用户',       1),
  ('/api/v1/users/{id}/status',         'PUT',    'user',   '修改用户状态',   1),
  ('/api/v1/users/{id}/password',       'PUT',    'user',   '重置密码',       1),
  ('/api/v1/users/{id}/roles',          'PUT',    'user',   '分配角色',       1),
  -- role
  ('/api/v1/roles',                     'GET',    'role',   '角色列表',       1),
  ('/api/v1/roles',                     'POST',   'role',   '创建角色',       1),
  ('/api/v1/roles/tree',                'GET',    'role',   '角色树',         1),
  ('/api/v1/roles/resources',           'GET',    'role',   '权限资源目录',   1),
  ('/api/v1/roles/{id}',                'GET',    'role',   '角色详情',       1),
  ('/api/v1/roles/{id}',                'PUT',    'role',   '修改角色',       1),
  ('/api/v1/roles/{id}',                'DELETE', 'role',   '删除角色',       1),
  ('/api/v1/roles/{id}/menus',          'GET',    'role',   '角色菜单',       1),
  ('/api/v1/roles/{id}/menus',          'PUT',    'role',   '分配菜单',       1),
  ('/api/v1/roles/{id}/actions',        'GET',    'role',   '角色按钮权限',   1),
  ('/api/v1/roles/{id}/actions',        'PUT',    'role',   '分配按钮权限',   1),
  ('/api/v1/roles/{id}/apis',           'GET',    'role',   '角色API权限',    1),
  ('/api/v1/roles/{id}/apis',           'PUT',    'role',   '分配API权限',    1),
  ('/api/v1/roles/{id}/data-scopes',    'GET',    'role',   '角色数据范围',   1),
  ('/api/v1/roles/{id}/data-scopes',    'PUT',    'role',   '分配数据范围',   1),
  -- menu
  ('/api/v1/menus',                     'GET',    'menu',   '菜单列表',       1),
  ('/api/v1/menus',                     'POST',   'menu',   '创建菜单',       1),
  ('/api/v1/menus/tree',                'GET',    'menu',   '菜单树',         1),
  ('/api/v1/menus/{id}',                'GET',    'menu',   '菜单详情',       1),
  ('/api/v1/menus/{id}',                'PUT',    'menu',   '修改菜单',       1),
  ('/api/v1/menus/{id}',                'DELETE', 'menu',   '删除菜单',       1),
  ('/api/v1/menus/{id}/params',         'GET',    'menu',   '菜单路由参数',   1),
  ('/api/v1/menus/{id}/params',         'POST',   'menu',   '新增路由参数',   1),
  ('/api/v1/menus/params/{id}',         'DELETE', 'menu',   '删除路由参数',   1),
  ('/api/v1/menus/{id}/actions',        'GET',    'menu',   '菜单按钮权限',   1),
  ('/api/v1/menus/{id}/actions',        'POST',   'menu',   '新增按钮权限',   1),
  ('/api/v1/menu-actions/{id}',         'PUT',    'menu',   '修改按钮权限',   1),
  ('/api/v1/menu-actions/{id}',         'DELETE', 'menu',   '删除按钮权限',   1),
  -- system
  ('/api/v1/system/runtime',            'GET',    'system', '系统运行信息',   1)
ON DUPLICATE KEY UPDATE api_group = VALUES(api_group), api_desc = VALUES(api_desc), api_status = 1, deleted_at = NULL;

-- Grant super_admin (role_id = 1) allow policy for every API resource.
INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT
  'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `deleted_at` = NULL, `policy_status` = 1;
