-- 007_seed_api_mgmt_apis.sql
-- API 管理和白名单管理接口资源 + super_admin 权限策略

INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  ('/api/v1/apis/groups',         'GET',    'api', 'API分组列表',   1),
  ('/api/v1/apis',                'GET',    'api', 'API资源列表',   1),
  ('/api/v1/apis',                'POST',   'api', '创建API资源',   1),
  ('/api/v1/apis/{id}',           'PUT',    'api', '修改API资源',   1),
  ('/api/v1/apis/{id}',           'DELETE', 'api', '删除API资源',   1),
  ('/api/v1/api-skip-rules',      'GET',    'api', '白名单列表',    1),
  ('/api/v1/api-skip-rules',      'POST',   'api', '创建白名单',    1),
  ('/api/v1/api-skip-rules/{id}', 'DELETE', 'api', '删除白名单',    1)
ON DUPLICATE KEY UPDATE api_group = VALUES(api_group), api_desc = VALUES(api_desc), api_status = 1, deleted_at = NULL;

-- Grant super_admin (role_id=1) allow policies
INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT 'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `api_group` = 'api' AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `deleted_at` = NULL, `policy_status` = 1;
