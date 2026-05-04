-- 006_seed_audit_apis.sql
-- Audit log API resources and super_admin permission policies

INSERT INTO `gbp_api_resources` (`api_path`, `api_method`, `api_group`, `api_desc`, `api_status`) VALUES
  ('/api/v1/audit/login-logs',              'GET',    'audit', '登录日志列表', 1),
  ('/api/v1/audit/login-logs/cleanup',      'DELETE', 'audit', '登录日志清理', 1),
  ('/api/v1/audit/operation-logs',          'GET',    'audit', '操作日志列表', 1),
  ('/api/v1/audit/operation-logs/{id}',     'GET',    'audit', '操作日志详情', 1),
  ('/api/v1/audit/operation-logs/cleanup',  'DELETE', 'audit', '操作日志清理', 1)
ON DUPLICATE KEY UPDATE api_group = VALUES(api_group), api_desc = VALUES(api_desc), api_status = 1, deleted_at = NULL;

-- Grant super_admin (role_id=1) allow policies for audit APIs
INSERT INTO `gbp_permission_policies`
  (`subject_type`, `subject_id`, `resource_type`, `resource_id`, `resource_key`, `action`, `effect`, `policy_status`)
SELECT 'role', 1, 'api', id, CONCAT(api_method, ':', api_path), api_method, 'allow', 1
FROM `gbp_api_resources`
WHERE `api_group` = 'audit' AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `deleted_at` = NULL, `policy_status` = 1;
